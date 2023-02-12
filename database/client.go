package database

import (
	"jwt-authentication-golang/models"
	"log"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectionString struct {
        Host     string
        Port     int
        User     string
        Password string
        Dbname   string
}

var Instance *gorm.DB

var dbError error

func Connect(connectionString ConnectionString) () {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", connectionString.Host, connectionString.User, connectionString.Password, connectionString.Dbname, connectionString.Port)
	Instance, dbError = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database!")
}

func Migrate() {
	Instance.AutoMigrate(&models.User{})
	log.Println("Database Migration Completed!")
}
