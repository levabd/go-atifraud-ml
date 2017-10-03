package main

import (
	"time"
	"fmt"
	"log"
)
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
func main() {
	defer timeTrack(time.Now(), "maon")


	start_minute:= time.Now().Minute()
	start_second:= time.Now().Second()
	start_millisecond:= start_second/1000
	fmt.Println(fmt.Printf("%f",start_millisecond))

	start_microsecond:= start_millisecond/1000
	start_nanoseconds:= start_microsecond/1000
	println(fmt.Sprintf("start time  %d, %d, %f, %f, %f",
		start_minute,
		start_second,
		start_millisecond,
		start_microsecond, start_nanoseconds))

	//start_time := time.Now()
	//start_nanosecond := start_time.Nanosecond()
	//println(start_time.Minute(), start_time.Second(), start_nanosecond)
	//
	//services.PrepareData(carbon.Now().SubMonths(2).Unix(), carbon.Now().Unix())
	//
	//end := time.Now()
	//println(end.Minute(), end.Second(), end.Nanosecond(), end.Nanosecond()-start_nanosecond)
}
