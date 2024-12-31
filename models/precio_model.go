package models

import "time"

type Precio struct {
	ID                uint         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt         time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time    `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         *time.Time   `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombrePrecio      string       `gorm:"column:nombre_precio" json:"nombre_precio"`
	DescripcionPrecio string       `gorm:"column:descripcion_precio" json:"descripcion_precio"`
	ValorPrefabricada float64      `gorm:"column:valor_prefabricada" json:"valor_prefabricada"` // Usamos float64 para soportar decimales
	PrefabricadaID    uint         `gorm:"column:prefabricada_id" json:"prefabricada_id"`
	Prefabricada      Prefabricada `gorm:"foreignKey:PrefabricadaID"`
	Incluye           []Incluye    `gorm:"foreignKey:PrecioID;constraint:OnDelete:CASCADE"`
}

func (Precio) TableName() string {
	return "precios"
}
