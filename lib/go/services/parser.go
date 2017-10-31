package services

import (
	"regexp"
	"strings"
	"github.com/buger/jsonparser"
	"fmt"
	"path/filepath"
	"os"
	"time"
	h "github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"bufio"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
	"github.com/jinzhu/gorm"
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"github.com/levabd/go-atifraud-ml/lib/go/udger"
)

func timeIsWrong(logTimestamp string, startLogTime int64, finishLogTime int64, ) bool {

	timeInLogs := time.Unix(0, h.StrToInt64(logTimestamp))

	start := time.Unix(0, startLogTime)
	finish := time.Unix(0, finishLogTime)

	return timeInLogs.Before(start) || timeInLogs.After(finish)
}

func HandleLogLine(
	line string,
	filterCrawlers bool,
	needUaParsing bool,
	startLogTime int64,
	finishLogTime int64,
	udger *udger.Udger) (bool, map[string]interface{}, map[string]interface{}, map[string]interface{}) {

	var (
		elements                          = strings.SplitN(line, ",", 3)
		mainRow    map[string]interface{} = make(map[string]interface{})
		valueRow   map[string]interface{} = make(map[string]interface{})
		orderedRow map[string]interface{} = make(map[string]interface{})
	)

	// get data
	timestamp := line[:10]
	i := strings.Index(line[11:], ",")
	ip := line[11:11+i]
	headers := line[11+i+1:]

	//println("timestamp", timestamp)
	//println("ip", ip)
	//println("headers: ", headers)

	if timestamp == "" || ip == "" {
		log.Printf("parsers.HandleLogLine: one is empty(timestamp, ip, headers)  %s, %s, %s", timestamp, ip, headers)
		return false, nil, nil, nil
	}

	if timeIsWrong(timestamp, startLogTime, finishLogTime) {
		log.Printf("parsers.HandleLogLine: timeIsWrong: log_time - %s, start_log_time - %v, finish_log_time %v",
			string(elements[0]), startLogTime, finishLogTime)

		return false, nil, nil, nil
	}

	if len(elements) < 3 {
		log.Printf("parsers.HandleLogLine: not enoth elements %+v", elements)
		return false, nil, nil, nil
	}

	if filterCrawlers && (timestamp == "" || ip == "" || headers == "") {
		log.Printf("parsers.HandleLogLine: elements[0] == empty || elements[1] == empty || elements[2] == empty %+v", elements)
		return false, nil, nil, nil
	}

	if headers != "" {

		if headers[:1] != "'" {
			log.Printf("parsers.HandleLogLine: string(elements[2][0]) != ' elements[2] = %s", headers)
			return false, nil, nil, nil
		}

		// replace all \/ by single /
		var re = regexp.MustCompile(`(?m)\\/`)
		var substitution = `/`

		str_replaced := re.ReplaceAllString(headers, substitution)
		jsonToParse := strings.TrimPrefix(strings.TrimSuffix(str_replaced, "'"), "'")

		data := []byte(jsonToParse)
		i := 0
		jsonparser.ObjectEach(data, func(
			key []byte,
			value []byte,
			dataType jsonparser.ValueType,
			offset int) error {
			orderedRow[ string(key)] = i
			valueRow[ string(key)] = string(value)
			i = i + 1
			return nil
		})

		if valueRow["User-Agent"] == nil && valueRow["user-agent"] == nil && valueRow["user_agent"] == nil {
			log.Printf("parsers.ParseAndStoreSingleGzLogInDb: No user agent in headers %+v", valueRow)
			return false, nil, nil, nil
		}

		ua := ""

		if _ua, ok := valueRow["User-Agent"].(string); ok {
			ua = _ua
		}

		if _ua, ok := valueRow["user-agent"].(string); ok && ua != "" {
			ua = _ua
		}

		if _ua, ok := valueRow["user_agent"].(string); ok && ua != ""{
			ua = _ua
		}

		if ua == "" {
			return false, nil, nil, nil
		}

		if filterCrawlers && udger.IsCrawler(ip, ua) {
			//log.Println("ua is crawler: ", elements[1], ua)
			return false, nil, nil, nil
		}

		if needUaParsing {
			uaObj := udger.ParseData["user_agent"]

			if uaObj["ua_family_code"] == "" {
				log.Println("ua_family_code is empty: ", elements[1], ua)
				return false, nil, nil, nil
			}

			mainRow["ua_family_code"] = uaObj["ua_family_code"]
			mainRow["ua_version"] = uaObj["ua_family_code"] + " " + uaObj["ua_version"]
			mainRow["ua_class_code"] = uaObj["ua_class_code"]
			mainRow["device_class_code"] = uaObj["device_class_code"]
			mainRow["os_family_code"] = uaObj["os_family_code"]
			mainRow["os_code"] = uaObj["os_code"]
		}

		mainRow["timestamp"] = timestamp
		mainRow["ip"] = ip
		mainRow["User_Agent"] = valueRow["User-Agent"].(string)

		return true, mainRow, valueRow, orderedRow
	}

	return false, mainRow, valueRow, orderedRow
}

