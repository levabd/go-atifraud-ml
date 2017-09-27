package parsers

import (
	"strings"
	"github.com/buger/jsonparser"
	"fmt"
	"github.com/udger/udger"
	"sync"
	"path/filepath"
	"os"
	"time"
	"strconv"
	"bytes"
	"log"
)

var instance *udger.Udger
var once sync.Once

func GetUdgerInstance(path string) *udger.Udger {
	once.Do(func() {

		var path_to_udger_db string = path
		current_dir, err := filepath.Abs("./")
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(-1)
		}
		dir_ := filepath.Join(current_dir, path_to_udger_db)

		u, err := udger.New(dir_)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(-1)
		}

		instance = u
	})
	return instance
}

func convertTimestamp(t string) int64 {
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func timeIsWrong(log_timestamp string, start_log_time int64, finish_log_time int64) bool {
	time_in_logs := time.Unix(0, convertTimestamp(log_timestamp))

	start := time.Unix(0, start_log_time)
	finish := time.Unix(0, finish_log_time)

	return time_in_logs.Before(start) || time_in_logs.After(finish)
}

func HandleLogLine(line string, filter_crawlers bool, need_ua_parsing bool, start_log_time int64, finish_log_time int64) (bool, map[string]string, map[string]string, map[string]int) {
	var (
		result      bool              = true
		elements                      = strings.SplitN(line, ",", 3)
		main_row    map[string]string = make(map[string]string)
		ordered_row map[string]int    = make(map[string]int)
		value_row   map[string]string = make(map[string]string)
	)

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

		// define crawler in User-aGENT
		if ua, ok := value_row["User-Agent"]; ok {

			is_crawler := UaIsCrawler(elements[1], ua)

			fmt.Println("is_crawler: ",is_crawler, "\n")

			if filter_crawlers && is_crawler {
				result = false
				return result, nil, nil, nil
			}

			if need_ua_parsing {
				ua_obj := GetUa(ua)

				var buffer bytes.Buffer
				buffer.WriteString(ua_obj["ua_family_code"])
				buffer.WriteString(ua_obj["ua_version"])
				main_row["ua_family_code"] = ua_obj["ua_family_code"]
				main_row["ua_version"] = buffer.String()
				main_row["ua_class_code"] =  ua_obj["ua_class_code"]
				main_row["device_class_code"] = ua_obj["device_class_code"]
				main_row["os_family_code"] =  ua_obj["os_family_code"]
				main_row["os_code"] =  ua_obj["os_code"]
			}
		}

		main_row["timestamp"] = elements[0]
		main_row["ip"] = elements[1]
		main_row["User_Agent"] = value_row["User-Agent"]
	}

	return result, main_row, value_row, ordered_row
}

func ParseSingleLog(path_to_log string, filter_crawlers bool, parse_ua bool, start_log_time int64, finish_log_time int64)  {

	current_dir, err := filepath.Abs("./")
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}
	dir_ := filepath.Join(current_dir, path_to_log)
	log.Print(dir_)
}

