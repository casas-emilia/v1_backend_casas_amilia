package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función auxiliar para manejar errores
func handleErrorEmpresa(c *gin.Context, err error, statusCode int, message string) {
	if err != nil {
		c.JSON(statusCode, gin.H{"error": message})
	}
}

func CrearEmpresa(c *gin.Context) {
	var request dto.CrearEmpresaRequest

	// Validamos el body
	if err := c.ShouldBindJSON(&request); err != nil {
		handleErrorEmpresa(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Creación de Empresa
	empresa := models.Empresa{
		NombreEmpresa:      request.NombreEmpresa,
		DescripcionEmpresa: request.DescripcionEmpresa,
		HistoriaEmpresa:    request.HistoriaEmpresa,
		MisionEmpresa:      request.MisionEmpresa,
		VisionEmpresa:      request.VisionEmpresa,
		UbicacionEmpresa:   request.UbicacionEmpresa,
		CelularEmpresa:     request.CelularEmpresa,
		EmailEmpresa:       request.EmailEmpresa,
	}

	// Guardamos la Empresa en la Base de Datos
	if err := configs.DB.Create(&empresa).Error; err != nil {
		handleErrorEmpresa(c, err, http.StatusInternalServerError, "No se pudo crear Empresa")
		return
	}

	// Retornamos mensaje con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Empresa creada con éxito"})

	// Response
	response := dto.EmpresaResponse{
		ID:                 empresa.ID,
		NombreEmpresa:      empresa.NombreEmpresa,
		DescripcionEmpresa: empresa.DescripcionEmpresa,
		HistoriaEmpresa:    empresa.HistoriaEmpresa,
		MisionEmpresa:      empresa.MisionEmpresa,
		VisionEmpresa:      empresa.VisionEmpresa,
		UbicacionEmpresa:   empresa.UbicacionEmpresa,
		CelularEmpresa:     empresa.CelularEmpresa,
		EmailEmpresa:       empresa.EmailEmpresa,
	}
	// Responder con empresa creada
	c.JSON(http.StatusCreated, gin.H{"empresa": response})
}

// Función para Obtener Todas las Empresas
func ObtenerEmpresas(c *gin.Context) {
	var empresas []models.Empresa
	var empresasResponse []dto.EmpresaResponse

	// Buscar todas las Empresas con sus Servicios y Redes en la base de datos
	if err := configs.DB.
		Where("deleted_at IS NULL").
		Preload("Servicio", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar Serviios eliminadas lógicamente
		}).
		Preload("Red", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar RedesSociales eliminadas lógicamente
		}).
		Find(&empresas).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudieron obtener las Empresas")
		return
	}

	// Response
	for _, empresa := range empresas {
		var serviciosResponse []dto.ServicioResponse
		var redesResponse []dto.RedResponse

		for _, servicio := range empresa.Servicio {
			serviciosResponse = append(serviciosResponse, dto.ServicioResponse{
				ID:                  servicio.ID,
				NombreServicio:      servicio.NombreServicio,
				DescripcionServicio: servicio.DescripcionServicio,
				EmpresaID:           servicio.EmpresaID,
			})
		}

		for _, red := range empresa.Red {
			redesResponse = append(redesResponse, dto.RedResponse{
				RedSocial: red.RedSocial,
				Link:      red.Link,
				EmpresaID: red.EmpresaID,
			})
		}

		empresasResponse = append(empresasResponse, dto.EmpresaResponse{
			ID:                 empresa.ID,
			UpdatedAt:          empresa.UpdatedAt,
			NombreEmpresa:      empresa.NombreEmpresa,
			DescripcionEmpresa: empresa.DescripcionEmpresa,
			HistoriaEmpresa:    empresa.HistoriaEmpresa,
			MisionEmpresa:      empresa.MisionEmpresa,
			VisionEmpresa:      empresa.VisionEmpresa,
			UbicacionEmpresa:   empresa.UbicacionEmpresa,
			CelularEmpresa:     empresa.CelularEmpresa,
			EmailEmpresa:       empresa.EmailEmpresa,
			Servicios:          serviciosResponse,
			Redes:              redesResponse,
		})
	}

	// Responder/mostrar Empresas
	c.JSON(http.StatusOK, gin.H{"empresas": empresasResponse})
}

