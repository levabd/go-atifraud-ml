package models

import (
	"time"
	"testing"
	"github.com/stretchr/testify/assert"

)

func TestTrimmedLogOrderData(t *testing.T) {
	ordered := Model{}.JsonStrToMap(`{
	"Random-header": 8,
	"From": 3,
	"Host": 0,
	"Accept": 5,
	"Connection": 1,
	"User-Agent": 2,
	"Accept-Encoding": 4}`)
	log := Log{
		Timestamp: time.Now(),
		OrderData: ordered,
	}

	trimmed_order_data := log.TrimOrderData()

	println("models.TestTrimmedLogValueData: trimmed_order_data - ", len(trimmed_order_data))
	println("models.TestTrimmedLogValueData: ordered - ", len(ordered))
	assert.Equal(t, 1, len(ordered)-len(trimmed_order_data), "1 order header must be cut")
	assert.NotEqual(t, len(trimmed_order_data), len(ordered), "trimmed_order_data must length be != ordered length")
}

func TestTrimmedLogValueData(t *testing.T) {
	headers := Model{}.JsonStrToMap(`{
	"Host":"www.popugaychik.com",
	"Connection":"Keep-alive",
	"Accept":"text\/html,application\/xhtml+xml,application\/xml;q=0.9,*\/*;q=0.8",
	"From":"googlebot(at)googlebot.com",
	"User-Agent":"Mozilla\/5.0 (compatible; Googlebot\/2.1; +http:\/\/www.google.com\/bot.html)",
	"Accept-Encoding":"gzip,deflate,br","If-Modified-Since":"Sat, 12 Aug 2017 08:16:35 GMT"}`)
	log := Log{
		Timestamp: time.Now(),
		ValueData: headers,
	}

	trimmed_headers := log.TrimValueData()
	println("models.TestTrimmedLogValueData: headers - ", len(headers))
	println("models.TestTrimmedLogValueData: trimmed_headers - ", len(trimmed_headers))
	assert.Equal(t, 5, len(headers)-len(trimmed_headers), "5 headers must be cut")
	assert.NotEqual(t, len(trimmed_headers), len(headers), "trimmed_headers must length be != headers length")
}
