package models

import "time"

// 基本模型的定义
type NFModel struct {
	status string `gorm:"column:status"`
	createTime time.Time `gorm:"column:create_time"`
	statusTime time.Time `gorm:"column:status_time"`
}

