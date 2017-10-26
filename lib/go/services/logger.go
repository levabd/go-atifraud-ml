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

	absPath, err := filepath.Abs("./.")
	if err != nil {
		log.Fatal(err)
		return
	}
	var filePath string

	if helpers.IsTesting() {
		filePath = filepath.Join(absPath, "..", "..", "..",  "data", "logs_go", "filename.log")
		fmt.Println("Run app not for testing")
	} else {
		filePath = filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "logs_go", "filename.log")
		fmt.Println("Run app not for production")
	}

	fmt.Println("services.init: path to log file - ", filePath)
	errorLogFile, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %s, error: %v", filePath, err)
		os.Exit(1)
	}

	Logger = log.New(errorLogFile, "applog: ", log.Lshortfile|log.LstdFlags)
}