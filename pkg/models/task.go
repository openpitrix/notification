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

type TaskWithNfInfo struct {
	NotificationId string
	JobID          string
	TaskID         string
	Title          string
	ShortContent   string
	Content        string
	EmailAddr      string
}

const (
	GetTaskwithNfContentbyIDSQL = "SELECT  t3.notification_id,t2.job_id,t1.task_id,t3.title,t3.short_content,  t3.content,t1.email_addr " +
		"	FROM task t1,job t2,notification t3 where t1.job_id=t2.job_id and t2.notification_id=t3.notification_id  and t1.task_id=? "
)
