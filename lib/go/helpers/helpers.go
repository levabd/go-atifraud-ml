package helpers

import (
	"path/filepath"
	"os"
	"regexp"
	"compress/gzip"
	"io/ioutil"
	"strconv"
	"time"
	"log"
	"github.com/joho/godotenv"
	"flag"
)

var env_is_loaded = false
var env_is_testing = false

func IsTesting() bool {
	if env_is_testing {
		return true
	}

	return flag.Lookup("test.v") != nil
}

func LoadEnv() error {

	if env_is_loaded {
		return nil
	}

	path_to_env, err := filepath.Abs("./.")
	if err != nil {
		log.Fatal(err)
		return err
	}

	if env_is_testing = IsTesting(); env_is_testing {
		path_to_env = filepath.Join(path_to_env, "..", ".env")
	} else {
		path_to_env = filepath.Join(path_to_env, "lib", "go", ".env")
	}
	if _, err := os.Stat(path_to_env); os.IsNotExist(err) {
		return err
	}

	godotenv.Load(path_to_env)
	env_is_loaded = true

	return nil
}

func GetFileFromDirWithExt(path string, ext string) []string {
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

func UnixTimestampStrToTime(str string) time.Time {
	if str == "" {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}
	tm := time.Unix(StrToInt64(str), 0)
	return tm
}

func StrToInt64(t string) int64 {
	if t == "" {
		return 0
	}
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func StrToInt(t string) int64 {
	if t == "" {
		return 0
	}
	i, err := strconv.ParseInt(t, 10, 0)
	if err != nil {
		panic(err)
	}

	return i
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


func GetRandomKeyFromMap(_map map[string]interface{}) string {
	var _key string
	for _key, _ = range _map {
		break
	}

	return _key
}