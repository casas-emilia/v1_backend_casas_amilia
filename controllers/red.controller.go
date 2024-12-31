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

// Función para crear datos de una Red Social de la Empresa
func CrearRed(c *gin.Context) {
	var request dto.CrearRedRequest
	idParam := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
	}

	// Validamos el Body
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Creamos la Red
	red := models.Red{
		RedSocial: request.RedSocial,
		Link:      request.Link,
		EmpresaID: uint(empresaID),
	}

	// Agregamos los datos de la RedSocial a la Base de Datos
	if err := configs.DB.Create(&red).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo guardar los datos de la RedSocial")
		return
	}

	// Mostrmos un mensaje de éxito
	c.JSON(http.StatusOK, gin.H{"message": "RedSocial creada con éxito"})

	// Response
	response := dto.RedResponse{
		ID:        red.ID,
		RedSocial: red.RedSocial,
		Link:      red.Link,
		EmpresaID: red.EmpresaID,
	}

	// Mostramos el Response
	c.JSON(http.StatusCreated, gin.H{"red": response})
}

// Función paea Obtener todas las redes sociales de la empresa
func ObtenerRedes(c *gin.Context) {
	var redes []models.Red
	var redesResponse []dto.RedResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Buscamos todas las redes sociales de la Empresa
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").Find(&redes).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Redes sociales no encontradas")
		return
	}

	// Cilco para las redes sociales
	for _, red := range redes {
		redesResponse = append(redesResponse, dto.RedResponse{
			ID:        red.ID,
			RedSocial: red.RedSocial,
			Link:      red.Link,
			EmpresaID: red.EmpresaID,
		})
	}

	// Reponder/enviar redes sociales
	c.JSON(http.StatusOK, gin.H{"redes_sociales": redesResponse})
}

// Función para Obtener datos de una red social de acuerdo a su ID
func ObtenerRed(c *gin.Context) {
	var red models.Red
	var redResponse dto.RedResponse

	idParamRed := c.Param("redID")
	redID, err := strconv.ParseUint(idParamRed, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Red inválido")
		return
	}

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Buacar Red Social de acuerdo a su ID
	if err := configs.DB.Where("deleted_at IS NULL").Where("empresa_id = ?", empresaID).First(&red, redID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Red social no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Red Social")
		return
	}

	//
	redResponse = dto.RedResponse{
		ID:        red.ID,
		RedSocial: red.RedSocial,
		Link:      red.Link,
		EmpresaID: red.EmpresaID,
	}

	// Responder/mostrar Red Social
	c.JSON(http.StatusOK, gin.H{"red_social": redResponse})
}

// Función para actualizar datos de una red social de acuerdo a su ID
func ActualizarRed(c *gin.Context) {
	var red models.Red
	var redResponse dto.RedResponse
	var request dto.ActualizarRedRequest

	idParamRed := c.Param("redID")
	redID, err := strconv.ParseUint(idParamRed, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Red inválido")
		return
	}

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Bind Json del request a la estructura DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de Datos")
		return
	}

	// Buacar Red social de acuerdo al ID enviado
	if err := configs.DB.Where("deleted_at IS NULL").Where("empresa_id = ?", empresaID).First(&red, redID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Red social no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Red social")
		return
	}

	// Actualizar datos de la Red social
	red.RedSocial = request.RedSocial
	red.Link = request.Link

	// Guardar los datos en la base de Datos
	if err := configs.DB.Save(&red).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo actualizar los datos de la Red Social")
		return
	}

	redResponse = dto.RedResponse{
		ID:        red.ID,
		RedSocial: red.RedSocial,
		Link:      red.Link,
		EmpresaID: red.EmpresaID,
	}

	// Mostrar/enviar response y mensaje de exito
	c.JSON(http.StatusOK, gin.H{
		"message":    "Datos actualizados exitosamente",
		"red_social": redResponse,
	})
}

// Función para eliminar una red social de acuerdo al ID enviado
func EliminarRed(c *gin.Context) {
	var red models.Red

	idParamRed := c.Param("redID")
	redID, err := strconv.ParseUint(idParamRed, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Red inválido")
		return
	}

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Buscar red slcial con el ID enviado
	if err := configs.DB.Unscoped().Where("empresa_id = ?", empresaID).First(&red, redID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Servicio no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos del Servicio")
		return
	}

	// Verificar si la red social ya esta eliminada lógicamente
	if red.DeletedAt != nil && !red.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "Los datos de la Red Social ya se encuentran eliminados")
		return
	}

	// Establecer fecha y hora para la eliminación lógica
	now := time.Now()
	red.DeletedAt = &now

	// Guardar fecha y hora de la eliminación lógica en la base de datos
	if err := configs.DB.Save(&red).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo eliminar la Red Social")
		return
	}

	// Mostrar/enviar un mensaje de éxito de la eliminación lógica
	c.JSON(http.StatusOK, gin.H{"message": "Red Social eliminada exitosamente"})
}
