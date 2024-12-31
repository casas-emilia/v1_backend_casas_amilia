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

// Función para crear un Servicio
func CrearServicio(c *gin.Context) {
	var request dto.CrearServicioRequest
	idParam := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	// Validamos el Body
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Creamos el Servicio
	servicio := models.Servicio{
		NombreServicio:      request.NombreServicio,
		DescripcionServicio: request.DescripcionServicio,
		EmpresaID:           uint(empresaID),
	}

	// Agregamos el Servicio a la Base de Datos
	if err := configs.DB.Create(&servicio).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo crear el Servicio")
		return
	}

	// Mostramos un mensaje con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Servicio creada con éxito"})

	// Response
	response := dto.ServicioResponse{
		ID:                  servicio.ID,
		NombreServicio:      servicio.NombreServicio,
		DescripcionServicio: servicio.DescripcionServicio,
		EmpresaID:           servicio.EmpresaID,
	}

	// Mostramos el Response
	c.JSON(http.StatusCreated, gin.H{"servicio": response})
}

// Función para obtener todos los servicio de la Empresa
func ObtenerServicios(c *gin.Context) {
	var servicios []models.Servicio
	var serviciosResponse []dto.ServicioResponse
	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Buscar todos los servicios de la Empresa de acuerdoa su ID
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").Find(&servicios).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Servicios no encontrados")
		return
	}

	for _, servicio := range servicios {
		serviciosResponse = append(serviciosResponse, dto.ServicioResponse{
			ID:                  servicio.ID,
			NombreServicio:      servicio.NombreServicio,
			DescripcionServicio: servicio.DescripcionServicio,
			EmpresaID:           servicio.EmpresaID,
		})
	}

	// Mostrar/enviar respuesta
	c.JSON(http.StatusOK, gin.H{"servicios": serviciosResponse})
}

// Obtener un Servicio por su ID
func ObtenerServicio(c *gin.Context) {
	var servicio models.Servicio
	var ServicioResponse dto.ServicioResponse

	idParamServicio := c.Param("servicioID")
	servicioID, err := strconv.ParseUint(idParamServicio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Buscar el Servicio en la base de datos de acuerdo a su ID
	if err := configs.DB.Where("deleted_at IS NULL").Where("empresa_id = ?", empresaID).First(&servicio, servicioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Servicio no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos del Servicio")
		return
	}

	ServicioResponse = dto.ServicioResponse{
		ID:                  servicio.ID,
		NombreServicio:      servicio.NombreServicio,
		DescripcionServicio: servicio.DescripcionServicio,
		EmpresaID:           servicio.EmpresaID,
	}

	// Mostrar/enviar respuesta
	c.JSON(http.StatusOK, gin.H{"servicio": ServicioResponse})
}

// Actualizar un Servicio
func ActualizarServicio(c *gin.Context) {
	var servicio models.Servicio
	var servicioResponse dto.ServicioResponse
	var request dto.ActualizarServicioRequest

	idParamServicio := c.Param("servicioID")
	servicioID, err := strconv.ParseUint(idParamServicio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Bind json del request a la estructura dto
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de Datos")
		return
	}

	// Buscar datos del servicio de acuerdo a su ID
	if err := configs.DB.Where("deleted_at IS NULL").Where("empresa_id = ?", empresaID).First(&servicio, servicioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Servicio no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos del Servicio")
		return
	}

	// Actualizar datos del Servicio
	servicio.NombreServicio = request.NombreServicio
	servicio.DescripcionServicio = request.DescripcionServicio

	// Guardar los datos en la Base de Datos
	if err := configs.DB.Save(&servicio).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo actualizar los datos del Servicio")
		return
	}

	servicioResponse = dto.ServicioResponse{
		ID:                  servicio.ID,
		NombreServicio:      servicio.NombreServicio,
		DescripcionServicio: servicio.DescripcionServicio,
		EmpresaID:           servicio.EmpresaID,
	}

	// mostrar/enviar mensaje de éxito
	c.JSON(http.StatusOK, gin.H{
		"messge":   "Datos del Servicio Actualizado satisfactoriamente",
		"servicio": servicioResponse,
	})
}

// Eliminar lógicamente un Servicio
func EliminarServicio(c *gin.Context) {
	var servicio models.Servicio

	idParamServicio := c.Param("servicioID")
	servicioID, err := strconv.ParseUint(idParamServicio, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}
	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Buscar el Servicio de acuerdo a su ID
	if err := configs.DB.Unscoped().Where("empresa_id = ?", empresaID).First(&servicio, servicioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Servicio no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos del Servicio")
		return
	}

	// Verificar si el servicio ya se encuentra eliminado lógicamente
	if servicio.DeletedAt != nil && !servicio.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "El Servicio ya se encuentra eliminado")
		return
	}

	// Establecer facha y hora de la eliminación lógica
	now := time.Now()
	servicio.DeletedAt = &now

	// Guardar fecha y hora de ekliminación en la base de datos
	if err := configs.DB.Save(&servicio).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo eliminar Servicio")
		return
	}

	// Mostrar un mensaje de eliminación lógica exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Servicio eliminado exitosamente"})
}
