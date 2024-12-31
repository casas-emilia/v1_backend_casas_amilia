package models

import "time"

type Usuario struct {
	ID              uint        `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt       time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       *time.Time  `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	PrimerNombre    string      `gorm:"column:primer_nombre" json:"primer_nombre"`
	SegundoNombre   string      `gorm:"column:segundo_nombre" json:"segundo_nombre"`
	PrimerApellido  string      `gorm:"column:primer_apellido" json:"primer_apellido"`
	SegundoApellido string      `gorm:"column:segundo_apellido" json:"segundo_apellido"`
	Image           string      `gorm:"column:image" json:"image"`
	EmpresaID       uint        `gorm:"column:empresa_id" json:"empresa_id"`
	Empresa         Empresa     `gorm:"foreignKey:EmpresaID"`
	Contacto        []Contacto  `gorm:"foreignKey:UsuarioID;constraint:OnDelete:CASCADE"`
	Credencial      *Credencial `gorm:"foreignKey:UsuarioID;constraint:OnDelete:CASCADE"`
	// Noticia         []Noticia     `gorm:"foreignKey:UsuarioID;constraint:OnDelete:CASCADE"`
	Rol_usuario []Rol_usuario `gorm:"foreignKey:UsuarioID;constraint:OnDelete:CASCADE"`
}

func (Usuario) TableName() string {
	return "usuarios"
}
