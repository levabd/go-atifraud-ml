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

	var need_all_file_parsing bool = false

	if len(args) == 0 || len(args) == 1 {
		return false
	}

	if args[1] == "true" {
		need_all_file_parsing = true
	}

	return need_all_file_parsing
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

	need_all_file_parsing := needAllFilesParsing(os.Args)
	start_log_time := helpers.StrToInt64(os.Getenv("PARSER_TIME_START"))
	finish_log_time := helpers.StrToInt64(os.Getenv("PARSER_TIME_END"))

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

	if need_all_file_parsing {

		logs_dir := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs")

		files := helpers.GetFileFromDirWithExt(logs_dir, "gz")
		files_to_handle := len(files)

		for i := 0; i < files_to_handle; i++ {
			fmt.Println(fmt.Printf("parse_gz_logs.go: main - File %s name ", files[i]))
			s.Logger.Printf("parse_gz_logs.go: main - File %s name ", files[i])

			gz_log := m.GzLog{}
			db.Where("file_name = ?", files[i]).First(&gz_log)

			if gz_log.ID != 0 {
				fmt.Println(fmt.Printf("File %s already loaded to DB ", files[i]))
				s.Logger.Printf("File %s already loaded to DB ", files[i])
				continue
			}

			e := s.ParseAndStoreSingleGzLogInDb(
				filepath.Join(logs_dir, files[i]),
				true,
				true,
				start_log_time,
				finish_log_time,
				false)
			if e != nil {
				fmt.Println(fmt.Sprintf("parse_gz_logs.go: main - Failed to parse ind store log from: %s ", files[i]))
				s.Logger.Fatalf("parse_gz_logs.go: main - Failed to parse ind store log from: %s ", files[i])
			}
			db.Create(&m.GzLog{FileName: files[i]})
			s.Logger.Printf("parse_gz_logs.go: main - File %s was parsed and stored in DB ", files[i])
		}
		fmt.Println(fmt.Sprintf("parse_gz_logs.go: main - Parsed and saved %v files", files_to_handle))
		s.Logger.Printf("parse_gz_logs.go: main - Parsed and saved %v files", files_to_handle)
		StartEducation()
		return
	}

	// single latest log gz file parsing
	full_file_path, file_name, err := s.GetLatestLogFilePath()
	if err != nil {
		fmt.Printf("parse_gz_logs.go: main - Getting latest log file failure: %s ", err)
		s.Logger.Fatalf("parse_gz_logs.go: main - Getting latest log file failure: %s ", err)
		return
	}
	println("full_file_path", full_file_path)
	// store new loaded log
	db.Create(&m.GzLog{FileName: file_name})

	e := s.ParseAndStoreSingleGzLogInDb(
		full_file_path,
		true,
		true,
		start_log_time,
		finish_log_time,
		false)

	if e != nil {
		fmt.Printf("parse_gz_logs.go: main - Failed to ParseAndStoreSingleGzLogInDb: %s", e)
		s.Logger.Fatalf("parse_gz_logs.go: main - Failed to ParseAndStoreSingleGzLogInDb: %s", e)
	}

	fmt.Printf("parse_gz_logs.go: main - Successfully parse file: %s ", e)
	s.Logger.Println("parse_gz_logs.go: main - Successfully parse file: ", file_name)

	StartEducation()
}

func StartEducation() {
	start_time := os.Getenv("PARSER_TIME_START")
	end_time := os.Getenv("PARSER_TIME_END")

	if start_time == "" || end_time == "" {
		panic("PARSER_TIME_START and PARSER_TIME_END must be set in env file")
	}

	trimmed_value_data, trimmed_order_data, pair_dict_list := s.PrepareData(helpers.StrToInt64(start_time), helpers.StrToInt64(end_time))

	println(len(trimmed_value_data), len(trimmed_order_data), len(pair_dict_list))
}

func init() {
	helpers.LoadEnv()
}
