package typescore

import "time"

type ParamsSendMail struct {
	From    string
	To      string
	Subject string
}

type NotifyCategory string

const (
	InfoNotifyCategory      NotifyCategory = "info"   // Информационное уведомление(от админа)
	DeviceNewNotifyCategory NotifyCategory = "ip_new" // Новый IP
)

type NotifyParams struct {
	Text         *string         // Текст уведомления
	Title        *string         // Заголовок уведомления
	ImageURLPath *string         // Ссылка на картинку в уведомлении
	UsersIDs     []*string       // Список идентификаторов пользователей
	Category     *NotifyCategory // Категория уведомления

	IsTelegram bool // Отправка уведомления в Telegram
	IsEmail    bool // Отправка уведомления на email

	Emergency bool // Экстренная отправка уведомления игнорирует запреты пользователя
}

// Notification - структура для хранения данных о системных уведомлениях.
type Notification struct {
	ID                   *uint64    `gorm:"column:id;index;type:bigserial;autoIncrement;not null;unique"  ignore_update_db:"true" json:"id" db:"id" mapstructure:"id"`                                             // Автоматически заполняемое уникальное поле для большого числа записей
	UniqUUID             *string    `gorm:"type:uuid;primaryKey;column:uniq_uuid;default:gen_random_uuid();not null" db:"uniq_uuid"  ignore_update_db:"true"  json:"uniq_uuid,omitempty" mapstructure:"uniq_uuid"` // Уникальный идентификатор уведомления
	Body                 *string    `gorm:"type:text;column:body" json:"body" db:"body"`                                                                                                                           // Информационное тело уведомления
	Title                *string    `gorm:"type:text;column:title" json:"title" db:"title"`                                                                                                                        // Информационный заголовок уведомления
	ImageURL             *string    `gorm:"type:text;column:image_url" json:"image_url" db:"image_url"`                                                                                                            // Ссылка на картинку в уведомлении
	DocsURL              *string    `gorm:"type:text;column:docs_url" json:"docs_url" db:"docs_url"`                                                                                                               // Ссылка на документы в уведомлении
	TypesNotif           *string    `gorm:"type:text;index;column:types_notif" json:"types_notif" db:"types_notif" mapstructure:"types_notif"`                                                                     // Тип уведомления
	IsTelegram           *bool      `gorm:"default:false;column:is_telegram" json:"is_telegram,omitempty" db:"is_telegram" mapstructure:"is_telegram"`                                                             // Отправка уведомления в Telegram
	IsEmail              *bool      `gorm:"default:false;column:is_email" json:"is_email,omitempty" db:"is_email" mapstructure:"is_email"`                                                                         // Отправка уведомления на email
	SuccessfullyReceived []*string  `gorm:"type:text[];column:successfully_received" json:"successfully_received" db:"successfully_received" mapstructure:"successfully_received"`                                 // Массив успешных получателей уведомления
	ErrorReceived        []*string  `gorm:"type:text[];column:error_received" json:"error_received" db:"error_received"`                                                                                           // Массив неуспешных получателей уведомления
	SendDate             *time.Time `gorm:"type:timestamp;column:send_date;default:CURRENT_TIMESTAMP" json:"send_date" db:"send_date"`                                                                             // Дата отправки уведомления пользователям
}
