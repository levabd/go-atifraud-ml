package services

import (
	h "github.com/levabd/go-atifraud-ml/lib/go/helpers"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"

	"testing"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/uniplaces/carbon"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"os"
)

func init() {

	h.LoadEnv()
}

func TestDataOrderingParsing(t *testing.T) {
	localAssert := assert.New(t)

	result, mainRow, valueRow, orderedRow := HandleLogLine(
		`1503090009,62.84.44.222,'
		{"Cache-Control":"no-cache",
		"Connection":"Keep-Alive",
		"Pragma":"no-cache",
		"Accept":"*\/*",
		"Accept-Encoding":"gzip, deflate",
		"From":"bingbot(at)microsoft.com",
		"Host":"www.vypekajem.com",
		"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"}'`,
		true,
		true,
		1503090000,
		1503090015)

	valueRowLen := len(valueRow)
	localAssert.Equal(true, result, "result is must be true")
	localAssert.Equal(9, len(mainRow), "main_row must contain 9 values")

	localAssert.Equal(0, orderedRow["Cache-Control"], "Cache-Control must be on 0 position")
	localAssert.Equal(1, orderedRow["Connection"], "Connection must be on 1 position")
	localAssert.Equal(2, orderedRow["Pragma"], "Pragma must be on 2 position")
	localAssert.Equal(3, orderedRow["Accept"], "Accept must be on 3 position")
	localAssert.Equal(4, orderedRow["Accept-Encoding"], "Accept-Encoding must be on 4 position")
	localAssert.Equal(5, orderedRow["From"], "From must be on 5 position")
	localAssert.Equal(6, orderedRow["Host"], "Host must be on 6 position")
	localAssert.Equal(7, orderedRow["User-Agent"], "User-Agent must be on 7 position")

	localAssert.Equal(9, len(mainRow), "main_row must contain 9 values")
	localAssert.Equal(len(orderedRow), valueRowLen, "value_row_len must equal to  len(ordered_row)")
}

func TestTimeCheckingParsing(t *testing.T) {
	localAssert := assert.New(t)

	result, _, _, _ := HandleLogLine(
		`1503080000,62.84.44.222,'
		{"Cache-Control":"no-cache",
		"Connection":"Keep-Alive",
		"Pragma":"no-cache",
		"Accept":"*\/*",
		"Accept-Encoding":"gzip, deflate",
		"From":"bingbot(at)microsoft.com",
		"Host":"www.vypekajem.com",
		"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"}'`,
		true,
		true,
		1503090000,
		1503090015)

	localAssert.Equal(false, result, "False must be given if line timestamp less than start_log_time")

	result, _, _, _ = HandleLogLine(
		`1503090019,62.84.44.222,'
		{"Cache-Control":"no-cache",
		"Connection":"Keep-Alive",
		"Pragma":"no-cache",
		"Accept":"*\/*",
		"Accept-Encoding":"gzip, deflate",
		"From":"bingbot(at)microsoft.com",
		"Host":"www.vypekajem.com",
		"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"}'`,
		true,
		true,
		1503090001,
		1503090015)

	localAssert.Equal(false, result, "False must be given if line timestamp greater than finish_log_time")

	result, _, _, _ = HandleLogLine(
		`1503090002,62.84.44.222,'
		{"Cache-Control":"no-cache",
		"Connection":"Keep-Alive",
		"Pragma":"no-cache",
		"Accept":"*\/*",
		"Accept-Encoding":"gzip, deflate",
		"From":"bingbot(at)microsoft.com",
		"Host":"www.vypekajem.com",
		"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"}'`,
		true,
		true,
		1503090000,
		1503090015)

	localAssert.Equal(true, result, "True must be given if line timestamp in neeeded period")
}

