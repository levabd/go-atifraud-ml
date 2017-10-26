package main

import  (
	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"github.com/uniplaces/carbon"
	"time"
	"github.com/valyala/fasthttp"
	"log"
	"fmt"
	"github.com/levabd/go-atifraud-ml/lib/go/models"
	"encoding/json"
)

var (
	crawlers, bots, humans                int
	crawlers_time, bots_time, humans_time []time.Duration

	c    = &fasthttp.Client{}
	req  fasthttp.Request
	resp fasthttp.Response
)

func main() {

	var (
		logs []models.Log
	)

	// single latest log gz file parsing
	//fileName, err := services.GetLatestLogFile()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//path := filepath.Join(filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs"), fileName)
	//u, err := udger.New(os.Getenv("DB_FILE_PATH_UDGER"))
	//if err != nil {
	//	log.Fatalf("server.go - init: Failed to initiate udger helper: %s", err)
	//}

	//logs = services.ParseAndGetDataFromSingleGzFile(
	//	path,
	//	carbon.Now().SubMonths(1).Unix(),
	//	carbon.Now().Unix(),
	//	u)
	start, _ := carbon.Create(2017, 9,1,0,0,0,0,"Asia/Almaty")
	end, _ := carbon.Create(2017, 9,11,13,0,0,0,"Asia/Almaty")

	logs = services.GetLogsInPeriod(start.Unix(), end.Unix())[:10000]

	//if err != nil {
	//	panic(err)
	//	return
	//}

	for i := 0; i < len(logs); i++ {
		d := logs[i].ValueData
		valueString, _:= json.Marshal(d)
		//println(string(valueString))
		makeRequest(string(valueString))
	}

	defer timeTrack(time.Now(), "main")
	defer resp.ConnectionClose()
	println(fmt.Sprintf("humans %v, crawlers %v, bots %v ", humans, crawlers, bots))
	//println(fmt.Sprintf("time   ,    time %s,     time %s", countAverage(humans_time), countAverage(crawlers_time), countAverage(bots_time)))
}

func countAverage(arr []time.Duration) time.Duration{
	var sum int64

	for _, t := range arr {
		sum += int64(t)
	}

	return time.Duration(sum/int64(len(arr)))
}

func timeTrack(start time.Time, name string) (elapsed time.Duration) {
	log.Println(fmt.Sprintf("%s took %s", name, time.Since(start)))

	return time.Since(start)
}

func makeRequest(body string) string {

	start := time.Now()
	//defer timeTrack(time.Now(), "main for "+ip)
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
