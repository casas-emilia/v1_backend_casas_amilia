package dto

import "time"

type CrearPrecioRequest struct {
	NombrePrecio      string  `json:"nombre_precio" binding:"required"`
	DescripcionPrecio string  `json:"descripcion_precio" binding:"required"`
	ValorPrefabricada float64 `json:"valor_prefabricada" binding:"required"`
	//PrefabricadaID    uint    `json:"prefabricada_id" binding:"required"`
}

type ActualizarPrecioRequest struct {
	NombrePrecio      string  `json:"nombre_precio" binding:"required"`
	DescripcionPrecio string  `json:"descripcion_precio" binding:"required"`
	ValorPrefabricada float64 `json:"valor_prefabricada" binding:"required"`
}

type PrecioResponse struct {
	ID                uint              `json:"id"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	NombrePrecio      string            `json:"nombre_precio" binding:"required"`
	DescripcionPrecio string            `json:"descripcion_precio" binding:"required"`
	ValorPrefabricada float64           `json:"valor_prefabricada" binding:"required"`
	PrefabricadaID    uint              `json:"prefabricada_id" binding:"required"`
	Incluyes          []IncluyeResponse `json:"incluyes"`
}
