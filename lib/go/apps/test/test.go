package main

import (
	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"github.com/uniplaces/carbon"
	"time"
)

func main() {

	start_time := time.Now()
	start_nanosecond := start_time.Nanosecond()
	println(start_time.Minute(), start_time.Second(), start_nanosecond)

	services.PrepareData(carbon.Now().SubMonths(2).Unix(), carbon.Now().Unix())

	end := time.Now()
	println(end.Minute(), end.Second(), end.Nanosecond(), end.Nanosecond()-start_nanosecond)
}
