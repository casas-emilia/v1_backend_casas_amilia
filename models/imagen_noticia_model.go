package models

import "time"

type Imagen_noticia struct {
	ID        uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	Image     string     `gorm:"column:image" json:"image"`
	NoticiaID uint       `gorm:"column:noticia_id" json:"noticia_id"`
	Noticia   Noticia    `gorm:"foreignKey:NoticiaID"`
}

func (Imagen_noticia) TableName() string {
	return "imagenes_noticias"
}
