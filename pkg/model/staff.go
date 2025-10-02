package model

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username   string    `gorm:"column:username;unique;not null" json:"username"`
	Password   string    `gorm:"column:password;not null" json:"-"` 
	Fullname   string    `gorm:"column:fullname" json:"fullname"` 
	Email      string    `gorm:"column:email;unique;not null" json:"email"`
	Role       string    `gorm:"column:role;not null" json:"role"`
	CreateAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	BusinessID uuid.UUID `gorm:"column:business_id" json:"business_id"`
}

type StaffRequest struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password" binding:"required,min=6"`
	Fullname   string    `json:"fullname" binding:"required"`
	Email      string    `json:"email" binding:"required,email"`
	Role       string    `json:"role" binding:"required"`
	BusinessID uuid.UUID `json:"business_id" binding:"required"`
}

type StaffUpdateRequest struct {
	Username   string    `json:"username"`
	Email      string    `json:"email" binding:"omitempty,email"`
	Role       string    `json:"role"`
	BusinessID uuid.UUID `json:"business_id"`
}

type GetListStaffRequest struct {
	Username   *string `json:"username,omitempty" form:"username"`
	Email      *string `json:"email,omitempty" form:"email"`
	Role       *string `json:"role,omitempty" form:"role"`
	BusinessID *string `json:"business_id,omitempty" form:"business_id"`
	Page       int     `json:"page" form:"page"`
	PageSize   int     `json:"page_size" form:"page_size"`
	Sort       string  `json:"sort" form:"sort"`
	Keyword     string  `json:"keyword" form:"keyword"`
}

type GetListStaffResponse struct {
	Data []Staff                `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Staff Staff  `json:"staff"`
}
