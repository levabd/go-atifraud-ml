package main

import (
	"github.com/jinzhu/gorm"
	"github.com/levabd/go-atifraud-ml/lib/go/models"
	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"strings"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/uniplaces/carbon"
	"time"
)

var (
	bots_not_test_data, humans_not_test_data , bots_test_data, humans_test_data int
	connection                       = &fasthttp.Client{}
	req                              = fasthttp.AcquireRequest()
	resp                             = fasthttp.AcquireResponse()
	valuesFeaturesOrder              = services.LoadFittedValuesFeaturesOrder()
	threshold = 0.03
)

func main() {
	db, err := gorm.Open("postgres", models.GetDBConnectionStr())
	if err != nil {
		panic("user_agent_helpers.go - LoadFittedUserAgentCodes: Failed to connect to database")
	}
	defer db.Close()
	if !db.HasTable(&models.TestLog{}) {
		db.AutoMigrate(&models.TestLog{})
	}

	test_logs := []models.TestLog{}
	db.Order("timestamp").Find(&test_logs)

	for _, log := range test_logs {
		d := log.ValueData
		valueString, _ := json.Marshal(d)
		response :=	handleHeader(valueString, log.UaFamilyCode)

		if response == "human" {
			humans_test_data++
		} else if response == "bot" {
			bots_test_data++
		}
	}

	start, _ := carbon.Create(2017, 9, 10, 23, 0, 0, 0, "Asia/Almaty")
	end, _ := carbon.Create(2017, 9, 13, 10, 0, 0, 0, "Asia/Almaty")

	human_logs := []models.Log{}
	db.Order("timestamp").Where("timestamp BETWEEN ? AND ?",  time.Unix(start.Unix(), 0),time.Unix( end.Unix(), 0)).Find(&human_logs)
	human_logs = human_logs[:30000]

	for _, log := range human_logs {
		d := log.ValueData
		valueString, _ := json.Marshal(d)
		response:=		handleHeader(valueString, log.UaFamilyCode)

		if response == "human" {
			humans_not_test_data++
		} else if response == "bot" {
			bots_not_test_data++
		}
	}

	println(fmt.Sprintf("Threshold %f", threshold))
	println(fmt.Sprintf("Test sample size 60000"))
	println(fmt.Sprintf("1 Error: %f%%. human as bot: human %v, bots %v", float64(bots_not_test_data)/float64(60000)*100, humans_not_test_data, bots_not_test_data))
	println(fmt.Sprintf("2 Error: %f%%. bot as human: human %v, bots %v", float64(humans_test_data)/float64(60000)*100, humans_test_data, bots_test_data))
}

func handleHeader(response []byte, uaFamilyCode string ) string {
	userAgent, valueData, orderData := handleLogLine(response)

	if userAgent == "" {
		return "unknown_no_user_agent"
	}

	//start:= time.Now()
	trimmedValue, trimmedOrder := trimData(valueData, orderData)
	fullFeatures := services.GetSingleFullFeatures(trimmedOrder, trimmedValue, valuesFeaturesOrder)

	// transform to sparse matrix
	var sparseArray []string
	for index_column, value := range fullFeatures {
		if value == 1 {
			sparseArray = append(sparseArray, fmt.Sprintf("%v", index_column))
		}
	}

	var predictionResults []map[string]float64

	_response := doRequest("http://0.0.0.0:8081/?positions=" + strings.Join(sparseArray, ","))
	err := json.Unmarshal([]byte(_response), &predictionResults)
	if err != nil {
		panic(err)
	}

	for _, obj := range predictionResults {
		for key, prediction := range obj {
			if prediction <= threshold {
				continue
			}
			if key == uaFamilyCode  {
				return "human"
			}
		}
	}

	return "bot"
}

func doRequest(url string) []byte {
	req.SetRequestURI(url)
	connection.Do(req, resp)
	return resp.Body()
}

func handleLogLine(line []byte) (string, map[string]interface{}, map[string]interface{}) {
	//defer timeTrack(time.Now(), "handleLogLine")

	var (
		valueRow   map[string]interface{} = make(map[string]interface{})
		orderedRow map[string]interface{} = make(map[string]interface{})
	)

	if len(line) == 0 || ! isJSON(line) {
		return "", valueRow, orderedRow
	}

	i := 0
	jsonparser.ObjectEach(line, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		orderedRow[ string(key)] = i
		valueRow[ string(key) ] = string(value)
		i = i + 1
		return nil
	})

	// define crawler in User-Agent
	if userAgent, ok := valueRow["User-Agent"].(string); ok {
		return userAgent, valueRow, orderedRow
	} else {
		return "", valueRow, orderedRow
	}
}

func trimData(valueData map[string]interface{}, orderData map[string]interface{}) (map[string]interface{}, map[string]interface{}) {
	var headerModel = models.TestLog{ValueData: valueData, OrderData: orderData}
	return headerModel.TrimValueData(), headerModel.TrimOrderData()
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}
