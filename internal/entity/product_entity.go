package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Category    datatypes.JSON `gorm:"type:json" json:"category"`
	Brand       string         `gorm:"type:varchar(100);not null" json:"brand"`
	Color       datatypes.JSON `gorm:"type:json" json:"color"`
	Specs       datatypes.JSON `gorm:"type:json" json:"specs"`
	Price       float64        `gorm:"type:decimal(12,2);not null" json:"price"`
	CreatedBy   uuid.UUID      `gorm:"type:char(36);not null;column:created_by" json:"created_by"`
	CreatedAt   time.Time      `gorm:"type:timestamp;not null;column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;not null;column:updated_at;autoUpdateTime" json:"updated_at"`
	Images      []ProductImage `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"images"`
}

func (Product) TableName() string {
	return "products"
}
