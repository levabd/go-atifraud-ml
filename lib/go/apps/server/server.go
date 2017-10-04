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

func handleHeader(headerBodyIp []byte, response []byte) bool {
	//defer timeTrack(time.Now(), "handleHeader")

	_, value_data, order_data := handleLogLine(response)
	_, trimmed_order := trimData(value_data, order_data)

	//fmt.Printf("value_data :%+v\n", value_data)
	//fmt.Printf("trimmed_order :%+v\n", trimmed_order)

	pair_dict := s.GetPairDictForHeaders(trimmed_order)

	//fmt.Printf("pair_dict :%+v\n", pair_dict)

	// todo make prediction
	return len(pair_dict)>0
}

func handleLogLine( line []byte) (string, map[string]interface{}, map[string]interface{}) {
	var (
		value_row   map[string]interface{} = make(map[string]interface{})
		ordered_row map[string]interface{} = make(map[string]interface{})
	)

	if len(line) == 0 || ! isJSON(line) {
		return "", value_row, ordered_row
	}

	i := 0
	jsonparser.ObjectEach(line, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		ordered_row[ string(key)] = i
		value_row[ string(key) ] = string(value)
		i = i + 1
		return nil
	})

	// define crawler in User-Agent
	if user_agent, ok := value_row["User-Agent"].(string); ok {
		return user_agent, value_row, ordered_row
	} else {
		return "", value_row, ordered_row
	}

	return "", value_row, ordered_row
}

func trimData(value_data map[string]interface{}, order_data map[string]interface{}) (map[string]interface{}, map[string]interface{}) {

	var header_model = m.Log{ValueData: value_data, OrderData: order_data}
	return header_model.TrimValueData(), header_model.TrimOrderData()
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}