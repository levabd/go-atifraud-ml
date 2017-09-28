package parsers

import (
	"testing"
	"github.com/joho/godotenv"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"time"
	"path/filepath"
	"os"
)

func TestDataOrderingParsing(t *testing.T)  {
	result, main_row, value_row, ordered_row :=HandleLogLine(
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

	value_row_len := len(value_row)
	helpers.AssertEqual(t, result, true, "result is must be true")
	helpers.AssertEqual(t, len(main_row), 9, "main_row must contain 9 values")

	helpers.AssertEqual(t, ordered_row["Cache-Control"], 0, "Cache-Control must be on 0 position")
	helpers.AssertEqual(t, ordered_row["Connection"], 1, "Connection must be on 1 position")
	helpers.AssertEqual(t, ordered_row["Pragma"], 2, "Pragma must be on 2 position")
	helpers.AssertEqual(t, ordered_row["Accept"], 3, "Accept must be on 3 position")
	helpers.AssertEqual(t, ordered_row["Accept-Encoding"], 4, "Accept-Encoding must be on 4 position")
	helpers.AssertEqual(t, ordered_row["From"], 5, "From must be on 5 position")
	helpers.AssertEqual(t, ordered_row["Host"], 6, "Host must be on 6 position")
	helpers.AssertEqual(t, ordered_row["User-Agent"], 7, "User-Agent must be on 7 position")

	helpers.AssertEqual(t, len(main_row), 9, "main_row must contain 9 values")
	helpers.AssertEqual(t, value_row_len, len(ordered_row), "value_row_len must equal to  len(ordered_row)")
}

func TestTimeCheckingParsing(t *testing.T)  {
	result, _, _, _ :=HandleLogLine(
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


	helpers.AssertEqual(t,result, false, "False must be given if line timestamp less than start_log_time")

	result, _, _, _ =HandleLogLine(
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

	helpers.AssertEqual(t, result, false, "False must be given if line timestamp greater than finish_log_time")

	result, _, _, _ =HandleLogLine(
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

	helpers.AssertEqual(t, result, true, "True must be given if line timestamp in neeeded period")
}

func TestOneFileParsing(t *testing.T)  {

	start := time.Date(2016, 01, 17, 20, 34, 58, 651387237, time.UTC)
	finish := time.Date(2018, 03, 17, 20, 34, 58, 651387237, time.UTC)

	path := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data","unit_tests_files", "2017-02-01.log")

	main_table, value_table, ordered_table:= ParseSingleLog(path,
		true,
		true,
		start.Unix(),
		finish.Unix())

	helpers.AssertEqual(t,  len(main_table), 16, "main_table be 16 in len")
	helpers.AssertEqual(t,  len(value_table), 16, "value_table be 16 in len")
	helpers.AssertEqual(t,  len(ordered_table), 16, "ordered_table be 16 in len")
}




func TestGetUa(t *testing.T){
	ua:=GetUa(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36`)

	helpers.AssertEqual(t,  ua["ua_family_code"], "chrome", "ua_family_code must be chrome")
	helpers.AssertEqual(t,  ua["os_family_code"], "windows", "os_family_code must be windows")
	helpers.AssertEqual(t,  ua["ua_class_code"], "browser", "ua_class_code must be browser")
	helpers.AssertEqual(t,  ua["os_code"], "windows_10", "os_code must be windows_10")
	helpers.AssertEqual(t,  ua["device_class_code"], "desktop", "device_class_code must be desktop")
	helpers.AssertEqual(t,  ua["ua_version"], "61.0.3163.100", "ua_version must be 61.0.3163.100")
}


func TestFindBotByUaFamily(t *testing.T){
	is_crawler:=IsCrawlerByUdger( "12", `Mozilla/5.0 (compatible; Googlebot/2.1; startmebot/1.0; +https://start.me/bot)`)
	helpers.AssertEqual(t,  is_crawler, true, "is_crawler must be true")
}


func TestParsingLogic(t *testing.T){
		result, _, _, _ :=HandleLogLine(
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

	helpers.AssertEqual(t,  result, false, "result must be false")

	// disable filter crawler
	// disable ua parsing
	result, m, _, _ :=HandleLogLine(
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

	helpers.AssertEqual(t, result, true, "result must be true")
	helpers.AssertEqual(t, len(m), 3, "main data must be 3 in length")

	 // disable filter crawler
	// enable ua parsing
	r, m, v, o :=HandleLogLine(
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

	helpers.AssertEqual(t,  r, true, "result must be true")
	helpers.AssertEqual(t,   len(m),  9, "main data must be 9 in length")
	helpers.AssertEqual(t,   len(v),  8, "value_data must be 8 in length")
	helpers.AssertEqual(t,   len(o),  8, "order_data must be 8 in length")
}

func TestParseLatestLogGzFile(t *testing.T) {

	start := time.Date(2016, 01, 17, 20, 34, 58, 651387237, time.UTC)
	finish := time.Date(2018, 03, 17, 20, 34, 58, 651387237, time.UTC)

	main_table, value_table, ordered_table:=ParseLatestLogGzFile( true,
		true,
		start.Unix(),
		finish.Unix(),
		true)

	helpers.AssertEqual(t,  len(main_table), 16, "main_table be 16 in len")
	helpers.AssertEqual(t,  len(value_table), 16, "value_table be 16 in len")
	helpers.AssertEqual(t,  len(ordered_table), 16, "ordered_table be 16 in len")
}

func init() {
	godotenv.Load()
}