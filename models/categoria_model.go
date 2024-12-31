package models

import "time"

type Categoria struct {
	ID                   uint             `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt            time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt            *time.Time       `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombreCategoria      string           `gorm:"column:nombre_categoria" json:"nombre_categoria"`
	DescripcionCategoria string           `gorm:"column:descripcion_categoria" json:"descripcion_categoria"`
	Tipo_categoria       []Tipo_categoria `gorm:"foreignKey:CategoriaID;contraint:OnDelete:CASCADE"`
	Prefabricada         []Prefabricada   `gorm:"foreignKey:CategoriaID;constraint:OnDelete:CASCADE"`
}

func (Categoria) TableName() string {
	return "categorias"
}