func TestOneFileParsing(t *testing.T) {
	localAssert := assert.New(t)

	start := time.Date(2016, 01, 17, 20, 34, 58, 651387237, time.UTC)
	finish := time.Date(2018, 03, 17, 20, 34, 58, 651387237, time.UTC)

	path := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "unit_tests_files", "2017-02-01.log")

	mainTable, valueTable, orderedTable := ParseSingleLog(path,
		true,
		true,
		start.Unix(),
		finish.Unix())
	localAssert.Equal(8, len(mainTable), "main_table be 7 in len")
	localAssert.Equal(8, len(valueTable), "value_table be 7 in len")
	localAssert.Equal(8, len(orderedTable), "ordered_table be 7 in len")
}

func TestGetUa(t *testing.T) {
	localAssert := assert.New(t)

	ua := GetUa(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36`)

	localAssert.Equal("chrome", ua["ua_family_code"], "ua_family_code must be chrome")
	localAssert.Equal("windows", ua["os_family_code"], "os_family_code must be windows")
	localAssert.Equal("browser", ua["ua_class_code"], "ua_class_code must be browser")
	localAssert.Equal("windows_10", ua["os_code"], "os_code must be windows_10")
	localAssert.Equal("desktop", ua["device_class_code"], "device_class_code must be desktop")
	localAssert.Equal("61.0.3163.100", ua["ua_version"], "ua_version must be 61.0.3163.100")
}

func TestFindBotByUaFamily(t *testing.T) {
	localAssert := assert.New(t)

	isCrawler := IsCrawler("12", `Mozilla/5.0 (compatible; Googlebot/2.1; startmebot/1.0; +https://start.me/bot)`)
	localAssert.Equal(true, isCrawler, "is_crawler must be true")
}

func TestParsingLogic(t *testing.T) {
	localAssert := assert.New(t)

	result, _, _, _ := HandleLogLine(
		`1503090009,40.77.167.95,'
		{"Cache-Control":"no-cache",
		"Connection":"Keep-Alive",
		"Pragma":"no-cache",
		"Accept":"*\/*",
		"Accept-Encoding":"gzip, deflate",
		"From":"bingbot(at)microsoft.com",
		"Host":"www.vypekajem.com",
		"User-Agent":"Mozilla/5.0 (compatible; Googlebot/2.1; startmebot/1.0; +https://start.me/bot)"}'`,
		true,
		true,
		1503090000,
		1503090015)

	localAssert.Equal(false, result, "result must be false")

	// disable filter crawler
	// disable ua parsing
	result, main, _, _ := HandleLogLine(
		`1503090009,62.84.44.222,'
		{"Cache-Control":"no-cache",
		"Connection":"Keep-Alive",
		"Pragma":"no-cache",
		"Accept":"*\/*",
		"Accept-Encoding":"gzip, deflate",
		"From":"bingbot(at)microsoft.com",
		"Host":"www.vypekajem.com",
		"User-Agent":"Mozilla/5.0 (compatible; Googlebot/2.1; startmebot/1.0; +https://start.me/bot)"}'`,
		false,
		false,
		1503090000,
		1503090015)

	localAssert.Equal(true, result, "result must be true")
	localAssert.Equal(3, len(main), "main data must be 3 in length")

	// disable filter crawler
	// enable ua parsing
	r, main, value, order := HandleLogLine(
		`1503090009,62.84.44.222,'
		{"Cache-Control":"no-cache",
		"Connection":"Keep-Alive",
		"Pragma":"no-cache",
		"Accept":"*\/*",
		"Accept-Encoding":"gzip, deflate",
		"From":"bingbot(at)microsoft.com",
		"Host":"www.vypekajem.com",
		"User-Agent":"Mozilla/5.0 (compatible; Googlebot/2.1; startmebot/1.0; +https://start.me/bot)"}'`,
		false,
		true,
		1503090000,
		1503090015)

	localAssert.Equal(true, r, "result must be true")
	localAssert.Equal(9, len(main), "main data must be 9 in length")
	localAssert.Equal(8, len(value), "value_data must be 8 in length")
	localAssert.Equal(8, len(order), "order_data must be 8 in length")
}

func TestParseLatestLogGzFile(t *testing.T) {
	m.TruncateTable("gz_logs")
	localAssert := assert.New(t)

	start := time.Date(2016, 01, 17, 20, 34, 58, 651387237, time.UTC)
	finish := time.Date(2018, 03, 17, 20, 34, 58, 651387237, time.UTC)

	logsDir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "unit_tests_files")
	files := h.GetFileFromDirWithExt(logsDir, "gz")

	err := ParseAndStoreSingleGzLogInDb(filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "unit_tests_files", files[len(files)-1]), true,
		true,
		start.Unix(),
		finish.Unix(),
		true)

	localAssert.Equal(nil, err, "error must be nil")
}

