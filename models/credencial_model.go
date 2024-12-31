package models

import "time"

type Credencial struct {
	ID        uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	Email     string     `gorm:"unique;not null;column:email" json:"email"`
	Password  string     `gorm:"not null" json:"password"`
	UsuarioID uint       `gorm:"unique;not null;column:usuario_id" json:"usuario_id"`
	Usuario   *Usuario   `gorm:"foreignKey:UsuarioID"`
}

func (Credencial) TableName() string {
	return "credenciales"
}
