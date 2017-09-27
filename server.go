package main

import (
	"fmt"
	"github.com/levabd/go-atifraud-ml/lib/parsers"
)

func main() {

	//result, main_row, value_row, ordered_row :=parsers.HandleLogLine(
	//	`1506501921,95.181.252.91,'{"Host":"www.vypekajem.com","Connection":"keep-alive","Upgrade-Insecure-Requests":"1","User-Agent":"Mozilla\/5.0 (Windows NT 6.1) AppleWebKit\/537.36 (KHTML, like Gecko) Chrome\/55.0.2883.87 Safari\/537.36","Accept":"text\/html,application\/xhtml+xml,application\/xml;q=0.9,image\/webp,*\/*;q=0.8","Referer":"https:\/\/www.google.ru\/","Accept-Encoding":"gzip, deflate, sdch","Accept-Language":"ru-RU,ru;q=0.8,en-US;q=0.6,en;q=0.4"}'`,
	//	true,
	//	true,
	//	1506501851,
	//	1506501951)

	//fmt.Println("result: ", result, main_row, value_row, ordered_row )

	parsers.ParseSingleLog('logs')
}