func TestGetLatestLogFilePath(t *testing.T) {

	m.TruncateTable("gz_logs")

	logsDir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "unit_tests_files")
	files := h.GetFileFromDirWithExt(logsDir, "gz")

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	// fill table with existed files
	tx := db.Begin()
	for _, fileName := range files {
		gzLog := m.GzLog{FileName: fileName}
		tx.Create(&gzLog)
	}
	tx.Commit()

	// must return error - because file in DB
	_, _, err = GetLatestLogFilePath()

	// clear table
	assert.NotEqual(t, nil, err, "Error must not be nil")
}

func TestGetLogsInPeriod(t *testing.T) {
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
		log := m.Log{Timestamp: timeStamp,}
		tx.Create(&log)
	}
	tx.Commit()

	logs := GetLogsInPeriod(
		carbon.Now().SubMonths(2).Unix(),
		carbon.Now().Unix())

	assert.Equal(t, 3, len(logs), "Error must not be nil")
	m.TruncateTable("logs")
}

func TestGetTrimmedLodMapsForPeriod(t *testing.T) {
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
		value :=m.Model{}.JsonStrToMap(`{"Host":"www.popugaychik.com","Connection":"Keep-alive","Accept":"text\/html,application\/xhtml+xml,application\/xml;q=0.9,*\/*;q=0.8","From":"googlebot(at)googlebot.com","User-Agent":"Mozilla\/5.0 (compatible; Googlebot\/2.1; +http:\/\/www.google.com\/bot.html)","Accept-Encoding":"gzip,deflate,br","If-Modified-Since":"Sat, 12 Aug 2017 08:16:35 GMT"}`)
		ordered :=m.Model{}.JsonStrToMap(`{"Random_header": 8,"From": 3, "Host": 0, "Accept": 5, "Connection": 1, "User-Agent": 2, "Accept-Encoding": 4}`)
		log := m.Log{Timestamp: timeStamp, ValueData: value, OrderData:ordered}
		tx.Create(&log)
	}
	tx.Commit()

	_, trimmedValueData, trimmedOrderData :=GetTrimmedLodMapsForPeriod(
		carbon.Now().SubMonths(2).Unix(),
		carbon.Now().SubMonths(1).Unix())

	assert.Equal(t, 2, len(trimmedValueData[0]),"log headeer value must be 2 in length after trimming")
	assert.Equal(t, 6, len(trimmedOrderData[0]),"log headeer value must be 6 in length after trimming")
	m.TruncateTable("logs")
}

func TestPrepareData(t *testing.T) {

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
		value :=m.Model{}.JsonStrToMap(`{"Host":"www.popugaychik.com","Connection":"Keep-alive","Accept":"text\/html,application\/xhtml+xml,application\/xml;q=0.9,*\/*;q=0.8","From":"googlebot(at)googlebot.com","User-Agent":"Mozilla\/5.0 (compatible; Googlebot\/2.1; +http:\/\/www.google.com\/bot.html)","Accept-Encoding":"gzip,deflate,br","If-Modified-Since":"Sat, 12 Aug 2017 08:16:35 GMT"}`)
		ordered :=m.Model{}.JsonStrToMap(`{"Random_header": 8,"From": 3, "Host": 0, "Accept": 5, "Connection": 1, "User-Agent": 2, "Accept-Encoding": 4}`)
		log := m.Log{Timestamp: timeStamp, ValueData: value, OrderData:ordered}
		tx.Create(&log)
	}
	tx.Commit()

	PrepareData(carbon.Now().SubMonths(2).Unix(), carbon.Now().SubMonths(1).Unix())
	m.TruncateTable("logs")

}

