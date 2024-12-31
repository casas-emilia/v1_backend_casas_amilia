package models

import "time"

type Portada struct {
	ID            uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombrePortada string     `gorm:"column:nombre_portada" json:"nombre_portada"`
	Image         string     `gorm:"column:image" json:"image"`
	EmpresaID     uint       `gorm:"column:empresa_id" json:"empresa_id"`
	Empresa       Empresa    `gorm:"foreignKey:EmpresaID"`
}

func (Portada) TableName() string {
	return "portadas"
}
