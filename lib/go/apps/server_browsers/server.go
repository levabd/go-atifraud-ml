package main

import  (
	"flag"
	"log"
	"github.com/valyala/fasthttp"
	"time"
	"fmt"
	"github.com/buger/jsonparser"
	"encoding/json"
	s "github.com/levabd/go-atifraud-ml/lib/go/services"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
	"strings"
	"github.com/jinzhu/gorm"
	"os"
	"path/filepath"
	"github.com/levabd/go-atifraud-ml/lib/go/udger"
)

var (
	addr                = flag.String("addr", "localhost:8082", "Thost:port to listen to")
	compress            = flag.Bool("compress", false, "Whether to enable transparent response compression")
	valuesFeaturesOrder = s.LoadFittedValuesFeaturesOrder()
	uaVersionStrings    = s.LoadFittedUaVersionDeCoder()
	db                  *gorm.DB
	udger_instance      *udger.Udger = nil
	connection                       = &fasthttp.Client{}
    req = fasthttp.AcquireRequest()
	resp = fasthttp.AcquireResponse()
)

func init() {
	_db, err := gorm.Open("sqlite3", os.Getenv("DB_FILE_PATH_UDGER"))
	if err != nil {
		log.Fatalf("server.go - init: Failed to estabblish database connection: %s", err)
	}
	db = _db

	pathToUdgerDb := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "db", "udgerdb_v3.dat")
	log.Println("path_to_udger_db: ", pathToUdgerDb)

	u, err := udger.New(os.Getenv("DB_FILE_PATH_UDGER"))
	if err != nil {
		log.Fatalf("server.go - init: Failed to initiate udger helper: %s", err)
	}

	udger_instance = u
}

func main() {
	defer db.Close()
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	log.Println("Strarted listening server on address: ", *addr)
	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
	defer  req.ConnectionClose()
	defer  resp.ConnectionClose()
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	if ! ctx.IsPost() {
		log.Printf("This is must be post request - will not be handled", )
		return
	}

	var body = ctx.PostBody()

	if len(body) == 0 {
		log.Println("Post body is")
		return
	}

	var agent = handleHeader(body)

	// todo USE agent
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.SetBodyString(agent)
}

//noinspection GoUnusedParameter
func handleHeader(response []byte) string {
	userAgent, valueData, orderData := handleLogLine(response)

	if userAgent == "" {
		return "unknown_no_user_agent"
	}

	if udger_instance.IsCrawler("", userAgent, true) {
		return "crawler"
	}

	trimmedValue, trimmedOrder := trimData(valueData, orderData)
	fullFeatures := s.GetSingleFullFeatures(trimmedOrder, trimmedValue, valuesFeaturesOrder)

	// transform to sparse matrix
	var sparseArray []string
	for index_column, value := range fullFeatures {
		if value == 1 {
			sparseArray = append(sparseArray, fmt.Sprintf("%v", index_column))
		}
	}

	_response := doRequest("http://0.0.0.0:8081/?positions=" + strings.Join(sparseArray, ","))
	var predictionResults []map[string]float64
	err := json.Unmarshal([]byte(_response), &predictionResults)
	if err != nil {
		panic(err)
	}

	for _, obj := range predictionResults {
		for key, prediction := range obj {
			if prediction <= 0.03 {
				continue
			}
			if key == udger_instance.ParseData["user_agent"]["ua_family_code"] {
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
	var headerModel = m.Log{ValueData: valueData, OrderData: orderData}
	return headerModel.TrimValueData(), headerModel.TrimOrderData()
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}
