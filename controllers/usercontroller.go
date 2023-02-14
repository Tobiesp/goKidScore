package controllers

import (
	"jwt-auth/database"
	"jwt-auth/middlewares"
	"jwt-auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	user.Enabled = true
	user.FailedLogin = 0
	user.Role = ""
	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}

func GetAllUsers(context *gin.Context) {
	var users []models.User
	database.Instance.Find(&users)
	if middlewares.IsAdmin(context) {
		usersview := transformAdminUser(users)
		context.IndentedJSON(http.StatusOK, usersview)
	} else {
		usersview := transformUser(users)
		context.IndentedJSON(http.StatusOK, usersview)
	}
}

func transformUser(users []models.User) []models.UserView {
	var usersView []models.UserView = make([]models.UserView, len(users))
	for i := range users {
		usersView[i].ID = int64(users[i].ID)
		usersView[i].Username = users[i].Username
		usersView[i].Name = users[i].Name
		usersView[i].Email = users[i].Email
	}
	return usersView
}

func transformAdminUser(users []models.User) []models.AdminUserView {
	var UserView []models.AdminUserView = make([]models.AdminUserView, len(users))
	for i := range users {
		UserView[i].ID = int64(users[i].ID)
		UserView[i].Email = users[i].Email
		UserView[i].Enabled = users[i].Enabled
		UserView[i].FailedLogin = users[i].FailedLogin
		UserView[i].LastLogon = users[i].LastLogon
		UserView[i].Name = users[i].Name
		UserView[i].Role = users[i].Role
		UserView[i].Username = users[i].Username
	}
	return UserView
}
