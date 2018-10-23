package models

import "time"

// User NotificationCenterPost.
type Task struct {
	TaskID string `gorm:"column:task_id"`
	JobID string `gorm:"column:job_id"`
	AddrsStr string `gorm:"column:addrs_str"`
	TaskAction string `gorm:"column:task_action"`
	Result string `gorm:"column:result"`
	ErrorCode string `gorm:"column:error_code"`
	Status int64 `gorm:"column:status"`
	//NFModel
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}