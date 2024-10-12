package models

import (
	"errors"

	"gorm.io/gorm"
)

const (
	RoleOwner   = "owner"
	RoleCashier = "cashier"
	RoleAdmin   = "admin"
)

var validRoles = []string{RoleOwner, RoleCashier, RoleAdmin}

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null" validate:"required"`
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string `json:"password" gorm:"not null" validate:"required,min=6"`
	Role     string `json:"role" gorm:"not null" validate:"required,oneof=owner cashier admin user"`
}

// func for validate user role
func (u *User) ValidateRole() error {
	for _, validRole := range validRoles {
		if u.Role == validRole {
			return nil
		}
	}
	return errors.New("invalid role: must be one of [owner, cashier, admin]")
}

type Token struct {
	gorm.Model
	UserID uint
	Token  string
}
