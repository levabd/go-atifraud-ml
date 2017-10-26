package models


type UACode struct {
	ID				uint
	UserAgent		string
	IntCode			int
	FloatCode		float64
}

func (UACode) TableName() string {
	return "ua_codes"
}