package models

import "time"

type Imagen_prefabricada struct {
	ID             uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	Image          string     `gorm:"column:image" json:"image"`
	PrefabricadaID uint       `gorm:"column:prefabricada_id" json:"prefabricada_id"`
}

func (Imagen_prefabricada) TableName() string {
	return "imagenes_prefabricadas"
}
