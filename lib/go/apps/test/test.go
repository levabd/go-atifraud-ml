package main

import (
	"time"
	"log"
	"fmt"

	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"github.com/uniplaces/carbon"
)
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
func main() {
	defer timeTrack(time.Now(), "maon")

	startTime := time.Now()
	startNanosecond := startTime.Nanosecond()
	println(startTime.Minute(), startTime.Second(), startNanosecond)

	userAgent, valueFeatures, ordersFeatures := services.PrepareData(carbon.Now().SubMonths(2).Unix(), carbon.Now().Unix())

	println(len(userAgent), len(valueFeatures), len(ordersFeatures))

	fmt.Println(len(valueFeatures[0]))
	fmt.Println(len(ordersFeatures[0]))

	end := time.Now()
	println(end.Minute(), end.Second(), end.Nanosecond(), end.Nanosecond() - startNanosecond)
}
