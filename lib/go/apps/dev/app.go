package main

import (
	"log"
	"github.com/valyala/fasthttp"
)

var (
	connection                       = &fasthttp.Client{}
	req                              = fasthttp.AcquireRequest()
	resp                             = fasthttp.AcquireResponse()
)

func main() {
	log.Println("Start reloading model on prediction server")

	_response := doRequest("http://0.0.0.0:8081/reload")

	if string(_response) == "reloaded" {
		log.Println("Prediction model reloaded on python server")
	} else {
		log.Println("Problem while reloading prediction model on python server", string(_response))
	}
}

func doRequest(url string) []byte {

	req.SetRequestURI(url)

	connection.Do(req, resp)

	return resp.Body()
}