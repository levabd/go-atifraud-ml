package main

import (
	"time"
	"fmt"
	"log"
	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"github.com/uniplaces/carbon"
)
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
func main() {
	defer timeTrack(time.Now(), "maon")

	start_time := time.Now()
	start_nanosecond := start_time.Nanosecond()
	println(start_time.Minute(), start_time.Second(), start_nanosecond)

	services.PrepareData(carbon.Now().SubMonths(2).Unix(), carbon.Now().Unix())

	end := time.Now()
	println(end.Minute(), end.Second(), end.Nanosecond(), end.Nanosecond()-start_nanosecond)
}
