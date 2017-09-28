package helpers

import (
	"path/filepath"
	"os"
	"regexp"
	"compress/gzip"
	"io/ioutil"
	"strconv"
	"testing"
	"fmt"
	"encoding/json"
	"time"
)

func GetFileFromDirWithExt(path string,  ext string) []string {
	var files []string
	filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

func ReadGzFile(filename string) ([]byte, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func UnixTimestampStrToTime(str string ) time.Time {
	if str == ""{
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}
	tm := time.Unix(StrToInt(str), 0)
	return tm
}

func StrToInt(t string) int64 {
	if t == ""{
		return 0
	}
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}


func GetMapValueByKey(main_row map[string]interface{}, key string) string {

	if main_row[key] == nil {
		return ""
	}

	if main_row[key] == "" {
		return ""
	}

	return main_row[key].(string)
}

func IsJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil

}
