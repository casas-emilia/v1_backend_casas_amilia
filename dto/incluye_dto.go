package dto

import "time"

type CrearIncluyeRequest struct {
	NombreIncluye string `json:"nombre_incluye" binding:"required"`
	//PrecioID      uint   `json:"precio_id" binding:"required"`
}

type ActualizarIncluyeRequest struct {
	NombreIncluye string `json:"nombre_incluye" binding:"required"`
}

type IncluyeResponse struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	NombreIncluye string    `json:"nombre_incluye"`
	PrecioID      uint      `json:"precio_id"`
}
