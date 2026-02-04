// Package models defines the data structures used throughout the FitClassMaster application.
package models

import "time"

// Role represents the user's authorization level.
type Role string

const (
	RoleMember  Role = "member"
	RoleTrainer Role = "trainer"
	RoleAdmin   Role = "admin"
)

// User represents a registered person in the system.
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id"`
	FirstName string    `gorm:"size:64"`
	LastName  string    `gorm:"size:64"`
	Email     string    `gorm:"size:255;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"` // Hashed password
	Role      Role      `gorm:"size:20;default:'member'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// FullName returns the concatenated first and last name of the user.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsMember returns true if the user has the 'member' role.
func (u *User) IsMember() bool { return u.Role == RoleMember }

// IsAdmin returns true if the user has the 'admin' role.
func (u *User) IsAdmin() bool { return u.Role == RoleAdmin }

// IsTrainer returns true if the user has the 'trainer' role.
func (u *User) IsTrainer() bool { return u.Role == RoleTrainer }
