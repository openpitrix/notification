package models

import "time"

// User NotificationCenterPost.
type NotificationCenterPost struct {
	NfPostID string `gorm:"column:nf_post_id"`
	NfPostType string `gorm:"column:nf_post_type"`
	AddrsStr string `gorm:"column:addrs_str"`
	Title string `gorm:"column:title"`
	Content string `gorm:"column:content"`
	ShortContent string `gorm:"column:short_content"`
	ExporedDays int64 `gorm:"column:expired_days"`
	Owner string `gorm:"column:owner"`
	Status string `gorm:"column:status"`
	//NFModel
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}


