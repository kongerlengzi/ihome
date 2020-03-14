package model

import (
	"08_ihome/apiv1.0/handler"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name string           `json:"name" gorm:"not null"`
	Password_hash string  `json:"password_hash" gorm:"not null"`
	Mobile string         `json:"mobile" gorm:"not null"`
	Real_name string      `json:"real_name"`
	Id_card string        `json:"id_card"`
	Avatar_url string     `json:"name"`
	Houses []Houses       `gorm:"foreignkey:UserID"`
	Orders []Orders
}

func ExistUserByID(mobile string) bool {
	var user User
	db.Select("mobile").Where("mobile=?",mobile).First(&user)
	if user.Mobile != "" {
		return true
	}
	return false
}

func AddUser(u handler.Register) bool {
	user := User {
		Name : u.Name,
		Password_hash : u.Pwd,
		Mobile : u.Mobile,
	}
	db.Create(&user)
	return true
}

func GetUser(mobile string) (user User) {
	db.Where("mobile = ?", mobile).First(&user)
	return
}

func Update(id int, data interface{}) bool {
	var user User
	db.Model(&user).Where("id = ?",id).Update(data)
	return true
}

func Login(mobile string, pwd string) bool {
	user := GetUser(mobile)
	if user.Password_hash != pwd {
		return false
	}
	return true
}