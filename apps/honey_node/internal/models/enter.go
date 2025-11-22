package models

// File: models/enter.go
// Description: 内嵌模型

import (
	"time"
)

// Model 内嵌模型
type Model struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
