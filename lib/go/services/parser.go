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

func timeIsWrong(log_timestamp string, start_log_time int64, finish_log_time int64) bool {

	time_in_logs := time.Unix(0, h.StrToInt64(log_timestamp))

	start := time.Unix(0, start_log_time)
	finish := time.Unix(0, finish_log_time)

	return time_in_logs.Before(start) || time_in_logs.After(finish)
}

func HandleLogLine(
	line string,
	filter_crawlers bool,
	need_ua_parsing bool,
	start_log_time int64,
	finish_log_time int64) (bool, map[string]interface{}, map[string]interface{}, map[string]interface{}) {

	var (
		result      bool                   = true
		elements                           = strings.SplitN(line, ",", 3)
		main_row    map[string]interface{} = make(map[string]interface{})
		value_row   map[string]interface{} = make(map[string]interface{})
		ordered_row map[string]interface{} = make(map[string]interface{})
	)

	if elements[0] == "" {
		result = false
		return false, nil, nil, nil
	}

	if timeIsWrong(elements[0], start_log_time, finish_log_time) {
		Logger.Printf("parsers.HandleLogLine: timeIsWrong: log_time - %s, start_log_time - %v, finish_log_time %v",
			string(elements[0]), start_log_time, finish_log_time)

		return false, nil, nil, nil
	}

	if len(elements) < 3{
		Logger.Printf("parsers.HandleLogLine: not enoth elements ")
		return false, nil, nil, nil
	}

	if filter_crawlers && (elements[0] == "" ||elements[1] == "" ||elements[2] == "") {
		Logger.Printf("parsers.HandleLogLine: not enoth elements ")
		return false, nil, nil, nil
	}

	if len(elements[2]) > 0 {
		if string(elements[2][0]) != "'" {
			Logger.Printf("parsers.HandleLogLine: string(elements[2][0]) != ' ")
			return false, nil, nil, nil
		}

		json_to_parse := strings.Replace(string(elements[2]), " ", "", -1)
		json_to_parse = strings.TrimPrefix(strings.TrimSuffix(json_to_parse, ""), "'")

		data := []byte(json_to_parse)
		i := 0
		jsonparser.ObjectEach(data, func(
			key []byte,
			value []byte,
			dataType jsonparser.ValueType,
			offset int) error {
			ordered_row[ string(key)] = i
			value_row[ string(key)] = string(value)
			i = i + 1
			return nil
		})

		if value_row["User-Agent"] ==nil{
			Logger.Printf("parsers.ParseAndStoreSingleGzLogInDb: No user agent in header %s ", json_to_parse)
			return false, nil, nil, nil
		}

		// define crawler in User-Agent
		if ua, ok := value_row["User-Agent"].(string); ok {
			if filter_crawlers && IsCrawler(elements[1], ua) {
				result = false
				return result, nil, nil, nil
			}

			if need_ua_parsing {
				ua_obj := GetUa(ua)

				var buffer bytes.Buffer
				buffer.WriteString(h.GetMapValueByKey(ua_obj, "ua_family_code"))
				buffer.WriteString(h.GetMapValueByKey(ua_obj, "ua_version"))
				main_row["ua_family_code"] = h.GetMapValueByKey(ua_obj, "ua_family_code")
				main_row["ua_version"] = buffer.String()
				main_row["ua_class_code"] = h.GetMapValueByKey(ua_obj, "ua_class_code")
				main_row["device_class_code"] = h.GetMapValueByKey(ua_obj, "device_class_code")
				main_row["os_family_code"] = h.GetMapValueByKey(ua_obj, "os_family_code")
				main_row["os_code"] = h.GetMapValueByKey(ua_obj, "os_code")
			}
		}

		main_row["timestamp"] = elements[0]
		main_row["ip"] = elements[1]
		main_row["User_Agent"] = value_row["User-Agent"].(string)
	}

	return result, main_row, value_row, ordered_row
}

func ParseSingleLog(
	path_to_log string,
	filter_crawlers bool,
	parse_ua bool,
	start_log_time int64,
	finish_log_time int64) ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {

	inFile, _ := os.Open(path_to_log)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var (
		main_table    []map[string]interface{}
		value_table   []map[string]interface{}
		ordered_table []map[string]interface{}
	)

	for scanner.Scan() {
		is_not_bot, main_row, value_row, ordered_row := HandleLogLine(
			scanner.Text(),
			filter_crawlers,
			parse_ua,
			start_log_time,
			finish_log_time)

		if is_not_bot {
			main_table = append(main_table, main_row)
			value_table = append(value_table, value_row)
			ordered_table = append(ordered_table, ordered_row)
		}
	}

	return main_table, value_table, ordered_table
}

