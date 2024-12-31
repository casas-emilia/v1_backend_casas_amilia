package dto

import "time"

type CrearPrefabricadaRequest struct {
	NombrePrefabricada string `json:"nombre_prefabricada" binding:"required"`
	M2                 int    `json:"m2" binding:"required"`
	Garantia           string `json:"garantia" binding:"required"`
	Eslogan            string `json:"eslogan"`
	Descripcion        string `json:"descripcion"`
	Destacada          bool   `json:"destacada"`
	Oferta             bool   `json:"oferta"`
	CategoriaID        uint   `json:"categoria_id" binding:"required"`
	//EmpresaID          uint   `json:"empresa_id" binding:"required"`
	EstiloID uint `json:"estilo_id" binding:"required"`
	TipoID   uint `json:"tipo_id" binding:"required"`
}

type ActualizarPrefabricadaRequest struct {
	NombrePrefabricada string `json:"nombre_prefabricada" binding:"required"`
	M2                 int    `json:"m2" binding:"required"`
	Garantia           string `json:"garantia" binding:"required"`
	Eslogan            string `json:"eslogan"`
	Descripcion        string `json:"descripcion"`
	Destacada          bool   `json:"destacada"`
	Oferta             bool   `json:"oferta"`
	CategoriaID        uint   `json:"categoria_id" binding:"required"`
	//EmpresaID          uint   `json:"empresa_id" binding:"required"`
	EstiloID uint `json:"estilo_id" binding:"required"`
	TipoID   uint `json:"tipo_id"`
}
type PrefabricadaResponse struct {
	ID                    uint                          `json:"id"`
	CreatedAt             time.Time                     `json:"created_at"`
	UpdatedAt             time.Time                     `json:"updated_at"`
	NombrePrefabricada    string                        `json:"nombre_prefabricada" binding:"required"`
	M2                    int                           `json:"m2" binding:"required"`
	Garantia              string                        `json:"garantia" binding:"required"`
	Eslogan               string                        `json:"eslogan"`
	Descripcion           string                        `json:"descripcion"`
	Destacada             bool                          `json:"destacada"`
	Oferta                bool                          `json:"oferta"`
	CategoriaID           uint                          `json:"categoria_id" binding:"required"`
	EmpresaID             uint                          `json:"empresa_id" binding:"required"`
	EstiloID              uint                          `json:"estilo_id" binding:"required"`
	TipoID                uint                          `json:"tipo_id"`
	ImagenesPrefabricadas []Imagen_prefabricadaResponse `json:"imagenes_prefabricadas"`
	Caracteristicas       []CaracteristicaResponse      `json:"caracteristicas"`
	Precios               []PrecioResponse              `json:"precios"`
}
