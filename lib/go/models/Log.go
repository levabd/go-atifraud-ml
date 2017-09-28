package models

import (
	"time"
)

type Log struct {
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
	ValueData       string
	OrderData       string
}