func ParseAndStoreSingleGzLogInDb(
	file_path string,
	filter_crawlers bool,
	parse_ua bool,
	start_log_time int64,
	finish_log_time int64,
	do_roolback bool) error {

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		Logger.Printf("parsers.ParseAndStoreSingleGzLogInDb: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}
	str_in_bytes, _ := h.ReadGzFile(file_path)
	lines := strings.Split(string(str_in_bytes), "\n")

	for _, line := range lines {
		can_be_used, main_row, value_row, ordered_row := HandleLogLine(
			line,
			filter_crawlers,
			parse_ua,
			start_log_time,
			finish_log_time)

		if can_be_used {
			tx := db.Begin()
			log := m.Log{
				Timestamp: h.UnixTimestampStrToTime(h.GetMapValueByKey(main_row, "timestamp")),
				Ip:        h.GetMapValueByKey(main_row, "ip", ),

				UserAgent:       h.GetMapValueByKey(main_row, "User_Agent"),
				UaFamilyCode:    h.GetMapValueByKey(main_row, "ua_family_code"),
				UaVersion:       h.GetMapValueByKey(main_row, "ua_version"),
				UaClassCode:     h.GetMapValueByKey(main_row, "ua_class_code"),
				DeviceClassCode: h.GetMapValueByKey(main_row, "device_class_code"),
				OsFamilyCode:    h.GetMapValueByKey(main_row, "os_family_code"),
				OsCode:          h.GetMapValueByKey(main_row, "os_code"),

				ValueData: value_row,
				OrderData: ordered_row,
			}
			tx.Create(&log)
			if do_roolback {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}

	return nil
}

func GetLatestLogFilePath() (string, string, error) {

	logs_dir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs")
	files := h.GetFileFromDirWithExt(logs_dir, "gz")

	if len(files) == 0 {
		Logger.Fatalf("parsers.GetLatestLogFilePath: There is no files(gz) in th dir %s", logs_dir)
		return "", "", nil
	} else {

		// if file already loaded - return error
		file_name := files[len(files)-1]
		db, err := gorm.Open("postgres", m.GetDBConnectionStr())
		if err != nil {
			Logger.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
		}
		defer db.Close()
		if !db.HasTable(&m.GzLog{}) {
			db.AutoMigrate(&m.GzLog{})
		}

		gz_log := m.GzLog{}
		db.Where("file_name = ?", file_name).First(&gz_log)

		if gz_log.ID != 0 {
			return filepath.Join(logs_dir, file_name), file_name, errors.New(fmt.Sprintf("File %s already loaded to DB", file_name))
		}

		return filepath.Join(logs_dir, file_name), file_name, nil
	}
}

func PrepareData(start_log_time int64, finish_log_time int64) {
	trimmed_value_data, trimmed_order_data :=GetTrimmedLodMapsForPeriod(start_log_time, finish_log_time)
	pair_dict_list := GetPairsDictList(trimmed_order_data)

	println(fmt.Sprintf("trimmed_value_data %v", len(trimmed_value_data)))
	println(fmt.Sprintf("pair_dict_list len %v", len(pair_dict_list )))
	//println(fmt.Sprintf("pair_dict_list value %v", pair_dict_list ))
}

func GetTrimmedLodMapsForPeriod(
	start_log_time int64,
	finish_log_time int64)([]map[string]interface{},[]map[string]interface{}){

	var(
		trimmed_value_data []map[string]interface{}
		trimmed_order_data []map[string]interface{}
		logs = GetLogsInPeriod(start_log_time, finish_log_time)
	)

	for _, log:= range logs{
		trimmed_value_data = append(trimmed_value_data, log.TrimValueData())
		trimmed_order_data = append(trimmed_order_data, log.TrimOrderData())
	}

	return trimmed_value_data, trimmed_order_data
}

func GetLogsInPeriod(start_log_time int64, finish_log_time int64) []m.Log {

	start := time.Unix(start_log_time, 0)
	end := time.Unix(finish_log_time, 0)

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
