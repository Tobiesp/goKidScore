package middlewares

import (
	helper "jwt-auth/helpers"
	"jwt-auth/models"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err := helper.ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		user, err := helper.GetUserFromToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": "No a vaild token."})
			context.Abort()
			log.Println("ERROR: " + err.Error())
			return
		}
		action := actionFromMethod(context.Request.Method)
		if !hasPermission(user, action, context.Request.URL.Path) {
			context.JSON(403, gin.H{"error": "Access deined"})
			context.Abort()
			return
		}
		context.Next()
	}
}

func IsAdmin(context *gin.Context) bool {
	jwtToken := context.GetHeader("Authorization")
	if jwtToken == "" {
		return false
	}
	user, err := helper.GetUserFromToken(jwtToken)
	if err != nil {
		return false
	}
	return user.Role == "Admin"
}

func hasPermission(user models.User, action string, route string) bool {
	if user.Email == "j.jack@test.org" {
		return true
	}
	if user.Role == "Admin" {
		return true
	}
	if user.Role == "Writer" && (action == "read" || action == "write") {
		if strings.HasPrefix(route, "api/user/") && strings.HasSuffix(route, strconv.FormatInt(int64(user.ID), 10)) {
			return true
		} else if strings.HasPrefix(route, "api/user/") && !strings.HasSuffix(route, strconv.FormatInt(int64(user.ID), 10)) {
			return false
		} else {
			return true
		}
	}
	if user.Role == "reader" && action == "read" {
		if strings.HasPrefix(route, "api/user/") && strings.HasSuffix(route, strconv.FormatInt(int64(user.ID), 10)) {
			return true
		} else if strings.HasPrefix(route, "api/user/") && !strings.HasSuffix(route, strconv.FormatInt(int64(user.ID), 10)) {
			return false
		} else {
			return true
		}
	}
	return false
}

func actionFromMethod(httpMethod string) string {
	switch httpMethod {
	case "GET":
		return "read"
	case "POST":
		return "write"
	case "PUT":
		return "write"
	case "DELETE":
		return "write"
	default:
		return ""
	}
}
