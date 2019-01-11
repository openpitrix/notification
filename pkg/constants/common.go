package constants

const (
	EtcdPrefix = "notification/"
	EmailQueue = "emailtask"
)

const (
	NfPostIDPrifix   = "nf-"
	JobPostIDPrifix  = "job-"
	TaskPostIDPrifix = "task-"
)

const (
	StatusNew      = "new"
	StatusSending  = "sending"
	StatusFinished = "finished"
	StatusFailed   = "failed"
)

const (
	ServiceName = "Notification"
	prefix      = "notification-"
	//ApiGatewayHost = prefix + "api-gateway"
	ApiGatewayHost = "127.0.0.1"
	ApiGatewayPort = 9200

	//NotificationManagerHost = prefix + "manager"
	NotificationManagerHost = "127.0.0.1"
	NotificationManagerPort = 9201
)
