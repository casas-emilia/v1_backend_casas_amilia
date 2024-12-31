package models

import "time"

type Incluye struct {
	ID            uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombreIncluye string     `gorm:"column:nombre_incluye" json:"nombre_incluye"`
	PrecioID      uint       `gorm:"column:precio_id" json:"precio_id"`
	Precio        Precio     `gorm:"foreignKey:PrecioID"`
}

func (Incluye) TableName() string {
	return "incluyes"
}
