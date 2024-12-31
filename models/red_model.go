package models

import "time"

type Red struct {
	ID        uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	RedSocial string     `gorm:"column:red_social" json:"red_social"`
	Link      string     `gorm:"column:link" json:"link"`
	EmpresaID uint       `gorm:"column:empresa_id" json:"empresa_id"`
	Empresa   Empresa    `gorm:"foreignKey:EmpresaID"`
}

func (Red) TableName() string {
	return "redes"
}
