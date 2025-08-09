package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID   uuid.UUID `gorm:"type:char(36);not null;index" json:"product_id"`
	ImageObject string    `gorm:"type:varchar(255);not null" json:"image_object"`
	Position    int       `gorm:"type:int;default:0" json:"position"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;autoUpdateTime" json:"updated_at"`
	Product     Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
