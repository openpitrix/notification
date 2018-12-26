package models

import "time"

type Task struct {
	TaskID     string    `gorm:"column:task_id"`
	JobID      string    `gorm:"column:job_id"`
	EmailAddr  string    `gorm:"column:email_addr"`
	TaskAction string    `gorm:"column:task_action"`
	ErrorCode  int64     `gorm:"column:error_code"`
	Status     string    `gorm:"column:status"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

//func (Task) TableName() string {
//	return "task"
//}

type TaskWNfInfo struct {
	Title        string
	ShortContent string
	Content      string
	TaskID       string
	AddrsStr     string
}
