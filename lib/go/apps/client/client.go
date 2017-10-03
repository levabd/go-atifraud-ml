package main

import (
	"github.com/valyala/fasthttp"
	"time"
	"log"

)
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	defer timeTrack(time.Now(), "main")

	c := &fasthttp.Client{}
	c.MaxIdleConnDuration =  100 * time.Second

	var req fasthttp.Request
	req.Header.SetMethod("POST")
	req.SetRequestURI("http://localhost:8082/")
	req.Header.Set("Host", "localhost")
	req.Header.Set("Body-Header-Ip", "62.84.44.222")
	req.SetBodyString(`{"Cache-Control":"no-cache","Connection":"Keep-Alive","Pragma":"no-cache","Accept":"*\/*","Accept-Encoding":"gzip, deflate","From":"bingbot(at)microsoft.com","Host":"www.vypekajem.com","User-Agent":"Mozilla\/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit\/537.51.1 (KHTML, like Gecko) Version\/7.0 Mobile\/11A465 Safari\/9537.53 (compatible; bingbot\/2.0; +http:\/\/www.bing.com\/bingbot.htm)"}`)
	var resp fasthttp.Response

	err := c.DoTimeout(&req, &resp, time.Second)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() == 200 { // OK
		useResponseBody(resp.Body())
	}
}

func useResponseBody(body []byte) {
	if string(body) =="true"{
		log.Println("resp.Body", string(body))
	}
}