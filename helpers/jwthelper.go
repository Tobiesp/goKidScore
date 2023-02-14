package helper

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"kids-score/config"
	"kids-score/database"
	"kids-score/models"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func getTimeoutPeriod() int64 {
	var defaultTimeout int64 = int64(time.Hour)
	var timeoutStr string = config.AppConfig.Jwt.Timeout
	timeout, err := strconv.ParseInt(timeoutStr, 0, 64)
	if err != nil {
		log.Println("Warning: Unable to parse value: " + timeoutStr)
		return defaultTimeout
	}
	if strings.HasSuffix(timeoutStr, "m") {
		return timeout * int64(time.Minute)
	} else if strings.HasSuffix(timeoutStr, "h") {
		return timeout * int64(time.Hour)
	} else if strings.HasSuffix(timeoutStr, "d") {
		return timeout * int64(time.Hour) * 24
	} else if strings.HasSuffix(timeoutStr, "M") {
		return timeout * int64(time.Hour) * 24 * 30
	} else if strings.HasSuffix(timeoutStr, "y") {
		return timeout * int64(time.Hour) * 24 * 365
	} else {
		log.Println("Warning: Unable to parse value: " + timeoutStr)
		return defaultTimeout
	}
}

func GenerateJWT(email string, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(time.Duration(getTimeoutPeriod()))
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func GetUserFromToken(signedToken string) (models.User, error) {
	var user models.User
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return user, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return user, err
	}
	record := database.Instance.Where("email = ?", claims.Email).First(&user)
	if record.Error != nil {
		return user, errors.New("Could not find user.")
	}
	return user, nil
}
