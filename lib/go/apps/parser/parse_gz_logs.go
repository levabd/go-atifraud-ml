package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"path/filepath"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
	s "github.com/levabd/go-atifraud-ml/lib/go/services"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
	"fmt"
)

func init() {
	err := helpers.LoadEnv()
	if err != nil {
		s.Logger.Fatalln(err)
	}
}

func needAllFilesParsing(args []string) bool {

	var needAllFileParsing bool = false

	if len(args) == 0 || len(args) == 1 {
		return false
	}

	if args[1] == "true" {
		needAllFileParsing = true
	}

	return needAllFileParsing
}

// Parse GZ file/files and store data in DB
// If first command argument is true - all files
// in gir data/logs/*gz will we parsed and stored in DB
// Example(run from project dir):
//     go run lib/go/commands/parse_gz_logs.go true
//
// If there is no command argument - only latest
// gz file in data/logs/*gz will be parsed and stored
// Example(run from project dir):
//     go run lib/go/commands/parse_gz_logs.go [false]
func main() {

	if os.Getenv("PARSER_TIME_END") == "" || os.Getenv("PARSER_TIME_START") == "" {
		println("parse_gz_logs.go: main - PARSER_TIME_END and PARSER_TIME_START are not specified. Please check if there is an .env file in lib/go dir with such keys.")
		s.Logger.Fatalf("parse_gz_logs.go: main - PARSER_TIME_END and PARSER_TIME_START are not specified. Please check if there is an .env file in lib/go dir with such keys.")
		return
	}

	needAllFileParsing := needAllFilesParsing(os.Args)
	startLogTime := helpers.StrToInt64(os.Getenv("PARSER_TIME_START"))
	finishLogTime := helpers.StrToInt64(os.Getenv("PARSER_TIME_END"))

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		fmt.Println(fmt.Printf("parse_gz_logs.go: main - Failed to connect database: %s ", err))
		s.Logger.Fatalf("parse_gz_logs.go: main - Failed to connect database: %s ", err)
	}

	if !db.HasTable(&m.GzLog{}) {
		db.AutoMigrate(&m.GzLog{})
	}
	if !db.HasTable(&m.Log{}) {
		db.AutoMigrate(&m.Log{})
	}

	if needAllFileParsing {

		logsDir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs")

		files := helpers.GetFileFromDirWithExt(logsDir, "gz")
		filesToHandle := len(files)

		for i := 0; i < filesToHandle; i++ {
			fmt.Println(fmt.Printf("parse_gz_logs.go: main - File %s name ", files[i]))
			s.Logger.Printf("parse_gz_logs.go: main - File %s name ", files[i])

			gzLog := m.GzLog{}
			db.Where("file_name = ?", files[i]).First(&gzLog)

			if gzLog.ID != 0 {
				fmt.Println(fmt.Printf("File %s already loaded to DB ", files[i]))
				s.Logger.Printf("File %s already loaded to DB ", files[i])
				continue
			}

			e := s.ParseAndStoreSingleGzLogInDb(
				filepath.Join(logsDir, files[i]),
				true,
				true,
				startLogTime,
				finishLogTime,
				false)
			if e != nil {
				fmt.Println(fmt.Sprintf("parse_gz_logs.go: main - Failed to parse ind store log from: %s ", files[i]))
				s.Logger.Fatalf("parse_gz_logs.go: main - Failed to parse ind store log from: %s ", files[i])
			}
			db.Create(&m.GzLog{FileName: files[i]})
			s.Logger.Printf("parse_gz_logs.go: main - File %s was parsed and stored in DB ", files[i])
		}
		fmt.Println(fmt.Sprintf("parse_gz_logs.go: main - Parsed and saved %v files", filesToHandle))
		s.Logger.Printf("parse_gz_logs.go: main - Parsed and saved %v files", filesToHandle)
		StartEducation()
		return
	}

	// single latest log gz file parsing
	fullFilePath, fileName, err := s.GetLatestLogFilePath()
	if err != nil {
		fmt.Printf("parse_gz_logs.go: main - Getting latest log file failure: %s ", err)
		s.Logger.Fatalf("parse_gz_logs.go: main - Getting latest log file failure: %s ", err)
		return
	}
	println("full_file_path", fullFilePath)
	// store new loaded log
	db.Create(&m.GzLog{FileName: fileName})

	e := s.ParseAndStoreSingleGzLogInDb(
		fullFilePath,
		true,
		true,
		startLogTime,
		finishLogTime,
		false)

	if e != nil {
		fmt.Printf("parse_gz_logs.go: main - Failed to ParseAndStoreSingleGzLogInDb: %s", e)
		s.Logger.Fatalf("parse_gz_logs.go: main - Failed to ParseAndStoreSingleGzLogInDb: %s", e)
	}

	fmt.Printf("parse_gz_logs.go: main - Successfully parse file: %s ", e)
	s.Logger.Println("parse_gz_logs.go: main - Successfully parse file: ", fileName)

	StartEducation()
}

func StartEducation() {
	startTime := os.Getenv("PARSER_TIME_START")
	endTime := os.Getenv("PARSER_TIME_END")

	if startTime == "" || endTime == "" {
		panic("PARSER_TIME_START and PARSER_TIME_END must be set in env file")
	}

	userAgent, valueFeatures, orderFeatures := s.PrepareData(helpers.StrToInt64(startTime), helpers.StrToInt64(endTime))

	println(len(userAgent), len(valueFeatures), len(orderFeatures))
}

func init() {
	helpers.LoadEnv()
}
