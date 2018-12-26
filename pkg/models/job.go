package models

import "time"

type Job struct {
	JobID          string    `gorm:"column:job_id"`
	NotificationId string    `gorm:"column:notification_id"`
	JobType        string    `gorm:"column:job_type"`
	AddrsStr       string    `gorm:"column:addrs_str"`
	JobAction      string    `gorm:"column:job_action"`
	ExeCondition   string    `gorm:"column:exe_condition"`
	TotalTaskCount int64     `gorm:"column:total_task_count"`
	TaskSuccCount  int64     `gorm:"column:task_succ_count"`
	ErrorCode      int64     `gorm:"column:error_code"`
	Status         string    `gorm:"column:status"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}
