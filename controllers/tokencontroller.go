package controllers

import (
	"kids-score/database"
	helper "kids-score/helpers"
	"kids-score/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// check if email exists and password is correct
	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		if !user.Enabled {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "User Locked"})
			context.Abort()
			return
		}
		user.IncreaseFailedLogin()
		if user.FailedLogin > 3 {
			user.Enabled = false
			database.Instance.Save(user)
			context.JSON(http.StatusUnauthorized, gin.H{"error": "User Locked"})
			context.Abort()
			return
		}
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}
	user.ResetFailedLogin()
	user.SetLastLogin()
	database.Instance.Save(user)
	tokenString, err := helper.GenerateJWT(user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
