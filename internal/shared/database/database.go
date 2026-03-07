package database

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func Init() *Database {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:%25Andika12345@tcp(127.0.0.1:3306)/popda_2026?parseTime=true&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Ping(); err != nil {
			log.Fatal("Database ping failed:", err)
		}
	} else {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Database connected successfully")
	return &Database{
		DB: db,
	}
}
