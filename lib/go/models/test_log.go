package models

import "time"

type TestLog struct {
	Model
	ID                   uint
	OriginalID           uint
	CreatedAt            time.Time
	Timestamp            time.Time
	Ip                   string
	UserAgent            string
	UaFamilyCode         string
	OriginalUaFamilyCode string
	UaVersion            string
	UaClassCode          string
	DeviceClassCode      string
	OsFamilyCode         string
	OsCode               string
	ValueData            JSONB  `sql:"type:JSONB NULL DEFAULT '{}'::JSONB"`
	OrderData            JSONB  `sql:"type:JSONB NULL DEFAULT '{}'::JSONB"`
}

func  (l *TestLog) TrimOrderData() map[string]interface{} {
	tmpRow := make(map[string]interface{})
	for key, value := range l.OrderData {

		if IsInImportantOrdersKeySet(key){
			tmpRow[key]=value
		}
	}
	return tmpRow
}

func  (l *TestLog) TrimValueData() map[string]interface{} {
	tmpRow := make(map[string]interface{})
	for key, value := range l.ValueData {
		if IsInImportantValueKeySet(key){
			tmpRow[key]=value
		}
	}
	return tmpRow
}
