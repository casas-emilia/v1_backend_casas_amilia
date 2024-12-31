package dto

import "time"

type CrearEmpresaRequest struct {
	NombreEmpresa      string `json:"nombre_empresa" binding:"required"`
	DescripcionEmpresa string `json:"descripcion_empresa"`
	HistoriaEmpresa    string `json:"historia_empresa"`
	MisionEmpresa      string `json:"mision_empresa"`
	VisionEmpresa      string `json:"vision_empresa"`
	UbicacionEmpresa   string `json:"ubicacion_empresa"`
	CelularEmpresa     string `json:"celular_empresa" binding:"required"`
	EmailEmpresa       string `json:"email_empresa" binding:"required,email"`
}

type ActualizarEmpresaRequest struct {
	NombreEmpresa      string `json:"nombre_empresa" binding:"required"`
	DescripcionEmpresa string `json:"descripcion_empresa"`
	HistoriaEmpresa    string `json:"historia_empresa"`
	MisionEmpresa      string `json:"mision_empresa"`
	VisionEmpresa      string `json:"vision_empresa"`
	UbicacionEmpresa   string `json:"ubicacion_empresa"`
	CelularEmpresa     string `json:"celular_empresa" binding:"required"`
	EmailEmpresa       string `json:"email_empresa" binding:"required,email"`
}

type EmpresaResponse struct {
	ID                 uint               `json:"id"`
	UpdatedAt          time.Time          `json:"updated_at"`
	NombreEmpresa      string             `json:"nombre_empresa"`
	DescripcionEmpresa string             `json:"descripcion_empresa"`
	HistoriaEmpresa    string             `json:"historia_empresa"`
	MisionEmpresa      string             `json:"mision_empresa"`
	VisionEmpresa      string             `json:"vision_empresa"`
	UbicacionEmpresa   string             `json:"ubicacion_empresa"`
	CelularEmpresa     string             `json:"celular_empresa"`
	EmailEmpresa       string             `json:"email_empresa"`
	Servicios          []ServicioResponse `json:"servicios"`
	Redes              []RedResponse      `json:"redes"`
}
