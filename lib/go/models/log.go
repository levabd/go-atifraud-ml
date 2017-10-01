package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"
)

type JSONB map[string]interface{}

type Log struct {
	Model
	ID              uint
	CreatedAt       time.Time
	Timestamp       time.Time
	Ip              string
	UserAgent       string
	UaFamilyCode    string
	UaVersion       string
	UaClassCode     string
	DeviceClassCode string
	OsFamilyCode    string
	OsCode          string
	ValueData       JSONB  `sql:"type:JSONB NULL DEFAULT '{}'::JSONB"`
	OrderData       JSONB  `sql:"type:JSONB NULL DEFAULT '{}'::JSONB"`
}

func (Log) TableName() string {
	return "logs"
}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

func  (l *Log) TrimOrderData() map[string]interface{} {
	tmp_row := make(map[string]interface{})
	for key, value := range l.OrderData {
		if IsInImportantOrdersKeySet(key){
			tmp_row[key]=value
		}
	}
	return tmp_row
}

func  (l *Log) TrimValueData() map[string]interface{} {
	tmp_row := make(map[string]interface{})
	for key, value := range l.ValueData {
		if IsInImportantValueKeySet(key){
			tmp_row[key]=value
		}
	}
	return tmp_row
}

func IsInImportantOrdersKeySet(category string) bool {
	switch category {
	case
		"Upgrade-Insecure-Requests",
		"Accept",
		"If-Modified-Since",
		"Host",
		"Connection",
		"User-Agent",
		"From",
		"Accept-Encoding":
		return true
	}
	return false
}

func IsInImportantValueKeySet(category string) bool {
	switch category {
	case
		"Accept",
		"Accept-Charset",
		"Accept-Encoding":
		return true
	}
	return false
}