func ParseLogLineWithCrawlers(
	line string,
	startLogTime int64,
	finishLogTime int64,
) (headers_return string) {

	var (
		elements = strings.SplitN(line, ",", 3)
	)

	// get data
	timestamp := line[:10]
	i := strings.Index(line[11:], ",")
	ip := line[11:11+i]
	headers := line[11+i+1:]

	if timestamp == "" || ip == "" {
		log.Printf("parsers.HandleLogLine: one is empty(timestamp, ip, headers)  %s, %s, %s", timestamp, ip, headers)
		return headers_return
	}

	if timeIsWrong(timestamp, startLogTime, finishLogTime) {
		log.Printf("parsers.HandleLogLine: timeIsWrong: log_time - %s, start_log_time - %v, finish_log_time %v",
			string(elements[0]), startLogTime, finishLogTime)

		return headers_return
	}

	if len(elements) < 3 {
		log.Printf("parsers.HandleLogLine: not enoth elements %+v", elements)
		return headers_return
	}

	if timestamp == "" || ip == "" || headers == "" {
		log.Printf("parsers.HandleLogLine: elements[0] == empty || elements[1] == empty || elements[2] == empty %+v", elements)
		return headers_return
	}

	if headers != "" {
		if headers[:1] != "'" {
			log.Printf("parsers.HandleLogLine: string(elements[2][0]) != ' elements[2] = %s", headers)
			return headers_return
		}

		// replace all \/ by single /
		var re = regexp.MustCompile(`(?m)\\/`)
		var substitution = `/`

		str_replaced := re.ReplaceAllString(headers, substitution)
		headers_return = strings.TrimPrefix(strings.TrimSuffix(str_replaced, "'"), "'")

		return headers_return
	}

	return headers_return
}

func ParseSingleLog(
	pathToLog string,
	filterCrawlers bool,
	parseUa bool,
	startLogTime int64,
	finishLogTime int64,
	udger *udger.Udger) ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {

	inFile, _ := os.Open(pathToLog)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var (
		mainTable    []map[string]interface{}
		valueTable   []map[string]interface{}
		orderedTable []map[string]interface{}
	)

	for scanner.Scan() {
		isHuman, mainRow, valueRow, orderedRow := HandleLogLine(
			scanner.Text(),
			filterCrawlers,
			parseUa,
			startLogTime,
			finishLogTime, udger)

		if isHuman {
			mainTable = append(mainTable, mainRow)
			valueTable = append(valueTable, valueRow)
			orderedTable = append(orderedTable, orderedRow)
		}
	}

	return mainTable, valueTable, orderedTable
}

