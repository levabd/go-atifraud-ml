package services

import (
	"github.com/jinzhu/gorm"
	"time"
	"github.com/uniplaces/carbon"
	"fmt"
	"testing"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPairsDictList(t *testing.T) {

	m.TruncateTable("logs")
	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}
	// fill table
	tx := db.Begin()
	for i := 0; i < 10; i++ {
		timeStamp := time.Unix(carbon.Now().SubMonths(i).Unix(), 0)
		value := m.Model{}.JsonStrToMap(fmt.Sprintf(`{"Host":"www.popugaychik.com%v","Connection":"Keep-alive","Accept":"text\/html,application\/xhtml+xml,application\/xml;q=0.9,*\/*;q=0.8","From":"googlebot(at)googlebot.com","User-Agent":"Mozilla\/5.0 (compatible; Googlebot\/2.1; +http:\/\/www.google.com\/bot.html)","Accept-Encoding":"gzip,deflate,br","If-Modified-Since":"Sat, 12 Aug 2017 08:16:35 GMT"}`, i))
		ordered := m.Model{}.JsonStrToMap(fmt.Sprintf(`{"Random_header": 8,"From": %v, "Host":  %v, "Accept": %v, "Connection": 1, "User-Agent": 2, "Accept-Encoding": 4}`, 3+i, 0+i, 5+i))

		fmt.Println()

		log := m.Log{Timestamp: timeStamp, ValueData: value, OrderData: ordered}
		tx.Create(&log)
	}
	tx.Commit()

	_, _, trimmedOrderData := GetTrimmedLodMapsForPeriod(
		carbon.Now().SubMonths(2).Unix(),
		carbon.Now().SubMonths(1).Unix())

	orderFeatures := GetOrderFeatures(trimmedOrderData)

	assert.Equal(t, 2, len(orderFeatures))
	assert.Equal(t, 15, len(orderFeatures[0]))
	m.TruncateTable("logs")
}

func TestDefineKeyFoPair(t *testing.T)  {

	var twoHeaders []map[string]interface{}

	first :=make(map[string]interface{})
	first["header1"] = "2"

	second :=make(map[string]interface{})
	second["header2"] = "5"

	twoHeaders = append(twoHeaders, first)
	twoHeaders = append(twoHeaders, second)

	assert.Equal(t, "header1 < header2", DefineKeyFoPair(twoHeaders))
}