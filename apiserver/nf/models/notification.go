package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 基本模型的定义
type Model struct {
	ID uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}


type NotificationCenterPost struct {
	//Id int `gorm: "primary_key;AUTO_INCREMENT:number"`
	Id int `gorm: "primary_key;AUTO_INCREMENT:number"`
	UserName string `gorm:"column:name"`
	Phone string
	gorm.Model

}