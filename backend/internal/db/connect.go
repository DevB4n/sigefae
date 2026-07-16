package db

import (
	"fmt"

	"sigefae/internal/env"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *env.Env) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBCharset,
		cfg.DBParseTime,
		cfg.DBLoc,
	)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = conn

	return conn, nil
}