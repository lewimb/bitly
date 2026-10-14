package domain

import (
	"time"

	"gorm.io/gorm"
)

type URLMapping struct {
	gorm.Model
	LongUrl   string     `gorm:"unique"`
	ShortCode string     `gorm:"type:varchar(7);uniqueIndex"`
	ExpiresAt *time.Time `gorm:"index"`
}

type ShortenResult struct {
	ID        uint
	ShortCode string
}

type ShortenPayload struct {
	LongUrl string `json:"long_url"`
}
