package models

import "time"

type Tipo_categoria struct {
	ID          uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	CategoriaID uint       `gorm:"column:categoria_id" json:"categoria_id"`
	TipoID      uint       `gorm:"column:tipo_id" json:"tipo_id"`
	Categoria   Categoria  `gorm:"foreignKey:CategoriaID"`
	Tipo        Tipo       `gorm:"foreignKey:TipoID"`
}

func (Tipo_categoria) TableName() string {
	return "tipos_categorias"
}
