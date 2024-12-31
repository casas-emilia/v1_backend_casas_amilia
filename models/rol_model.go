package models

import "time"

type Rol struct {
	ID             uint          `gomr:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt      time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time    `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombreRol      string        `gorm:"column:nombre_rol" json:"nombre_rol"`
	DescripcionRol string        `gorm:"column:descripcion_rol" json:"descripcion_rol"`
	Rol_usuario    []Rol_usuario `gorm:"foreignKey:RolID;constraint:OnDelete:CASCADE"` // Eliminar Rol_usuario si se elimina el rol
}

func (Rol) TableName() string {
	return "roles"
}