func ParseAndStoreSingleGzLogInDb(
	filePath string,
	filterCrawlers bool,
	parseUa bool,
	startLogTime int64,
	finishLogTime int64,
	doRoolback bool, udger *udger.Udger) error {

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		log.Printf("parsers.ParseAndStoreSingleGzLogInDb: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}
	bytesOfString, _ := h.ReadGzFile(filePath)
	lines := strings.Split(string(bytesOfString), "\n")

	i := 0
	for index, line := range lines {

		splittedLine := strings.SplitN(line, ",", 3)
		if len(splittedLine) < 3 {
			log.Printf("splittedLine has than 3 elements: %+v", splittedLine)
			continue
		}

		if splittedLine[2] == "" {
			log.Printf("header is empty: %+v %v", splittedLine, index)
			continue
		}

		canBeUsed, mainRow, valueRow, orderedRow := HandleLogLine(
			line,
			filterCrawlers,
			parseUa,
			startLogTime,
			finishLogTime, udger)

		if canBeUsed {
			i++
			tx := db.Begin()
			_log := m.Log{
				Timestamp: h.UnixTimestampStrToTime(h.GetMapValueByKey(mainRow, "timestamp")),
				Ip:        h.GetMapValueByKey(mainRow, "ip", ),

				UserAgent:       h.GetMapValueByKey(mainRow, "User_Agent"),
				UaFamilyCode:    h.GetMapValueByKey(mainRow, "ua_family_code"),
				UaVersion:       h.GetMapValueByKey(mainRow, "ua_version"),
				UaClassCode:     h.GetMapValueByKey(mainRow, "ua_class_code"),
				DeviceClassCode: h.GetMapValueByKey(mainRow, "device_class_code"),
				OsFamilyCode:    h.GetMapValueByKey(mainRow, "os_family_code"),
				OsCode:          h.GetMapValueByKey(mainRow, "os_code"),

				ValueData: valueRow,
				OrderData: orderedRow,
			}
			tx.Create(&_log)
			if doRoolback {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}

	return nil
}

func ParseAndGetDataFromSingleGzFile(
	filePath string,
	startLogTime int64,
	finishLogTime int64) (logs []string) {

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		log.Printf("parsers.ParseAndStoreSingleGzLogInDb: Failed to connect database: %s", err)
		return logs
	}
	defer db.Close()
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}

	bytesOfString, _ := h.ReadGzFile(filePath)
	lines := strings.Split(string(bytesOfString), "\n")
	lines = lines[:10]

	for index, line := range lines {

		splittedLine := strings.SplitN(line, ",", 3)
		if len(splittedLine) < 3 {
			log.Printf("splittedLine has than 3 elements: %+v", splittedLine)
			continue
		}

		if splittedLine[2] == "" {
			log.Printf("header is empty: %+v %v", splittedLine, index)
			continue
		}

		headers_return := ParseLogLineWithCrawlers(
			line,
			startLogTime,
			finishLogTime)

		if headers_return != "" {
			logs = append(logs, headers_return)
		}
	}

	return
}

func GetLatestLogFile() (string, error) {
	logsDir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs")
	files := h.GetFileFromDirWithExt(logsDir, "gz")
	if len(files) == 0 {
		return "", errors.New(fmt.Sprintf("No files in path %s", logsDir))
	} else {
		return files[len(files)-1], nil
	}
}

func GetLatestLogFilePath() (string, string, error) {

	logsDir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs")
	files := h.GetFileFromDirWithExt(logsDir, "gz")

	if len(files) == 0 {
		log.Fatalf("parsers.GetLatestLogFilePath: There is no files(gz) in th dir %s", logsDir)
		return "", "", nil
	} else {

		// if file already loaded - return error
		fileName := files[len(files)-1]
		db, err := gorm.Open("postgres", m.GetDBConnectionStr())
		if err != nil {
			log.Fatalf("parser.go - main: Failed to connect database: %s", err)
		}
		defer db.Close()
		if !db.HasTable(&m.GzLog{}) {
			db.AutoMigrate(&m.GzLog{})
		}

		gzLog := m.GzLog{}
		db.Where("file_name = ?", fileName).First(&gzLog)

		if gzLog.ID != 0 {
			return filepath.Join(logsDir, fileName), fileName, errors.New(fmt.Sprintf("File %s already loaded to DB", fileName))
		}

		return filepath.Join(logsDir, fileName), fileName, nil
	}
}

func PrepareData(startLogTime int64, finishLogTime int64) (

	intUAClasses [][]int, floatUAClasses [][]float64, floatFullFeatures [][]float64, intFullFeatures [][]int) {

	userAgentList, trimmedValueData, trimmedOrderData := GetTrimmedLodMapsForPeriod(startLogTime, finishLogTime)

	valuesFeaturesOrder := FitValuesFeaturesOrder(trimmedValueData)

	floatFullFeatures, intFullFeatures = GetFullFeatures(trimmedOrderData, trimmedValueData, valuesFeaturesOrder)

	userAgentIntCodes, _ := FitUserAgentCodes(userAgentList)

	intUAClasses, floatUAClasses = GetUAClassesOneVsRest(userAgentList, userAgentIntCodes)

	return intUAClasses, floatUAClasses, floatFullFeatures, intFullFeatures
}

func PrepareUaFamilyCodes(limit int64) (
	intUaVersionClasses [][]int,
	floatUaVersionClasses [][]float64,
	floatFullFeatures [][]float64,
	intFullFeatures [][]int,
	uaFamilyCodeList []string,
	headersId []uint) {

	headersId, uaFamilyCodeList, trimmedValueData, trimmedOrderData := GetTrimmedDataWithUaFamilyCode(limit)

	valuesFeaturesOrder := FitValuesFeaturesOrder(trimmedValueData)

	floatFullFeatures, intFullFeatures = GetFullFeatures(trimmedOrderData, trimmedValueData, valuesFeaturesOrder)

	return
}

func GetTrimmedLodMapsForPeriod(
	startLogTime int64,
	finishLogTime int64) ([]string, []map[string]interface{}, []map[string]interface{}) {

	var (
		trimmedValueData []map[string]interface{}
		trimmedOrderData []map[string]interface{}
		userAgent        []string
		logs             = GetLogsInPeriod(startLogTime, finishLogTime)
	)

	for _, _log := range logs {
		trimmedValueData = append(trimmedValueData, _log.TrimValueData())
		trimmedOrderData = append(trimmedOrderData, _log.TrimOrderData())
		userAgent = append(userAgent, _log.UserAgent)
	}

	return userAgent, trimmedValueData, trimmedOrderData
}

func GetTrimmedDataWithUaFamilyCode(limit int64, ) ([]uint, []string, []map[string]interface{}, []map[string]interface{}) {

	var (
		trimmedValueData      []map[string]interface{}
		trimmedOrderData      []map[string]interface{}
		userAgentFamilyCodeId []uint
		userAgentVersion      []string
		logs                  = GetLogsWithLimit(limit)
	)

	println("Logs in sample:", len(logs))
	println("Educate from:", logs[len(logs)-1].Timestamp.Format("2006-01-02 15:04:05")  )
	println("Educate till:", logs[0].Timestamp.Format("2006-01-02 15:04:05"))

	for _, _log := range logs {
		trimmedValueData = append(trimmedValueData, _log.TrimValueData())
		trimmedOrderData = append(trimmedOrderData, _log.TrimOrderData())
		userAgentVersion = append(userAgentVersion, _log.UaFamilyCode)
		userAgentFamilyCodeId = append(userAgentFamilyCodeId, _log.ID)
	}

	return userAgentFamilyCodeId, userAgentVersion, trimmedValueData, trimmedOrderData
}

func GetLogsInPeriod(startLogTime int64, finishLogTime int64) []m.Log {

	start := time.Unix(startLogTime, 0)
	end := time.Unix(finishLogTime, 0)

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		log.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}

	logs := []m.Log{}
	db.Order("timestamp").Where("timestamp BETWEEN ? AND ?", start, end).Find(&logs)
	return logs
}

func GetLogsWithLimit(limit int64) []m.Log {
	if limit == 0 {
		return []m.Log{}
	}
	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		log.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}

	logs := []m.Log{}
	db.Limit(limit).Order("timestamp DESC").Find(&logs)
	return logs
}

func init() {
	h.LoadEnv()
}
