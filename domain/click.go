package domain

import "gorm.io/gorm"

type URLClick struct {
	gorm.Model
	MappingID uint   `gorm:"index"`
	IPAddress string `gorm:"type:varchar(100)"`
	UserAgent string `gorm:"type:text"`
}
