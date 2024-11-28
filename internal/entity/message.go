package entity

import (
	"time"
)

type Message struct {
	ID        uint `gorm:"primaryKey"`
	Hostname  string
	Message   string
	CreatedAt time.Time
}
