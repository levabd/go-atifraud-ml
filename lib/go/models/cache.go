package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Features struct {
	LogId uint
	Row int
	Column int
}

type JSONB_INT []int

type Browsers struct {
	LogId uint
	Name string
}

func (j JSONB_INT) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB_INT) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}