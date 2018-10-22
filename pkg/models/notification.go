package models

import "time"

// User NotificationCenterPost.
type NotificationCenterPost struct {
	NfPostID string `gorm:"column:nf_post_id"`
	NfPostType string `gorm:"column:nf_post_type"`
	Title string `gorm:"column:title"`
	Content string `gorm:"column:content"`
	ShortContent string `gorm:"column:short_content"`
	ExporedDays int64 `gorm:"column:expired_days"`
	Owner string `gorm:"column:owner"`
	//NFModel
	Status string `gorm:"column:status"`
	CreateTime time.Time `gorm:"column:create_time"`
	StatusTime time.Time `gorm:"column:status_time"`
}


func NewNotificationCenterPost(content string) *NotificationCenterPost {
	return &NotificationCenterPost{
		RuntimeCredentialId: NewRuntimeCrentialId(),
		Content:             content,
		CreateTime:          time.Now(),
	}
}
