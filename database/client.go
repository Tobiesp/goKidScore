package database

import (
	"fmt"
	"kids-score/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

type ConnectionInfo struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
}

func Connect(connectionInfo ConnectionInfo) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		connectionInfo.Host, connectionInfo.Port, connectionInfo.User,
		connectionInfo.Password, connectionInfo.Dbname)
	Instance, dbError = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database!")
}

func Migrate() {
	Instance.AutoMigrate(&models.User{})
	log.Println("Database Migration Completed!")
	loadInitialData()
	log.Println("Initial data load complete!")
}

func loadInitialData() {
	loadInitialUserData()
	//TODO: Add load for groups
	//TODO: Add load for score types
}

func loadInitialUserData() {
	var user models.User
	var user2 models.User
	user.Email = "admin@kidscore.org"
	user.Name = "ksAdmin"
	user.Username = "ksAdmin"
	user.Enabled = true
	user.FailedLogin = 0
	user.Role = "Admin"
	if err := user.HashPassword(user.Password); err != nil {
		log.Println("Error: " + err.Error())
		return
	}
	count := int64(0)
	res := Instance.Model(&user2).
		Where("username = ?", user.Username).
		Count(&count)

	if res.Error != nil {
		log.Println("ERROR: " + res.Error.Error())
	}

	if count == 0 {
		// user does not exists
		Instance.Save(&user)
	}
}
