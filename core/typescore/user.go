package typescore

import "time"

// user Role - роли пользователей
type UserRoleTypes string

const (
	UserRole       UserRoleTypes = "user"        // пользователь
	AdminRole      UserRoleTypes = "admin"       // администратор
	SuperAdminRole UserRoleTypes = "super_admin" // супер администратор
	SupportRole    UserRoleTypes = "support"     // поддержка
)

// User - структура для управления пользователями + данные пользователя
type User struct {
	SystemID            *string        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:system_id" json:"system_id,omitempty" db:"system_id" mapstructure:"system_id"` // Системный идентификатор записи
	SerialID            *uint64        `gorm:"type:bigint;index;autoIncrement;unique;column:serial_id" json:"serial_id" db:"serial_id" mapstructure:"serial_id"`                    // Уникальный порядковый идентификатор записи
	Role                *UserRoleTypes `gorm:"type:varchar(20);index;default:'user';column:role" json:"role" db:"role"`                                                             // Роль пользователя
	Email               *string        `gorm:"type:varchar(255);index;unique;column:email" json:"email,omitempty" db:"email" mapstructure:"email"`                                  // Адрес электронной почты пользователя
	TelegramID          *int64         `gorm:"unique;index;column:telegram_id" json:"telegram_id" db:"telegram_id"`                                                                 // Идентификатор пользователя в Telegram
	Nickname            *string        `gorm:"type:varchar(50);index;unique;column:nickname" json:"nickname,omitempty" db:"nickname" mapstructure:"nickname"`                       // Псевдоним или никнейм пользователя
	FirstName           *string        `gorm:"type:varchar(50);column:first_name" json:"first_name,omitempty" db:"first_name"`                                                      // Имя пользователя
	LastName            *string        `gorm:"type:varchar(50);column:last_name" json:"last_name,omitempty" db:"last_name"`                                                         // Фамилия пользователя
	NotificationEnabled *bool          `gorm:"default:true;column:notification_enabled" json:"notification_enabled" db:"notification_enabled"`                                      // Включены ли разрешения на push-уведомления
	IsBlocked           *bool          `gorm:"default:false;column:is_blocked" json:"is_blocked" db:"is_blocked"`                                                                   // Залочен ли пользователь(заблокирован или нет)
	CreatedAt           *time.Time     `gorm:"default:CURRENT_TIMESTAMP;column:created_at" json:"created_at" db:"created_at"`                                                       // Дата и время создания записи
}
