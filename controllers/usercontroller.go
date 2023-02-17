package controllers

import (
	"kids-score/database"
	helper "kids-score/helpers"
	"kids-score/middlewares"
	"kids-score/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	if UserExists(user.Username) {
		context.JSON(http.StatusConflict, gin.H{"error": "user already registered"})
		context.Abort()
		return
	}
	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}

func GetAllUsers(context *gin.Context) {
	if middlewares.IsAdmin(context) {
		var users []models.User
		database.Instance.Find(&users)
		usersview := transformAdminUser(users)
		context.IndentedJSON(http.StatusOK, usersview)
	} else {
		context.JSON(http.StatusForbidden, gin.H{"error": "request not allowed"})
		context.Abort()
	}
}

func GetUser(context *gin.Context) {
	userid := context.Param("userid")
	if userid != "" {
		var user models.User
		err := database.Instance.Where("id = ?", userid).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				context.Abort()
				return
			}
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.IndentedJSON(http.StatusOK, transformUser(user))
	} else {
		context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		context.Abort()
		return
	}
}

func UpdateUserPassword(context *gin.Context) {
	userid := context.Param("userid")
	user, err := helper.GetUserFromToken(context.GetHeader("Authorization"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if userid != strconv.FormatInt(int64(user.ID), 10) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "action can't be preformed for user"})
		context.Abort()
		return
	}
	var userPasswordView models.UserPaswordView
	if err := context.ShouldBindJSON(&userPasswordView); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	passHash, err := helper.HashPassword(userPasswordView.OldPassword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if passHash != user.Password {
		context.JSON(http.StatusBadRequest, gin.H{"error": "not valid credentials"})
		context.Abort()
		return
	}
	if userPasswordView.NewPassword != userPasswordView.RepeatPassword {
		context.JSON(http.StatusBadRequest, gin.H{"error": "new password and repeat password must match"})
		context.Abort()
		return
	}
	if err := user.HashPassword(userPasswordView.NewPassword); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := database.Instance.Save(user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
}

func UserExists(username string) bool {
	var user models.User
	err := database.Instance.Where("id = ?", username).First(&user).Error
	return err != gorm.ErrRecordNotFound
}

// func transformUserList(users []models.User) []models.UserView {
// 	var usersView []models.UserView = make([]models.UserView, len(users))
// 	for i := range users {
// 		usersView[i].ID = int64(users[i].ID)
// 		usersView[i].Username = users[i].Username
// 		usersView[i].Name = users[i].Name
// 		usersView[i].Email = users[i].Email
// 	}
// 	return usersView
// }

func transformUser(user models.User) models.UserView {
	var usersView models.UserView
	usersView.ID = int64(user.ID)
	usersView.Username = user.Username
	usersView.Name = user.Name
	usersView.Email = user.Email
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
