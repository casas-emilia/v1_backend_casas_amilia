package models

import "time"

type Caracteristica struct {
	ID             uint         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt      time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time   `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	Clave          string       `gorm:"column:clave" json:"clave"`
	Valor          string       `gorm:"column:valor" json:"valor"`
	PrefabricadaID uint         `gorm:"column:prefabricada_id" json:"prefabricada_id"`
	Prefabricada   Prefabricada `gorm:"foreignKey:PrefabricadaID"`
}

func (Caracteristica) TableName() string {
	return "caracteristicas"
}
