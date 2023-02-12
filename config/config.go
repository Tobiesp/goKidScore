package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
                Host string
	        Port int
	        User string
	        Password string
	        DBName string
	}
	Server struct {
                Port int
	}
	Jwt struct {
		Timeout string
	}
}

var AppConfig *Config

func LoadAppConfig(){
	log.Println("Loading Server Configurations...")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal(err)
	}
}