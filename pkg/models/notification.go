package models

import "time"

type Notification struct {
	NotificationId string    `gorm:"column:notification_id"`
	ContentType    string    `gorm:"column:content_type"`
	SentType       string    `gorm:"column:sent_type"`
	AddrsStr       string    `gorm:"column:addrs_str"`
	Title          string    `gorm:"column:title"`
	Content        string    `gorm:"column:content"`
	ShortContent   string    `gorm:"column:short_content"`
	ExporedDays    int64     `gorm:"column:expired_days"`
	Owner          string    `gorm:"column:owner"`
	Status         string    `gorm:"column:status"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	DeletedAt      time.Time `gorm:"column:deleted_at;default:null"`
}
