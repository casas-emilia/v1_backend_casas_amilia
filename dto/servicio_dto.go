package dto

type CrearServicioRequest struct {
	NombreServicio      string `json:"nombre_servicio"`
	DescripcionServicio string `json:"descripcion_servicio"`
	EmpresaID           uint   `json:"empresa_id"`
}

type ActualizarServicioRequest struct {
	NombreServicio      string `json:"nombre_servicio"`
	DescripcionServicio string `json:"descripcion_servicio"`
	EmpresaID           uint   `json:"empresa_id"`
}

type ServicioResponse struct {
	ID                  uint   `json:"id"`
	NombreServicio      string `json:"nombre_servicio"`
	DescripcionServicio string `json:"descripcion_servicio"`
	EmpresaID           uint   `json:"empresa_id"`
}
