package models

import "time"

type Prefabricada struct {
	ID                  uint                  `gorm:"primarykey;autoIncrement;column:id" json:"id"`
	CreatedAt           time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time             `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt           *time.Time            `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombrePrefabricada  string                `gorm:"column:nombre_prefabricada" json:"nombre_prefabricada"`
	M2                  int                   `gorm:"column:m2" json:"m2"`
	Garantia            string                `gorm:"column:garantia" json:"garantia"`
	Eslogan             string                `gorm:"column:eslogan" json:"eslogan"`
	Descripcion         string                `gorm:"column:descripcion" json:"descripcion"`
	Destacada           bool                  `gorm:"column:destacada;default:false" json:"destacada"`
	Oferta              bool                  `gorm:"column:oferta;default:false" json:"oferta"`
	CategoriaID         uint                  `gorm:"column:categoria_id" json:"categoria_id"`
	EmpresaID           uint                  `gorm:"column:empresa_id" json:"empresa_id"`
	EstiloID            uint                  `gorm:"column:estilo_id" json:"estilo_id"`
	TipoID              uint                  `gorm:"column:tipo_id" json:"tipo_id"`
	Categoria           Categoria             `gorm:"foreignKey:CategoriaID"`
	Empresa             Empresa               `gorm:"foreignKey:EmpresaID"`
	Estilo              Estilo                `gorm:"foreignKey:EstiloID"`
	Tipo                Tipo                  `gorm:"foreignKey:TipoID"`
	Imagen_prefabricada []Imagen_prefabricada `gorm:"foreignKey:PrefabricadaID;constraint:OnDelete:CASCADE"`
	Caracteristica      []Caracteristica      `gorm:"foreignKey:PrefabricadaID;constraint:OnDelete:CASCADE"`
	Precio              []Precio              `gorm:"foreignKey:PrefabricadaID;constraint:OnDelete:CASCADE"`
}

func (Prefabricada) TableName() string {
	return "prefabricadas"
}
