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

func (m Model) JsonStrToMap(str string)  map[string]interface{}  {
	data := []byte(str)
	value_row:=make(map[string]interface{} )
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		value_row[ string(key)] = string(value)
		return nil
	})
	return value_row
}

func (m Model) IsJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil
}

func TruncateTable(table_name string) {
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
	tx.Exec("DELETE FROM " + table_name + ";")
	tx.Commit()
	defer db.Close()
}

