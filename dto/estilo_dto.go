package dto

import "time"

type CrearEstiloRequest struct {
	NombreEstilo      string `json:"nombre_estilo" binding:"required"`
	DescripcionEstilo string `json:"descripcion_estilo"`
}

type ActualizarEstiloRequest struct {
	NombreEstilo      string `json:"nombre_estilo" binding:"required"`
	DescripcionEstilo string `json:"descripcion_estilo"`
}

type EstiloResponse struct {
	ID                uint      `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	NombreEstilo      string    `json:"nombre_estilo"`
	DescripcionEstilo string    `json:"descripcion_estilo"`
}
