package models

import (
	"github.com/buger/jsonparser"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"bytes"
	"os"
)

type Model struct {}

func GetDBConnectionStr() string {
	var buffer bytes.Buffer
	buffer.WriteString("host=")
	buffer.WriteString(os.Getenv("DB_HOST"))
	buffer.WriteString(" ")
	buffer.WriteString("user=")
	buffer.WriteString(os.Getenv("DB_USERNAME"))
	buffer.WriteString(" ")
	buffer.WriteString("dbname=")
	buffer.WriteString(os.Getenv("DB_NAME"))
	buffer.WriteString(" ")
	buffer.WriteString("sslmode=")
	buffer.WriteString("disable")
	buffer.WriteString(" ")
	buffer.WriteString("password=")
	buffer.WriteString(os.Getenv("DB_PASSWORD"))
	return buffer.String()
}

//noinspection GoUnusedExportedFunction
func GetUdger() string {
	var buffer bytes.Buffer
	buffer.WriteString("host=")
	buffer.WriteString(os.Getenv("DB_HOST"))
	buffer.WriteString(" ")
	buffer.WriteString("user=")
	buffer.WriteString(os.Getenv("DB_USERNAME"))
	buffer.WriteString(" ")
	buffer.WriteString("dbname=")
	buffer.WriteString(os.Getenv("DB_NAME"))
	buffer.WriteString(" ")
	buffer.WriteString("sslmode=")
	buffer.WriteString("disable")
	buffer.WriteString(" ")
	buffer.WriteString("password=")
	buffer.WriteString(os.Getenv("DB_PASSWORD"))
	return buffer.String()
}

func (m Model) JsonStrToMap(str string)  map[string]interface{}  {
	data := []byte(str)
	valueRow:=make(map[string]interface{} )
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valueRow[ string(key)] = string(value)
		return nil
	})
	return valueRow
}

func (m Model) IsJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil
}

func TruncateTable(tableName string) {
	db, err := gorm.Open("postgres", GetDBConnectionStr())
	if err != nil {
		panic("failed to connect database")
	}

	if !db.HasTable(&GzLog{}) {
		db.AutoMigrate(&GzLog{})
	}
	if !db.HasTable(&Log{}) {
		db.AutoMigrate(&Log{})
	}

	// clear table
	tx := db.Begin()
	tx.Exec("DELETE FROM " + tableName + ";")
	tx.Commit()
	defer db.Close()
}

