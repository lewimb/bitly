package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Host                 string
	User                 string
	Password             string
	DBName               string
	Port                 string
	TimeZone             string
	SSLMode              string
	PreferSimpleProtocol bool
}

func NewDB() *DB {
	return &DB{
		Host:                 os.Getenv("HOST"),
		User:                 os.Getenv("USER"),
		Password:             os.Getenv("Password"),
		DBName:               os.Getenv("DB_NAME"),
		Port:                 os.Getenv("PORT"),
		TimeZone:             os.Getenv("TIMEZONE"),
		SSLMode:              os.Getenv("SSLMODE"),
		PreferSimpleProtocol: true,
	}
}

func (db *DB) DSN() string {
	return fmt.Sprintf("host=%s password=%s user=%s dbname=%s timezone=%s sslmode=%s",
		db.Host, db.Password, db.User, db.DBName, db.TimeZone, db.SSLMode)
}

func (d *DB) Connect() (db *gorm.DB, err error) {
	return gorm.Open(postgres.New(postgres.Config{DSN: d.DSN(), PreferSimpleProtocol: d.PreferSimpleProtocol}))
}
