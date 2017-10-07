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

var envIsLoaded = false
var envIsTesting = false



func IsTesting() bool {
	if envIsTesting {
		return true
	}

	return flag.Lookup("test.v") != nil
}

func LoadEnv() error {

	if envIsLoaded {
		return nil
	}

	pathToEnv, err := filepath.Abs("./.")
	if err != nil {
		log.Fatal(err)
		return err
	}

	if envIsTesting = IsTesting(); envIsTesting {
		pathToEnv = filepath.Join(pathToEnv, "..", ".env")
	} else {
		pathToEnv = filepath.Join(pathToEnv, "lib", "go", ".env")
	}
	if _, err := os.Stat(pathToEnv); os.IsNotExist(err) {
		return err
	}

	godotenv.Load(pathToEnv)
	envIsLoaded = true

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
		fl, err := strconv.ParseFloat(t, 0)
		if err != nil {
			panic(err)
		}

		return  int64(fl)
	}

	return i
}

func GetMapValueByKey(mainRow map[string]interface{}, key string) string {

	if mainRow[key] == nil {
		return ""
	}

	if mainRow[key] == "" {
		return ""
	}

	return mainRow[key].(string)
}


func GetRandomKeyFromMap(input map[string]interface{}) string {
	var key string
	for key = range input {
		break
	}

	return key
}