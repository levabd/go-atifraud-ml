package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
)

var errorLogFile *os.File
var Logger *log.Logger

func init() {
	helpers.LoadEnv()

	abs_path, err := filepath.Abs("./.")
	if err != nil {
		log.Fatal(err)
		return
	}
	var file_path string

	if helpers.IsTesting() {
		file_path = filepath.Join(abs_path, "..", "..", "..",  "data", "logs_go", "filename.log")
		fmt.Println("run under go test")
	} else {
		file_path = filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs_go", "filename.log")
		fmt.Println("normal run")
	}

	fmt.Println("services.init: path to log file - ", file_path)
	errorLogFile, err = os.OpenFile(file_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %s, error: %v", file_path, err)
		os.Exit(1)
	}

	Logger = log.New(errorLogFile, "applog: ", log.Lshortfile|log.LstdFlags)
}