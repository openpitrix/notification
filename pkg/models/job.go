package models

import "time"

// User NotificationCenterPost.
type Job struct {
	JobID string `gorm:"column:job_id"`
	NfPostID string `gorm:"column:nf_post_id"`
	JobType string `gorm:"column:job_type"`
	AddrsStr string `gorm:"column:addrs_str"`
	JobAction string `gorm:"column:job_action"`
	ExeCondition string `gorm:"column:exe_condition"`
	TotalTaskCount int64 `gorm:"column:total_task_count"`
	TaskSuccCount int64 `gorm:"column:task_succ_count"`
	Result string `gorm:"column:result"`
	ErrorCode int64 `gorm:"column:error_code"`
	Status string `gorm:"column:status"`
	//NFModel
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

