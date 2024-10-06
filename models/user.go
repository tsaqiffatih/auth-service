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
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Role     string `json:"role"`
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
