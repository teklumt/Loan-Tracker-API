package domain

type SystemLog struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Details   string `json:"details"`
}

type LogRepository interface {
	GetSystemLogs(Type, limit, page string) ([]SystemLog, error)
	CreateLog(log SystemLog) error
}

type LogService interface {
	GetSystemLogs(Type, limit, page string) ([]SystemLog, ErrorResponse)
	CreateLog(log SystemLog) ErrorResponse
}