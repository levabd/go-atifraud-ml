package main

import (

	"github.com/levabd/go-atifraud-ml/lib/go/parsers"
	"fmt"
	"time"
	"path/filepath"
	"os"
	"github.com/joho/godotenv"
)

func main() {




}

func checkOneFileParsing()  {
	start := time.Date(
		2016, 01, 17, 20, 34, 58, 651387237, time.UTC)
	finish := time.Date(
		2018, 03, 17, 20, 34, 58, 651387237, time.UTC)


	current_dir, err := filepath.Abs("./")
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}

	path := filepath.Join(current_dir, "data", "unit_tests_files", "2017-02-01.log")

	main_table, value_table, ordered_table:= parsers.ParseSingleLog(path,
		true,
		true,
		start.Unix(),
		finish.Unix())

	fmt.Println("result: ", len(main_table), len(value_table), len(ordered_table))
}

func checkOneLineParsing()  {
	result, main_row, value_row, ordered_row :=parsers.HandleLogLine(
		//`1506501921,95.181.252.91,'{"Host":"www.vypekajem.com","Connection":"keep-alive","Upgrade-Insecure-Requests":"1","User-Agent":"Mozilla\/5.0 (Windows NT 6.1) AppleWebKit\/537.36 (KHTML, like Gecko) Chrome\/55.0.2883.87 Safari\/537.36","Accept":"text\/html,application\/xhtml+xml,application\/xml;q=0.9,image\/webp,*\/*;q=0.8","Referer":"https:\/\/www.google.ru\/","Accept-Encoding":"gzip, deflate, sdch","Accept-Language":"ru-RU,ru;q=0.8,en-US;q=0.6,en;q=0.4"}'`,
		`1503090009,40.77.167.95,'{"Cache-Control":"no-cache","Connection":"Keep-Alive","Pragma":"no-cache","Accept":"*\/*","Accept-Encoding":"gzip, deflate","From":"bingbot(at)microsoft.com","Host":"www.vypekajem.com","User-Agent":"Mozilla\/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit\/537.51.1 (KHTML, like Gecko) Version\/7.0 Mobile\/11A465 Safari\/9537.53 (compatible; bingbot\/2.0; +http:\/\/www.bing.com\/bingbot.htm)"}'`,
		true,
		true,
		1503090000,
		1503090015)
	fmt.Println("result: ", result, main_row, value_row, ordered_row)
}

func init() {
	godotenv.Load("lib/go/parsers/.env")
}