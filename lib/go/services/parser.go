package services

import (
	"strings"
	"github.com/buger/jsonparser"
	"fmt"
	"path/filepath"
	"os"
	"time"
	"bytes"
	h "github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"bufio"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
	"github.com/jinzhu/gorm"
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func timeIsWrong(logTimestamp string, startLogTime int64, finishLogTime int64,) bool {

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
	finishLogTime int64) (bool, map[string]interface{}, map[string]interface{}, map[string]interface{}) {

	var (
		result      bool                   = true
		elements                           = strings.SplitN(line, ",", 3)
		mainRow    map[string]interface{} = make(map[string]interface{})
		valueRow   map[string]interface{} = make(map[string]interface{})
		orderedRow map[string]interface{} = make(map[string]interface{})
	)

	if elements[0] == "" {
		result = false
		return false, nil, nil, nil
	}

	if timeIsWrong(elements[0], startLogTime, finishLogTime) {
		Logger.Printf("parsers.HandleLogLine: timeIsWrong: log_time - %s, start_log_time - %v, finish_log_time %v",
			string(elements[0]), startLogTime, finishLogTime)

		return false, nil, nil, nil
	}

	if len(elements) < 3{
		Logger.Printf("parsers.HandleLogLine: not enoth elements ")
		return false, nil, nil, nil
	}

	if filterCrawlers && (elements[0] == "" ||elements[1] == "" ||elements[2] == "") {
		Logger.Printf("parsers.HandleLogLine: not enoth elements ")
		return false, nil, nil, nil
	}

	if len(elements[2]) > 0 {
		if string(elements[2][0]) != "'" {
			Logger.Printf("parsers.HandleLogLine: string(elements[2][0]) != ' ")
			return false, nil, nil, nil
		}

		jsonToParse := strings.Replace(string(elements[2]), " ", "", -1)
		jsonToParse = strings.TrimPrefix(strings.TrimSuffix(jsonToParse, ""), "'")

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

		if valueRow["User-Agent"] ==nil{
			Logger.Printf("parsers.ParseAndStoreSingleGzLogInDb: No user agent in header %s ", jsonToParse)
			return false, nil, nil, nil
		}

		// define crawler in User-Agent
		if ua, ok := valueRow["User-Agent"].(string); ok {
			if filterCrawlers && IsCrawler(elements[1], ua) {
				result = false
				return result, nil, nil, nil
			}

			if needUaParsing {
				uaObj := GetUa(ua)

				var buffer bytes.Buffer
				buffer.WriteString(h.GetMapValueByKey(uaObj, "ua_family_code"))
				buffer.WriteString(h.GetMapValueByKey(uaObj, "ua_version"))
				mainRow["ua_family_code"] = h.GetMapValueByKey(uaObj, "ua_family_code")
				mainRow["ua_version"] = buffer.String()
				mainRow["ua_class_code"] = h.GetMapValueByKey(uaObj, "ua_class_code")
				mainRow["device_class_code"] = h.GetMapValueByKey(uaObj, "device_class_code")
				mainRow["os_family_code"] = h.GetMapValueByKey(uaObj, "os_family_code")
				mainRow["os_code"] = h.GetMapValueByKey(uaObj, "os_code")
			}
		}

		mainRow["timestamp"] = elements[0]
		mainRow["ip"] = elements[1]
		mainRow["User_Agent"] = valueRow["User-Agent"].(string)
	}

	return result, mainRow, valueRow, orderedRow
}

func ParseSingleLog(
	pathToLog string,
	filterCrawlers bool,
	parseUa bool,
	startLogTime int64,
	finishLogTime int64) ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {

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
			finishLogTime)

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
	doRoolback bool) error {

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		Logger.Printf("parsers.ParseAndStoreSingleGzLogInDb: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}
	bytesOfString, _ := h.ReadGzFile(filePath)
	lines := strings.Split(string(bytesOfString), "\n")

	for _, line := range lines {
		canBeUsed, mainRow, valueRow, orderedRow := HandleLogLine(
			line,
			filterCrawlers,
			parseUa,
			startLogTime,
			finishLogTime)

		if canBeUsed {
			tx := db.Begin()
			log := m.Log{
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
			tx.Create(&log)
			if doRoolback {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}

	return nil
}

func GetLatestLogFilePath() (string, string, error) {

	logsDir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs")
	files := h.GetFileFromDirWithExt(logsDir, "gz")

	if len(files) == 0 {
		Logger.Fatalf("parsers.GetLatestLogFilePath: There is no files(gz) in th dir %s", logsDir)
		return "", "", nil
	} else {

		// if file already loaded - return error
		fileName := files[len(files)-1]
		db, err := gorm.Open("postgres", m.GetDBConnectionStr())
		if err != nil {
			Logger.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
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

func PrepareData(
	startLogTime int64, finishLogTime int64,
)([]string, [][]bool, [][]bool ){

	userAgent, trimmedValueData, trimmedOrderData := GetTrimmedLodMapsForPeriod(startLogTime, finishLogTime)
	orderFeatures := GetOrderFeatures(trimmedOrderData)
	valuesFeaturesOrder := FitValuesFeaturesOrder(trimmedValueData)
	valueFeatures := GetValueFeatures(trimmedValueData, valuesFeaturesOrder)

	return userAgent, valueFeatures, orderFeatures
}

func GetTrimmedLodMapsForPeriod(
	startLogTime int64,
	finishLogTime int64)([]string, []map[string]interface{},[]map[string]interface{}){

	var(
		trimmedValueData []map[string]interface{}
		trimmedOrderData []map[string]interface{}
		userAgent []string
		logs = GetLogsInPeriod(startLogTime, finishLogTime)
	)

	for _, log:= range logs{
		trimmedValueData = append(trimmedValueData, log.TrimValueData())
		trimmedOrderData = append(trimmedOrderData, log.TrimOrderData())
		userAgent = append(userAgent, log.UserAgent)
	}

	return userAgent, trimmedValueData, trimmedOrderData
}

func GetLogsInPeriod(startLogTime int64, finishLogTime int64) []m.Log {

	start := time.Unix(startLogTime, 0)
	end := time.Unix(finishLogTime, 0)

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		Logger.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}

	logs := []m.Log{}
	db.Where("timestamp BETWEEN ? AND ?", start, end).Find(&logs)
	return logs
}

func init() {
	h.LoadEnv()
}
