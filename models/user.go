package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserView struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AdminUserView struct {
	ID          int64     `json:"ID"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Enabled     bool      `json:"enabled"`
	LastLogon   time.Time `json:"lastLogon"`
	FailedLogin int       `json:"failedLogonCount"`
}

type User struct {
	gorm.Model
	Name        string    `json:"name"`
	Username    string    `json:"username" gorm:"unique"`
	Email       string    `json:"email" gorm:"unique"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`
	Enabled     bool      `json:"enabled"`
	LastLogon   time.Time `json:"lastLogon"`
	FailedLogin int       `json:"failedLogonCount"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) IncreaseFailedLogin() {
	user.FailedLogin++
}

func (user *User) ResetFailedLogin() {
	user.FailedLogin = 0
}

func (user *User) SetLastLogin() {
	user.LastLogon = time.Now()
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
