package dto

import "time"

type CrearTipoRequest struct {
	MaterialEstructura  string `json:"material_estructura" binding:"required"`
	DescripcionMaterial string `json:"descripcion_material" binding:"required"`
}

type ActualizarTipoRequest struct {
	MaterialEstructura  string `json:"material_estructura" binding:"required"`
	DescripcionMaterial string `json:"descripcion_material" binding:"required"`
}

type TipoResponse struct {
	ID                  uint      `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	MaterialEstructura  string    `json:"material_estructura"`
	DescripcionMaterial string    `json:"descripcion_material"`
}