// Función para obtener una Empresa de acuerdo a su ID
func ObtenerEmpresa(c *gin.Context) {
	var empresa models.Empresa
	var empresaResponse dto.EmpresaResponse
	idParam := c.Param("empresaID")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	// Buscar los Datos de la Empresa de acuerdo a su ID
	if err := configs.DB.
		Where("deleted_at IS NULL").
		Preload("Servicio", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar Serviios eliminadas lógicamente
		}).
		Preload("Red", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar RedesSociales eliminadas lógicamente
		}).
		First(&empresa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Empresa no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Empresa")
		return
	}

	// Construir la respuesta si se encuentra la Empresa
	var serviciosResponse []dto.ServicioResponse
	var redesResponse []dto.RedResponse

	for _, servicio := range empresa.Servicio {
		serviciosResponse = append(serviciosResponse, dto.ServicioResponse{
			ID:                  servicio.ID,
			NombreServicio:      servicio.NombreServicio,
			DescripcionServicio: servicio.DescripcionServicio,
			EmpresaID:           servicio.EmpresaID,
		})
	}

	for _, red := range empresa.Red {
		redesResponse = append(redesResponse, dto.RedResponse{
			ID:        red.ID,
			RedSocial: red.RedSocial,
			Link:      red.Link,
			EmpresaID: red.EmpresaID,
		})
	}
	empresaResponse = dto.EmpresaResponse{
		ID:                 empresa.ID,
		UpdatedAt:          empresa.UpdatedAt,
		NombreEmpresa:      empresa.NombreEmpresa,
		DescripcionEmpresa: empresa.DescripcionEmpresa,
		HistoriaEmpresa:    empresa.HistoriaEmpresa,
		MisionEmpresa:      empresa.MisionEmpresa,
		VisionEmpresa:      empresa.VisionEmpresa,
		UbicacionEmpresa:   empresa.UbicacionEmpresa,
		CelularEmpresa:     empresa.CelularEmpresa,
		EmailEmpresa:       empresa.EmailEmpresa,
		Servicios:          serviciosResponse,
		Redes:              redesResponse,
	}

	// Mostrar/enviar respuesta
	c.JSON(http.StatusOK, gin.H{"empresa": empresaResponse})
}

// Función para actualizar datos de Empresa
func ActualizarEmpresa(c *gin.Context) {
	var request dto.ActualizarEmpresaRequest
	var empresa models.Empresa
	var empresaResponse dto.EmpresaResponse
	idParam := c.Param("empresaID")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	// Bind Json del request a la estructura del dto
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Buscar Empresa de acuerdo a su ID
	if err := configs.DB.
		Where("deleted_at IS NULL").
		Preload("Servicio", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar Serviios eliminadas lógicamente
		}).
		Preload("Red", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar RedesSociales eliminadas lógicamente
		}).
		First(&empresa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Empresa no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Empresa")
		return
	}

	// Actulizar datos de la Empresa
	empresa.NombreEmpresa = request.NombreEmpresa
	empresa.DescripcionEmpresa = request.DescripcionEmpresa
	empresa.HistoriaEmpresa = request.HistoriaEmpresa
	empresa.MisionEmpresa = request.MisionEmpresa
	empresa.VisionEmpresa = request.VisionEmpresa
	empresa.UbicacionEmpresa = request.UbicacionEmpresa
	empresa.CelularEmpresa = request.CelularEmpresa
	empresa.EmailEmpresa = request.EmailEmpresa

	// Guardar los datos en la base de datos
	if err := configs.DB.Save(&empresa).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo actualizar los datos de la Empresa")
		return
	}

	// Response
	var serviciosResponse []dto.ServicioResponse
	var redesResponse []dto.RedResponse

	for _, servicio := range empresa.Servicio {
		serviciosResponse = append(serviciosResponse, dto.ServicioResponse{
			NombreServicio:      servicio.NombreServicio,
			DescripcionServicio: servicio.DescripcionServicio,
			EmpresaID:           servicio.EmpresaID,
		})
	}

	for _, redes := range empresa.Red {
		redesResponse = append(redesResponse, dto.RedResponse{
			RedSocial: redes.RedSocial,
			Link:      redes.Link,
			EmpresaID: redes.EmpresaID,
		})
	}

	empresaResponse = dto.EmpresaResponse{
		ID:                 empresa.ID,
		UpdatedAt:          empresa.UpdatedAt,
		NombreEmpresa:      empresa.NombreEmpresa,
		DescripcionEmpresa: empresa.DescripcionEmpresa,
		HistoriaEmpresa:    empresa.HistoriaEmpresa,
		MisionEmpresa:      empresa.MisionEmpresa,
		VisionEmpresa:      empresa.VisionEmpresa,
		UbicacionEmpresa:   empresa.UbicacionEmpresa,
		CelularEmpresa:     empresa.CelularEmpresa,
		EmailEmpresa:       empresa.EmailEmpresa,
		Servicios:          serviciosResponse,
		Redes:              redesResponse,
	}

	// Reponder/mostrar con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Datos de Empresa actulizados exitosamente"})
	// enviar/mostrar response
	c.JSON(http.StatusOK, gin.H{"Datos de Empresa": empresaResponse})
}

// Función para eliminar lógicamente una Empresa
func EliminarEmpresa(c *gin.Context) {
	var empresa models.Empresa
	idParam := c.Param("empresaID")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	// Buscar Empresa por su ID
	if err := configs.DB.Unscoped().First(&empresa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Empresa no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Empresa")
		return
	}

	// Verificar si la Empresa ya esta Eliminada lógicamente
	if empresa.DeletedAt != nil && !empresa.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "La empresa ya está eliminada")
		return
	}

	// Establecer la fecha y hora de la eliminación lógica
	now := time.Now()
	empresa.DeletedAt = &now

	// Guardar en la base de datos la fecha y hora de la eliminación lógica
	if err := configs.DB.Save(&empresa).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo eliminar la Empresa")
		return
	}

	// Responder/enviar con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Empresa eliminada exitosamente"})
}
