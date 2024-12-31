package models

import "time"

type Tipo struct {
	ID                  uint             `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt           time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt           *time.Time       `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	MaterialEstructura  string           `gorm:"column:material_estructura" json:"material_estructura"`
	DescripcionMaterial string           `gorm:"column:descripcion_material" json:"descripcion_material"`
	Tipo_categoria      []Tipo_categoria `gorm:"foreignKey:TipoID;constraint:OnDelete:CASCADE"`
	Prefabricada        []Prefabricada   `gorm:"foreignKey:TipoID;constraint:OnDelete:CASCADE"`
}

func (Tipo) TableName() string {
	return "tipos"
}
