package models

import "time"

type Recuperacion struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"unique;not null"`
	UsuarioID uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
