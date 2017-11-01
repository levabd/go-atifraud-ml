package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"path/filepath"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
	s "github.com/levabd/go-atifraud-ml/lib/go/services"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
	"github.com/levabd/go-atifraud-ml/lib/go/udger"
	"fmt"
	"log"
	"github.com/uniplaces/carbon"
	"time"
	"gopkg.in/cheggaaa/pb.v1"
	"github.com/valyala/fasthttp"
	"strings"
	"os/exec"
	"flag"
)

var (
	parseAllFiles = flag.Bool("all", false, "host:port to listen to")
	parseOneFile  = flag.String("one", "path", "full file path to parse")
	onlyEducate   = flag.Bool("educate", false, "educate model")
	sample        = flag.Int64("sample", 90000, "education sample")
	connection    = &fasthttp.Client{}
	req           = fasthttp.AcquireRequest()
	resp          = fasthttp.AcquireResponse()
)

func init() {
	err := helpers.LoadEnv()
	if err != nil {
		log.Fatalln(err)
	}
}

// Parse GZ file/files and store data in DB - then train model
func main() {

	flag.Parse()

	if *onlyEducate {
		log.Println("Only education mode")
		log.Println("Education sample", *sample)
		EducateModel()
		return
	}

	defer timeTrack(time.Now(), "parse logs")

	if os.Getenv("PARSER_TIME_END") == "" || os.Getenv("PARSER_TIME_START") == "" {
		println("parse_gz_logs.go: main - PARSER_TIME_END and PARSER_TIME_START are not specified. Please check if there is an .env file in lib/go dir with such keys.")
		log.Fatalf("parse_gz_logs.go: main - PARSER_TIME_END and PARSER_TIME_START are not specified. Please check if there is an .env file in lib/go dir with such keys.")
		return
	}

	startLogTime := helpers.StrToInt64(os.Getenv("PARSER_TIME_START"))
	finishLogTime := carbon.Now().Unix()

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		fmt.Println(fmt.Printf("parse_gz_logs.go: main - Failed to connect database: %s ", err))
		log.Fatalf("parse_gz_logs.go: main - Failed to connect database: %s ", err)
	}

	if !db.HasTable(&m.GzLog{}) {
		db.AutoMigrate(&m.GzLog{})
	}
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}

	_udger, err := udger.New(os.Getenv("DB_FILE_PATH_UDGER"))
	if err != nil {
		panic(err)
	}

	if *parseOneFile != "" && *parseOneFile != "path" {
		log.Println("Parse one file mode")

		if _, err := os.Stat(*parseOneFile); os.IsNotExist(err) {
			log.Fatalf("parse_gz_logs.go: main - Privided path to file is not valid - file does not exist: %s ", *parseOneFile)
			return
		}

		_, fileName := filepath.Split(*parseOneFile)

		log.Fatalf("parse_gz_logs.go: main - File name: %s ", fileName)

		if s.GzLogFileLoaded(fileName) {
			log.Fatalf("parse_gz_logs.go: main -File already loaded in DB: %s ", fileName)
			return
		}

		handleFile(*parseOneFile, startLogTime, finishLogTime, _udger, db, fileName)
		return
	}

	if *parseAllFiles {
		log.Println("Parse all files mode")

		logsDir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs")
		files := helpers.GetFileFromDirWithExt(logsDir, "gz")
		filesToHandle := len(files)

		for i := 0; i < filesToHandle; i++ {

			log.Println(fmt.Sprintf("parse_gz_logs.go: main - Handle file with name %s", files[i]))
			log.Printf("parse_gz_logs.go: main - File %s name ", files[i])

			gzLog := m.GzLog{}
			db.Where("file_name = ?", files[i]).First(&gzLog)

			if gzLog.ID != 0 {
				log.Println(fmt.Printf("File %s already loaded to DB ", files[i]))
				log.Printf("File %s already loaded to DB ", files[i])
				continue
			}

			e := s.ParseAndStoreSingleGzLogInDb(
				filepath.Join(logsDir, files[i]),
				true,
				true,
				startLogTime,
				finishLogTime,
				false, _udger)

			if e != nil {
				log.Println(fmt.Sprintf("parse_gz_logs.go: main - Failed to parse ind store log from: %s ", files[i]))
				log.Fatalf("parse_gz_logs.go: main - Failed to parse ind store log from: %s ", files[i])
			}
			db.Create(&m.GzLog{FileName: files[i]})
			log.Printf("parse_gz_logs.go: main - File %s was parsed and stored in DB ", files[i])
		}

		fmt.Println(fmt.Sprintf("parse_gz_logs.go: main - Parsed and saved %v files", filesToHandle))
		log.Printf("parse_gz_logs.go: main - Parsed and saved %v files", filesToHandle)
		EducateModel()
		return
	}

	log.Println("Parse the latest file mode")

	// single latest log gz file parsing
	fullFilePath, fileName, err := s.GetLatestLogFilePath()
	if err != nil {
		fmt.Printf("parse_gz_logs.go: main - Getting latest log file failure: %s ", err)
		log.Fatalf("parse_gz_logs.go: main - Getting latest log file failure: %s ", err)
		return
	}

	handleFile(fullFilePath, startLogTime, finishLogTime, _udger, db, fileName)

	EducateModel()
}

