package parsers

import (
	"strings"
	"github.com/buger/jsonparser"
	"fmt"
	"path/filepath"
	"os"
	"time"
	"bytes"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"bufio"
	"encoding/json"
	"github.com/levabd/go-atifraud-ml/lib/go/models"
	"github.com/jinzhu/gorm"
)

func timeIsWrong(log_timestamp string, start_log_time int64, finish_log_time int64) bool {

	time_in_logs := time.Unix(0, helpers.StrToInt(log_timestamp))

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
		result = false
		return false, nil, nil, nil
	}

	if filter_crawlers && elements[1] == "" {
		result = false
		return result, nil, nil, nil
	}

	if string(elements[2][0]) != "'" {
		result = false
		return result, nil, nil, nil
	}

	json_to_parse := strings.Replace(string(elements[2]), " ", "", -1)
	json_to_parse = strings.TrimPrefix(strings.TrimSuffix(json_to_parse, ""), "'")

	if len(elements[2]) > 0 {
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


		// define crawler in User-Agent
		if ua, ok := value_row["User-Agent"].(string); ok {
			if filter_crawlers && UaIsCrawler(elements[1], ua) {
				result = false
				return result, nil, nil, nil
			}

			if need_ua_parsing {
				ua_obj := GetUa(ua)

				var buffer bytes.Buffer
				buffer.WriteString(helpers.GetMapValueByKey(ua_obj, "ua_family_code"))
				buffer.WriteString(helpers.GetMapValueByKey(ua_obj, "ua_version"))
				main_row["ua_family_code"] = helpers.GetMapValueByKey(ua_obj, "ua_family_code")
				main_row["ua_version"] = buffer.String()
				main_row["ua_class_code"] = helpers.GetMapValueByKey(ua_obj, "ua_class_code" )
				main_row["device_class_code"] = helpers.GetMapValueByKey(ua_obj, "device_class_code" )
				main_row["os_family_code"] = helpers.GetMapValueByKey(ua_obj, "os_family_code" )
				main_row["os_code"] = helpers.GetMapValueByKey(ua_obj, "os_code" )
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
		_, main_row, value_row, ordered_row := HandleLogLine(
			scanner.Text(),
			filter_crawlers,
			parse_ua,
			start_log_time,
			finish_log_time)
		main_table = append(main_table, main_row)
		value_table = append(value_table, value_row)
		ordered_table = append(ordered_table, ordered_row)
	}

	return main_table, value_table, ordered_table
}

func ParseLatestLogGzFile(
	filter_crawlers bool,
	parse_ua bool,
	start_log_time int64,
	finish_log_time int64, do_roolback bool) ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {

	full_file_path := getLatestLogFilePath()

	str_in_bytes, _ := helpers.ReadGzFile(full_file_path)
	splitted_lines := strings.Split(string(str_in_bytes), "\n")

	var (
		main_table    []map[string]interface{}
		value_table   []map[string]interface{}
		ordered_table []map[string]interface{}
	)

	db, err := gorm.Open("sqlite3", os.Getenv("DB_FILE_ANTIFRAUD"))
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&models.Log{})
	tx := db.Begin()

	for _, line := range splitted_lines {
		_, main_row, value_row, ordered_row := HandleLogLine(
			line,
			filter_crawlers,
			parse_ua,
			start_log_time,
			finish_log_time)
		main_table = append(main_table, main_row)
		value_table = append(value_table, value_row)
		ordered_table = append(ordered_table, ordered_row)

		value_row_json, _ := json.Marshal(value_row);
		ordered_row_json, _ := json.Marshal(ordered_row);
		log := models.Log{
			Timestamp: helpers.UnixTimestampStrToTime(helpers.GetMapValueByKey(main_row, "timestamp")),
			Ip:        helpers.GetMapValueByKey(main_row, "ip", ),

			UserAgent:       helpers.GetMapValueByKey(main_row, "User_Agent"),
			UaFamilyCode:    helpers.GetMapValueByKey(main_row, "ua_family_code"),
			UaVersion:       helpers.GetMapValueByKey(main_row, "ua_version"),
			UaClassCode:     helpers.GetMapValueByKey(main_row, "ua_class_code"),
			DeviceClassCode: helpers.GetMapValueByKey(main_row, "device_class_code"),
			OsFamilyCode:    helpers.GetMapValueByKey(main_row, "os_family_code"),
			OsCode:          helpers.GetMapValueByKey(main_row, "os_code"),

			ValueData: string(value_row_json),
			OrderData: string(ordered_row_json),
		}
		db.Create(&log)
	}

	if do_roolback {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return main_table, value_table, ordered_table
}



func getLatestLogFilePath() string {

	logs_dir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs")

	helpers.GetFileFromDirWithExt(logs_dir, "gz")

	files := helpers.GetFileFromDirWithExt(logs_dir, "gz")

	if len(files) == 0 {
		panic(fmt.Sprintf("There is no files(gz) in th dir %s", logs_dir))
	} else {
		return filepath.Join(logs_dir, files[len(files)-1])
	}
}
