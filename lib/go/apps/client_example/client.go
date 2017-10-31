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
	req.SetRequestURI("http://0.0.0.0:8082/")
	req.Header.Set("Host", "localhost")
	req.Header.Set("X-Real-IP", "37.187.141.25")
	req.SetBodyString(`Host: servicer.mgid.com
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36
Upgrade-Insecure-Requests: 1
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8
DNT: 1
Accept-Encoding: gzip, deflate
Accept-Language: ru-UA,ru;q=0.9,en-US;q=0.8,en;q=0.7,ru-RU;q=0.6
Cache-Control: no-cache`)
	var resp fasthttp.Response

	s:=time.Now()
	err := c.DoTimeout(&req, &resp, time.Second)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(s)
	log.Printf("%s took %s", "request", elapsed)

	if resp.StatusCode() == 200 { // OK
		useResponseBody(resp.Body())
	} else{
		println(resp.StatusCode())
		println(resp.Body())
	}
}

func useResponseBody(body []byte) {
	log.Println("resp.Body", string(body))
}