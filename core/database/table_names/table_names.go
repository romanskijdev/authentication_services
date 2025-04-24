package dbcoretablenames

type TableName string

const (
	TableNameUsers        TableName = "users" // Пользователи
	TableNameNotification TableName = "notifications"
)

func (t TableName) ToString() string {
	return string(t)
}
