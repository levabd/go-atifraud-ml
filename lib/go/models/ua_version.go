package models

type UAVersion struct {
	ID				uint
	Version 		string
	IntCode			int
	FloatCode		float64
}

func (UAVersion) TableName() string {
	return "ua_browsers"
}