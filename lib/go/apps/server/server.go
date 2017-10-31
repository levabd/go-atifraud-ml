package main

import (
	"flag"
	"log"
	"github.com/valyala/fasthttp"
	"fmt"
	"encoding/json"
	s "github.com/levabd/go-atifraud-ml/lib/go/services"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
	"strings"
	"github.com/jinzhu/gorm"
	"os"
	"path/filepath"
	"github.com/levabd/go-atifraud-ml/lib/go/udger"
	"sync"
)

var (
	addr                 = flag.String("addr", "0.0.0.0:8082", "host:port to listen to")
	addrPredictionServer = flag.String("addrPredictionServer", "0.0.0.0:8081", "host:port to make request on prediction server")
	compress             = flag.Bool("compress", false, "Whether to enable transparent response compression")
	valuesFeaturesOrder  = s.LoadFittedValuesFeaturesOrder()
	db                   *gorm.DB
	udger_instance       *udger.Udger = nil
	connection                        = &fasthttp.Client{}
	req                               = *fasthttp.AcquireRequest()
	resp                              = *fasthttp.AcquireResponse()
	mux                  sync.Mutex
)

const (
	PredictionThreshold = 0.03
	ReturnCrawlerWord   = "crawler"
	ReturnHumanWord     = "human"
	ReturnBotWord       = "bot"
	ReturnUnknownWord   = "unknown_no_user_agent"
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
	log.Println("Prediction requests will be send to: ", *addrPredictionServer)

	server := &fasthttp.Server{
		Handler:            h,
		Concurrency:        100,
		MaxRequestsPerConn: 2000,
		ReadBufferSize:     2000,
	}

	if err := server.ListenAndServe(*addr); err != nil {
		log.Fatalf("Error in ListenAndServe: %server", err)
	}
	defer req.ConnectionClose()
	defer resp.ConnectionClose()
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

	ctx.Response.SetBodyString(agent)
}

func handleHeader(response []byte) string {
	userAgent, valueData, orderData := handleLogLine(response)

	if userAgent == "" {
		return ReturnUnknownWord
	}

	if udger_instance.IsCrawler("", userAgent, true) {
		return ReturnCrawlerWord
	}

	if isHuman(valueData, orderData) {
		return ReturnHumanWord
	}

	return ReturnBotWord
}

func isHuman(valueData map[string]interface{}, orderData map[string]interface{}) bool {

	mux.Lock()
	defer mux.Unlock()

	trimmedValue, trimmedOrder := trimData(valueData, orderData)
	fullFeatures := s.GetSingleFullFeatures(trimmedOrder, trimmedValue, valuesFeaturesOrder)

	// transform to sparse matrix
	var sparseArray []string
	for index_column, value := range fullFeatures {
		if value == 1 {
			sparseArray = append(sparseArray, fmt.Sprintf("%v", index_column))
		}
	}

	predictionResults := getPredictionResults(sparseArray)
	for _, obj := range predictionResults {
		for key, prediction := range obj {
			if prediction <= PredictionThreshold {
				continue
			}
			if key == udger_instance.ParseData["user_agent"]["ua_family_code"] {
				return true
			}
		}
	}

	return false
}

func getPredictionResults(sparseArray []string) []map[string]float64 {

	var predictionResults []map[string]float64

	if len(sparseArray) == 0 {
		return predictionResults
	}

	_response := doRequest("http://" + *addrPredictionServer + "/?positions=" + strings.Join(sparseArray, ","))

	err := json.Unmarshal([]byte(_response), &predictionResults)

	if err != nil {
		log.Fatalln(err)
		return predictionResults
	}

	return predictionResults
}

func doRequest(url string) []byte {

	req.SetRequestURI("")

	req.SetRequestURI(url)

	connection.Do(&req, &resp)

	return resp.Body()
}

func handleLogLine(headers []byte) (ua string, valueRow map[string]interface{}, orderRow map[string]interface{}) {

	valueRow = make(map[string]interface{})
	orderRow = make(map[string]interface{})

	if len(headers) == 0 {
		log.Printf("Request body is empty. Body: %s", string(headers))
		return "", valueRow, orderRow
	}

	line_splitted := strings.Split(string(headers), "\n")

	if len(line_splitted) == 0 {
		log.Printf("Can't split request body with \n symbol. Body: %s", string(headers))
		return "", valueRow, orderRow
	}

	for index, _line := range line_splitted {
		i := strings.Index(_line, ":")
		if i > -1 {
			orderRow[ _line[:i] ] = index
			valueRow[ _line[:i] ] = _line[i+1:]
		} else {
			log.Printf("Index of ':' not found in _line: %s", _line)
		}
	}

	ua = ""

	if _ua, ok := valueRow["User-Agent"].(string); ok {
		ua = _ua
	}

	if _ua, ok := valueRow["user-agent"].(string); ok && ua != "" {
		ua = _ua
	}

	if _ua, ok := valueRow["user_agent"].(string); ok && ua != "" {
		ua = _ua
	}

	if ua == "" {
		log.Printf("User-agent is absent. Body: %s", string(headers))
		return "", valueRow, orderRow
	}

	return ua, valueRow, orderRow
}

func trimData(
	valueData map[string]interface{},
	orderData map[string]interface{},
) (map[string]interface{}, map[string]interface{}) {
	var headerModel = m.Log{ValueData: valueData, OrderData: orderData}
	return headerModel.TrimValueData(), headerModel.TrimOrderData()
}
