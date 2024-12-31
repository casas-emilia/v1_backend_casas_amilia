package models

import "time"

type Rol_usuario struct {
	ID        uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	UsuarioID uint       `gorm:"column:usuario_id" json:"usuario_id"`
	RolID     uint       `gorm:"column:rol_id" json:"rol_id"`
	Usuario   Usuario    `gorm:"foreignKey:UsuarioID"`
	Rol       Rol        `gorm:"foreignKey:RolID"`
}

func (Rol_usuario) TableName() string {
	return "roles_usuarios"
}
