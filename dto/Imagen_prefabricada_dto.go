package dto

import "time"

type CrearImagen_prefabricadaRequest struct {
	//Image string `json:"image" binding:"required"`
	//PrefabricadaID uint   `json:"prefabricada_id" binding:"required"`
}

type ActualizarImagen_prefabricadaRequest struct {
	//Image string `json:"image" binding:"required"`
}

type Imagen_prefabricadaResponse struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_At"`
	Image          string    `json:"image"`
	PrefabricadaID uint      `json:"prefabricada_id"`
}
