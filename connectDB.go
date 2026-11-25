package main

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rndmBackend/internal/domain"
)

var dsn = "host=localhost user=postgres password=test123 dbname=r4n port=5432 sslmode=disable"


func GetDbConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&domain.User{})
	slog.Info("Succefull db connect")
	return db
}