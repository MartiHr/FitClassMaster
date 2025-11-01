package models

import "time"

type Role string

const (
	RoleMember  Role = "member"
	RoleTrainer Role = "trainer"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"size:64"`
	LastName  string `gorm:"size:64"`
	Email     string `gorm:"size:255";uniqueIndex;not null`
	Password  string `gorm:"size:255";not null`
	Role      Role   `gorm:"size:20;default:'member'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