func handleFile(
	fullFilePath string,
	startLogTime int64,
	finishLogTime int64,
	_udger *udger.Udger,
	db *gorm.DB,
	fileName string,
) {

	e := s.ParseAndStoreSingleGzLogInDb(
		fullFilePath,
		true,
		true,
		startLogTime,
		finishLogTime,
		false, _udger)

	// store new loaded log
	db.Create(&m.GzLog{FileName: fileName})

	if e != nil {
		fmt.Printf("parse_gz_logs.go: main - Failed to ParseAndStoreSingleGzLogInDb: %s", e)
		log.Fatalf("parse_gz_logs.go: main - Failed to ParseAndStoreSingleGzLogInDb: %s", e)
	}
	fmt.Printf("parse_gz_logs.go: main - Successfully parse file: %s ", e)
	log.Println("parse_gz_logs.go: main - Successfully parse file: ", fileName)
}

func EducateModel() {

	// PrepareDataForTrain
	prepareDataForTrain()

	// train
	train()

	// reload model on prediction server
	reloadModelOnPythonServer()
}

func reloadModelOnPythonServer() {

	log.Println("Start reloading model on prediction server")

	portList := []string{"5000", "5001", "5002", "5003", "5004", "5005", "5006", "5007"}

	for _, port := range portList {
		_response := doRequest("http://0.0.0.0:" + port + "/reload")

		if string(_response) == "reloaded" {
			log.Printf("Prediction model reloaded on python server: 0.0.0.0:%s", port)
		} else {
			log.Println("Problem while reloading prediction model on python server: 0.0.0.0:%s. Response: %s", port, string(_response))
		}
	}
}

func train() {
	log.Println("Start training")

	cmd := exec.Command("./lib/python/train")
	cmd.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error python exec: ", err)
		os.Exit(-1)
	}
	if strings.Contains(string(out), "Education finished") {
		log.Println("Education finished")
	} else {
		log.Println("Problem detected while onlyEducate model")
	}
}

func prepareDataForTrain() {
	log.Println("Start preparing data for train")

	_, _, _, intFullFeatures, uaFamilyCodesList, logIds := s.PrepareUaFamilyCodes(*sample)
	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		panic("user_agent_helpers.go - LoadFittedUserAgentCodes: Failed to connect to database")
	}
	defer db.Close()
	if !db.HasTable(&m.Features{}) {
		db.AutoMigrate(&m.Features{})
	}
	if !db.HasTable(&m.Browsers{}) {
		db.AutoMigrate(&m.Browsers{})
	}
	db.Exec("TRUNCATE TABLE features;")
	db.Exec("TRUNCATE TABLE browsers;")
	tx := db.Begin()
	bar := pb.StartNew(len(intFullFeatures))
	bar.SetRefreshRate(time.Second)
	for index_row, featureValues := range intFullFeatures {
		for index_column, value := range featureValues {
			if value == 1 {
				cacheFeatures := m.Features{
					LogId:  logIds[index_row],
					Row:    index_row,
					Column: index_column,
				}
				tx.Create(&cacheFeatures)
			}
		}
		bar.Increment()
	}
	tx.Commit()
	tx = db.Begin()
	bar = pb.StartNew(len(uaFamilyCodesList))
	bar.SetRefreshRate(time.Second)
	for i, name := range uaFamilyCodesList {
		bar.Increment()
		cacheFeatures := m.Browsers{Name: name, LogId: logIds[i]}
		tx.Create(&cacheFeatures)
	}
	tx.Commit()

	log.Println("Finish preparing data for train")
}

func doRequest(url string) []byte {

	req.SetRequestURI(url)

	connection.Do(req, resp)

	return resp.Body()
}

func init() {
	helpers.LoadEnv()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
