package models

import "time"

type Estilo struct {
	ID                uint           `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         *time.Time     `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombreEstilo      string         `gorm:"column:nombre_estilo" json:"nombre_estilo"`
	DescripcionEstilo string         `gorm:"column:descripcion_estilo" json:"descripcion_estilo"`
	Prefabricada      []Prefabricada `gorm:"foreignKey:EstiloID;constraint:OnDelete:CASCADE"`
}

func (Estilo) TableName() string {
	return "estilos"
}
