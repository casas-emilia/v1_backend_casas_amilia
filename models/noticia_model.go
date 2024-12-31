package models

import "time"

type Noticia struct {
	ID                uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	TituloNoticia     string     `gorm:"column:titulo_noticia" json:"titulo_noticia"`
	DesarrolloNoticia string     `gorm:"column:desarrollo_noticia" json:"desarrollo_noticia"`
	// UsuarioID         uint             `gorm:"column:usuario_id" json:"usuario_id"`
	EmpresaID uint `gorm:"column:empresa_id" json:"empresa_id"`
	// Usuario           Usuario          `gorm:"foreignKey:UsuarioID"`
	Empresa        Empresa          `gorm:"foreignKey:EmpresaID"`
	Imagen_noticia []Imagen_noticia `gorm:"foreignKey:NoticiaID;constraint:OnDelete:CASCADE"`
}

func (Noticia) TableName() string {
	return "noticias"
}
