package main

import (
	"fmt"
	"kids-score/config"
	"kids-score/controllers"
	"kids-score/database"
	"kids-score/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Configurations from config.json using Viper
	config.LoadAppConfig()
	// Initialize Database
	database.Connect(buildDBConfig())
	database.Migrate()
	// Initialize Router
	router := initRouter()
	var routerString string = fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	router.Run(routerString)
}

func buildDBConfig() database.ConnectionInfo {
	var DBconfig database.ConnectionInfo = database.ConnectionInfo{}
	DBconfig.Host = config.AppConfig.Database.Host
	DBconfig.Dbname = config.AppConfig.Database.Dbname
	DBconfig.Password = config.AppConfig.Database.Password
	DBconfig.Port = config.AppConfig.Database.Port
	DBconfig.User = config.AppConfig.Database.User
	return DBconfig
}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/token", controllers.GenerateToken)
		api.POST("/user/register", controllers.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/users", controllers.GetAllUsers)
			secured.GET("/user/:userid", controllers.GetUser)
			secured.POST("/user/:userid/password", controllers.UpdateUserPassword)
			secured.GET("/ping", controllers.Ping)
		}
	}
	return router
}
