package model

import (
	"time"

	"github.com/google/uuid"
)

type Business struct {
	ID          uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name"`
	Description string    `gorm:"type:text"`
	Address     string    `json:"address"`
	BusinessType        string    `json:"type"`
	Status      string    `json:"status"`
	CreateAt    time.Time `gorm:"column:created_at"`
	Staffs []Staff `gorm:"foreignKey:BusinessID"`
	WorkerName string `json:"woker_name"`
}

type BusinessRequest struct {
	ID          uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string  `json:"name"`
	Description string `gorm:"type:text"`
	Address     string `json:"address"`
	BusinessType string  `json:"type"`
	Status      string `json:"status"`
}

type UriParse struct {
	ID []string `json:"id" uri:"id"`
}

type GetListBusinessRequest struct {
	Name        *string `json:"name,omitempty" form:"name"`
	Page        int     `json:"page" form:"page" form:"page"`
	PageSize    int     `json:"page_size" form:"page_size"`
	Sort        string  `json:"sort" form:"sort"`
	ManagerID   *string `json:"manager_id" form:"manager_id"`
	Address     *string `json:"address" form:"address"`
	Type        *string `json:"type" form:"type"`
	Description *string `json:"description" form:"description"`
}

type GetListBusinessResponse struct {
	Data []Business             `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}
