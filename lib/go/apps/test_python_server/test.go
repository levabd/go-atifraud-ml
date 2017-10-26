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

// To correctly perform this test you need:
// -------------------------------------------------------------------------------
// 1 - decide what education sample will be (for example 100 000 or 90 000)
// 2 - find rows which make this sample (our sample - 90 000) in DB:
//     example of queries (different periods of time - windows):
//     SELECT count(*)
//     FROM logs
//     -- where timestamp > '2017-08-08 00:00:00' and timestamp < '2017-08-28 15:00:00' - education(train) sample
//     -- where timestamp > '2017-08-28 15:00:00' and timestamp < '2017-09-02 08:00:00' - sample of bots
//     -- where timestamp > '2017-09-02 08:00:00' and timestamp < '2017-09-09 08:00:00' - sample of human
//
//     -- where timestamp > '2017-08-28 00:00:00' and timestamp < '2017-09-07 23:00:00' - education(train) sample
//     -- where timestamp > '2017-09-07 23:00:00' and timestamp < '2017-09-10 23:00:00' - sample of bots
//     -- where timestamp > '2017-09-10 23:00:00' and timestamp < '2017-09-13 10:00:00' - sample of human
//
//     -- where timestamp > '2017-09-07 23:00:00' and timestamp < '2017-09-15 14:00:00' - education(train) sample
//     -- where timestamp > '2017-09-15 14:00:00' and timestamp < '2017-09-18 12:00:00' - sample of bots
//     -- where timestamp > '2017-09-18 12:00:00' and timestamp < '2017-09-20 17:00:00' - sample of human
//
//     -- where timestamp > '2017-09-15 14:00:00' and timestamp < '2017-09-23 08:00:00' - education(train) sample
//     -- where timestamp > '2017-09-23 08:00:00' and timestamp < '2017-09-26 00:00:00' - sample of bots
//     -- where timestamp > '2017-09-26 00:00:00' and timestamp < '2017-09-28 17:00:00' - sample of human
//
//     -- where timestamp > '2017-09-23 08:00:00' and timestamp < '2017-10-01 14:00:00' - education(train) sample
//     -- where timestamp > '2017-10-01 14:00:00' and timestamp < '2017-10-03 22:00:00' - sample of bots
//     -- where timestamp > '2017-10-03 22:00:00' and timestamp < '2017-10-06 13:00:00' - sample of human
//
// 3 - Change lib/go/apps/prepare_data_for_train/prepare_data_for_train.go:36 by setting one of the needed period (in our case - education(train) sample, see above in query example)
// 4 - Run lib/go/apps/prepare_data_for_train/prepare_data_for_train.go
//     go run lib/go/apps/prepare_data_for_train/prepare_data_for_train.go
//
// 5 - Restart prediction_server (lib/python/prediction_server) - to reload model
// 6 - Change lib/go/apps/prepare_data_for_test/prepare_data_for_test.go:44 by setting one of the needed period (in our case - sample of bots, see above in query example)
// 7 - Run lib/go/apps/prepare_data_for_test/prepare_data_for_test.go - bots will be generated
//     go run lib/go/apps/prepare_data_for_test/prepare_data_for_test.go
// 8 - Change this file (lib/go/apps/test_python_server/test.go:92) by setting one of the needed period (in our case - sample of human, see above in query example)
// 9 - Run lib/go/apps/test_python_server/test.go
//     go run lib/go/apps/test_python_server/test.go
//
// 10- Analise the results


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

	trimmedValue, trimmedOrder := trimData(valueData, orderData)
	fullFeatures := services.GetSingleFullFeatures(trimmedOrder, trimmedValue, valuesFeaturesOrder)

	// transform to sparse matrix
	featuresSparseMatrix := getSparseMatrix(fullFeatures)

	predictionResults := getPredictionResults(featuresSparseMatrix)

	for _, obj := range predictionResults {
		for key, prediction := range obj {
			if prediction <= threshold {
				continue
			}
			if key == uaFamilyCode {
				return "human"
			}
		}
	}

	return "bot"
}
func getSparseMatrix(fullFeatures []float64) []string {
	var featuresSparseMatrix []string
	for index_column, value := range fullFeatures {
		if value == 1 {
			featuresSparseMatrix = append(featuresSparseMatrix, fmt.Sprintf("%v", index_column))
		}
	}
	return featuresSparseMatrix
}

func getPredictionResults(sparseArray []string) []map[string]float64 {
	var predictionResults []map[string]float64
	_response := doRequest("http://0.0.0.0:8081/?positions=" + strings.Join(sparseArray, ","))
	err := json.Unmarshal([]byte(_response), &predictionResults)
	if err != nil {
		panic(err)
	}
	return predictionResults
}

func doRequest(url string) []byte {

	req.SetRequestURI(url)

	connection.Do(req, resp)

	return resp.Body()
}

func handleLogLine(line []byte) (string, map[string]interface{}, map[string]interface{}) {

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
