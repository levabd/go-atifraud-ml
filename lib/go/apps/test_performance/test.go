package main

import (
	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"github.com/uniplaces/carbon"
	"time"
	"github.com/valyala/fasthttp"
	"log"
	"fmt"
	"github.com/levabd/go-atifraud-ml/lib/go/models"
)

var (
	crawlers, bots, humans                int
	crawlers_time, bots_time, humans_time []time.Duration

	c    = &fasthttp.Client{}
	req  fasthttp.Request
	resp fasthttp.Response
)

func main() {

	var logs []models.Log

	start, _ := carbon.Create(2017, 9, 1, 0, 0, 0, 0, "Asia/Almaty")
	end, _ := carbon.Create(2017, 9, 11, 13, 0, 0, 0, "Asia/Almaty")

	logs = services.GetLogsInPeriod(start.Unix(), end.Unix())[:10000]

	for i := 0; i < len(logs); i++ {
		d := logs[i].ValueData
		var str string
		for key, item := range d {
			str += key + ":" + item.(string)+ "\n"
		}
		str+="user-agent:"+logs[i].UserAgent+ "\n"
		//valueString, _ := json.Marshal(d)
		makeRequest(str)
	}

	defer timeTrack(time.Now(), "main")
	defer resp.ConnectionClose()

	println(fmt.Sprintf("humans %v, crawlers %v, bots %v ", humans, crawlers, bots))
}

func makeRequest(body string) string {

	start := time.Now()

	req.Header.SetMethod("POST")
	req.SetRequestURI("http://localhost:8082")
	req.Header.Set("Host", "localhost")
	req.SetBodyString(body)

	err := c.Do(&req, &resp)
	if err != nil {
		println("error", err)
	}

	if string(resp.Body()) == "human" {
		humans_time = append(humans_time, timeTrack(start, "human"))
		humans++
	} else if string(resp.Body()) == "crawler" {
		crawlers_time = append(crawlers_time, timeTrack(start, "crawler"))
		crawlers++
	} else if string(resp.Body()) == "bot" {
		bots_time = append(bots_time, timeTrack(start, "bot"))
		bots++
	}

	return string(resp.Body())
}

func timeTrack(start time.Time, name string) (elapsed time.Duration) {
	log.Println(fmt.Sprintf("%s took %s", name, time.Since(start)))
	return time.Since(start)
}
