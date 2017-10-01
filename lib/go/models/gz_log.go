package models

import (
	"time"
)

type GzLog struct {
	ID             uint
	FileName       string
	CreatedAt      time.Time
}

