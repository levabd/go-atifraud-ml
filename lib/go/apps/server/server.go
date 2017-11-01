package main

import (
	"flag"
	"log"
	"math/rand"
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
	"errors"
)

type PredictionBackends struct {
	servers []string
	n       int
}

func (b *PredictionBackends) Choose() string {
	idx := b.n % len(b.servers)
	b.n++
	return b.servers[idx]
}

func (b *PredictionBackends) String() string {
	return strings.Join(b.servers, ", ")
}

var (
	addr                = flag.String("addr", "0.0.0.0:8082", "host:port to listen to")
	compress            = flag.Bool("compress", false, "Whether to enable transparent response compression")
	predictionServers   = flag.String("predictionServers", "", "The backend servers to predictionServers connections across, separated by commas")
	valuesFeaturesOrder = s.LoadFittedValuesFeaturesOrder()
	db                  *gorm.DB
	udger_instance      *udger.Udger = nil
	connections         []fasthttp.Client
	reqs                []fasthttp.Request
	resps               []fasthttp.Response
	req                 = *fasthttp.AcquireRequest()
	resp                = *fasthttp.AcquireResponse()
	numberOfConnections = 8
	predictionBackends  *PredictionBackends
)

const (
	PredictionThreshold = 0.03
	ReturnCrawlerWord   = "crawler"
	ReturnHumanWord     = "human"
	ReturnErrorWord     = "error"
	ReturnBotWord       = "bot"
	ReturnUnknownWord   = "unknown_no_user_agent"
)

func init() {

	flag.Parse()

	if *predictionServers == "" {
		log.Println("You forget to spesicfy predictionServers option")
		return
	}

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

	// create 4 connections
	servers := strings.Split(*predictionServers, ",")
	if len(servers) == 1 && servers[0] == "" {
		log.Fatalln("please specify backend servers with -predictionBackends")
	}
	predictionBackends = &PredictionBackends{servers: servers}

	numberOfConnections := len(servers)
	connections = make([]fasthttp.Client, numberOfConnections)
	reqs = make([]fasthttp.Request, numberOfConnections)
	resps = make([]fasthttp.Response, numberOfConnections)

	for i := 0; i < numberOfConnections; i++ {
		connections[i] = fasthttp.Client{}
		reqs[i] = *fasthttp.AcquireRequest()
		resps[i] = *fasthttp.AcquireResponse()
	}
}

func main() {

	if *predictionServers == "" {
		return
	}

	defer db.Close()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	log.Println("Strarted listening server on address:", *addr)
	log.Println("Prediction requests will be send set of servers:", *predictionServers)

	server := &fasthttp.Server{
		Handler:        h,
		Concurrency:    1000,
		ReadBufferSize: 3000,
	}

	if err := server.ListenAndServe(*addr); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {

	if ! ctx.IsPost() {
		log.Printf("This is must be post request - will not be handled", )
		return
	}

	var body = ctx.PostBody()
	var clientIp = ctx.Request.Header.Peek("X-Real-IP")

	if len(body) == 0 {
		log.Println("Post body is")
		return
	}

	var agent, err = handleHeader(clientIp, body)
	if err != nil {
		ctx.Response.SetBodyString(fmt.Sprintln(err))
		ctx.SetStatusCode(424)
		return
	}

	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.SetBodyString(agent)
}

func handleHeader(clientIp []byte, response []byte) (string, error) {

	userAgent, valueData, orderData := handleLogLine(response)

	if userAgent == "" && string(clientIp) == "" {
		return ReturnUnknownWord, nil
	}
	var parsedData = udger_instance.GetNewParsedData()

	if udger_instance.IsCrawler(string(clientIp), userAgent, parsedData) {
		return ReturnCrawlerWord, nil
	}

	_isHuman, err := isHuman(valueData, orderData, parsedData)
	if err != nil {
		return ReturnErrorWord, err
	}

	if _isHuman {
		return ReturnHumanWord, nil
	}

	return ReturnBotWord, nil
}

func isHuman(valueData map[string]interface{}, orderData map[string]interface{}, parsedData map[string]map[string]string, ) (bool, error) {

	trimmedValue, trimmedOrder := trimData(valueData, orderData)
	fullFeatures := s.GetSingleFullFeatures(trimmedOrder, trimmedValue, valuesFeaturesOrder)

	// transform to sparse matrix
	var sparseArray []string
	for index_column, value := range fullFeatures {
		if value == 1 {
			sparseArray = append(sparseArray, fmt.Sprintf("%v", index_column))
		}
	}

	predictionResults, err := getPredictionResults(sparseArray)
	if err != nil {
		log.Println(err)
		return false, err
	}

	for _, obj := range predictionResults {
		for key, prediction := range obj {
			if prediction.(float64) <= PredictionThreshold {
				continue
			}
			if key == parsedData["user_agent"]["ua_family_code"] {
				return true, nil
			}
		}
	}

	return false, nil
}

func getPredictionResults(sparseArray []string) ([]map[string]interface{}, error) {

	var predictionResults []map[string]interface{}

	if len(sparseArray) == 0 {
		return predictionResults, nil
	}

	var params = strings.Join(sparseArray, ",")

	_response := doRequest(params)

	if len(_response) == 0 {
		return predictionResults, errors.New(fmt.Sprintf("Prediction server responde with empty body, params in request: %s", params))
	}

	err := json.Unmarshal([]byte(_response), &predictionResults)
	if err != nil {
		return predictionResults, errors.New(fmt.Sprintf("Uncorrect JSON from prediction server response: %s", err))
	}

	return predictionResults, nil
}

func doRequest(params string) []byte {

	r := rand.Intn(numberOfConnections)

	reqs[r].SetRequestURI("http://" + predictionBackends.Choose() + "/?positions=" + params)
	connections[r].Do(&reqs[r], &resps[r])

	//defer resps[r].SetConnectionClose()
	//defer reqs[r].SetConnectionClose()

	return resps[r].Body()
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

	if _ua, ok := valueRow["user-agent"].(string); ok && ua == "" {
		ua = _ua
	}

	if _ua, ok := valueRow["user_agent"].(string); ok && ua == "" {
		ua = _ua
	}

	if ua == "" {
		log.Printf("User-agent is absent. Body: %+v", valueRow)
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
