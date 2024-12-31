package models

import "time"

type Empresa struct {
	ID                 uint           `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt          *time.Time     `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	NombreEmpresa      string         `gorm:"column:nombre_empresa" json:"nombre_empresa"`
	DescripcionEmpresa string         `gorm:"column:descripcion_empresa" json:"descripcion_empresa"`
	HistoriaEmpresa    string         `gorm:"column:historia_empresa" json:"historia_empresa"`
	MisionEmpresa      string         `gorm:"column:mision_empresa" json:"mision_empresa"`
	VisionEmpresa      string         `gorm:"column:vision_empresa" json:"vision_empresa"`
	UbicacionEmpresa   string         `gorm:"column:ubicacion_empresa" json:"ubicacion_empresa"`
	CelularEmpresa     string         `gorm:"column:celular_empresa" json:"celular_empresa"`
	EmailEmpresa       string         `gorm:"not null;column:email" json:"email"`
	Usuario            []Usuario      `gorm:"foreignKey:EmpresaID;constraint:OnDelete:CASCADE"`
	Red                []Red          `gorm:"foreignKey:EmpresaID;constraint:OnDelete:CASCADE"`
	Servicio           []Servicio     `gorm:"foreignKey:EmpresaID;constraint:OnDelete:CASCADE"`
	Portada            []Portada      `gorm:"foreignKey:EmpresaID;constraint:OnDelete:CASCADE"`
	Prefabricada       []Prefabricada `gorm:"foreignKey:EmpresaID;constraint:OnDelete:CASCADE"`
	Noticia            []Noticia      `gorm:"foreignKey:EmpresaID;constraint:OnDelete:CASCADE"`
}

func (Empresa) TableName() string {
	return "empresa"
}
