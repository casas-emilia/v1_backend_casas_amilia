package models

import "time"

type Servicio struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt           *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombreServicio      string     `gorm:"column:nombre_servicio" json:"nombre_servicio"`
	DescripcionServicio string     `gorm:"column:descripcion_servicio" json:"descripcion_servicio"`
	EmpresaID           uint       `gorm:"column:empresa_id" json:"empresa_id"`
	Empresa             Empresa    `gorm:"foreignKey:EmpresaID"`
}

func (Servicio) TableName() string {
	return "servicios"
}
