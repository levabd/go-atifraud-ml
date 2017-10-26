package main

import (
	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"github.com/uniplaces/carbon"
	"math/rand"
	"github.com/levabd/go-atifraud-ml/lib/go/models"
	"github.com/jinzhu/gorm"
)

var uaFamilyCodeList []string

func getRandom(logUaFamilyCode string) string {

	if len(uaFamilyCodeList) == 1 {
		return uaFamilyCodeList[0]
	}

	randomUaFamilyCode := uaFamilyCodeList[rand.Intn(len(uaFamilyCodeList))]

	if logUaFamilyCode == randomUaFamilyCode {
		return getRandom(logUaFamilyCode)
	}

	// remove item
	extractedElementIndex := getIndex(uaFamilyCodeList, randomUaFamilyCode)
	uaFamilyCodeList = append(uaFamilyCodeList [:extractedElementIndex], uaFamilyCodeList [extractedElementIndex+1:]...)

	return randomUaFamilyCode
}

func getIndex(array []string, value string) int {
	for p, v := range array {
		if (v == value) {
			return p
		}
	}
	return -1
}

func main() {
	testNumberOfRows := 30000

	start, _ := carbon.Create(2017, 9,7,23,0,0,0,"Asia/Almaty")
	end, _ := carbon.Create(2017, 9, 10, 23, 0, 0, 0, "Asia/Almaty")

	logs := services.GetLogsInPeriod(start.Unix(), end.Unix())[:testNumberOfRows ]

	println("len(logs)", len(logs))
	println("len(testSelect)", len(logs))

	uaFamilyCodeSum := make(map[string]int)
	for _, log := range logs {
		if _, ok := uaFamilyCodeSum[log.UaFamilyCode]; ok {
			uaFamilyCodeSum[log.UaFamilyCode]++
		} else {
			uaFamilyCodeSum[log.UaFamilyCode] = 1
		}
	}

	uaFamilyCodeList = make([]string, 0)
	for key, _ := range uaFamilyCodeSum {
		uaFamilyCodeList = append(uaFamilyCodeList, key)
	}

	uaFamilyCodeReplacementList := make(map[string]string)

	testLogs := make([]models.Log, len(logs))

	for index, log := range logs {

		testLogs[index] = log

		if value, ok := uaFamilyCodeReplacementList[testLogs[index].UaFamilyCode]; ok {
			testLogs[index].UaFamilyCode = value
			continue
		}

		randomUaFamilyCode := getRandom(testLogs[index].UaFamilyCode)

		uaFamilyCodeReplacementList[testLogs[index].UaFamilyCode] = randomUaFamilyCode

		testLogs[index].UaFamilyCode = randomUaFamilyCode
	}
	println("uaFamilyCodeReplacementList len", len(uaFamilyCodeReplacementList))
	println("uaFamilyCodeList len", len(uaFamilyCodeList))
	println("uaFamilyCodeList", uaFamilyCodeList[0])

	db, err := gorm.Open("postgres", models.GetDBConnectionStr())
	if err != nil {
		panic("user_agent_helpers.go - LoadFittedUserAgentCodes: Failed to connect to database")
	}
	defer db.Close()
	if !db.HasTable(&models.TestLog{}) {
		db.AutoMigrate(&models.TestLog{})
	}
	db.Exec("TRUNCATE TABLE test_logs;")

	// store
	tx := db.Begin()

	for i := 0; i < len(testLogs); i++ {
		cacheFeatures := models.TestLog{
			OriginalID:           logs[i].ID,
			OriginalUaFamilyCode: logs[i].UaFamilyCode,
			CreatedAt:            testLogs[i].CreatedAt,
			Timestamp:            testLogs[i].Timestamp,
			Ip:                   testLogs[i].Ip,
			UserAgent:            testLogs[i].UserAgent,
			UaFamilyCode:         testLogs[i].UaFamilyCode,
			UaVersion:            testLogs[i].UaVersion,
			UaClassCode:          testLogs[i].UaClassCode,
			DeviceClassCode:      testLogs[i].DeviceClassCode,
			OsFamilyCode:         testLogs[i].OsFamilyCode,
			OsCode:               testLogs[i].OsCode,
			ValueData:            testLogs[i].ValueData,
			OrderData:            testLogs[i].OrderData,
		}
		tx.Create(&cacheFeatures)

		//fmt.Println(fmt.Sprintf("human ua/bot: %s = %s", logs[i].UaFamilyCode, testLogs[i].UaFamilyCode))
	}

	tx.Commit()

}
