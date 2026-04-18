package model

import "time"

type AdminRole int16

const (
	AdminRoleUnspecified AdminRole = 0
	AdminRoleSuperAdmin AdminRole = 1
	AdminRoleAdmin      AdminRole = 2
	AdminRoleReadonly   AdminRole = 3
)

// AdminUser admin account
type AdminUser struct {
	ID           string     `gorm:"column:id;primaryKey"`
	Username     string     `gorm:"column:username;uniqueIndex;not null"`
	PasswordHash string     `gorm:"column:password_hash;not null"`
	Email        string     `gorm:"column:email"`
	Role         AdminRole  `gorm:"column:role;type:smallint;not null;default:2"`
	Status       int8       `gorm:"column:status;not null;default:1"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (AdminUser) TableName() string { return "admin_users" }

// IsActive returns whether the admin account is active
func (a *AdminUser) IsActive() bool { return a.Status == 1 }
