package models

import "time"

type Contacto struct {
	ID               uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	EmailLaboral     string     `gorm:"column:email_laboral" json:"email_laboral"`
	CelularLaboral   string     `gorm:"column:celular_laboral" json:"celular_laboral"`
	DireccionLaboral string     `gorm:"column:direccion_laboral" json:"direccion_laboral"`
	UsuarioID        uint       `gorm:"usuario_id" json:"usuario_id"`
	Usuario          Usuario    `gorm:"foreignKey:UsuarioID"`
}

func (Contacto) TableName() string {
	return "contactos"
}
