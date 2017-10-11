package main

import (
	"flag"
	"log"
	"github.com/valyala/fasthttp"
	"time"
	"fmt"
	"github.com/buger/jsonparser"
	"encoding/json"
	s "github.com/levabd/go-atifraud-ml/lib/go/services"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
)

var (
	addr     = flag.String("addr", "localhost:8082", "Thost:port to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
	debug    = flag.Bool("debug", false, "Whether to debug  transparent response compression")
	valuesFeaturesOrder = s.LoadFittedValuesFeaturesOrder()
	userAgentIntCodes, userAgentFloatCodes = s.LoadFittedUserAgentCodes()
	userAgentStrings = s.LoadFittedUserAgentDeCoder()
)

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	log.Println("Strarted listening server on address: ", *addr)
	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	if *debug {
		defer timeTrack(time.Now(), "requestHandler")
	}

	if ! ctx.IsPost() {
		log.Printf("This is must be post request - will not be handled", )
		return
	}

	var body = ctx.PostBody()
	var headerBodyIp = ctx.Request.Header.Peek("Body-Header-Ip")

	if len(headerBodyIp) == 0{
		log.Println("`Body-Header-Ip` header was not provided")
		return
	}

	if len(body) == 0 {
		log.Println("Post body is")
		return
	}

	var isCrawler = handleHeader(headerBodyIp, body)

	// todo USE isCrawler
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.SetBodyString(fmt.Sprintf("%t", isCrawler))
	log.Printf("Response is %t", isCrawler)
}

//noinspection GoUnusedParameter
func handleHeader(headerBodyIp []byte, response []byte) bool {
	//defer timeTrack(time.Now(), "handleHeader")

	mainData, valueData, orderData := handleLogLine(response)
	trimmedValue, trimmedOrder := trimData(valueData, orderData)

	//fmt.Printf("value_data :%+v\n", value_data)
	//fmt.Printf("trimmed_order :%+v\n", trimmed_order)

	fullFeatures := s.GetSingleFullFeatures(trimmedOrder, trimmedValue, valuesFeaturesOrder)

	// todo remove debug
	fmt.Println(mainData)
	fmt.Println(len(fullFeatures))

	// todo make prediction
	return len(fullFeatures)>0
}

func handleLogLine( line []byte) (string, map[string]interface{}, map[string]interface{}) {
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

	return "", valueRow, orderedRow
}

func trimData(valueData map[string]interface{}, orderData map[string]interface{}) (map[string]interface{}, map[string]interface{}) {

	var headerModel = m.Log{ValueData: valueData, OrderData: orderData}
	return headerModel.TrimValueData(), headerModel.TrimOrderData()
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}