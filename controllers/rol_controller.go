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

// Función para crear un Rol
func CrearRol(c *gin.Context) {
	var request dto.CrearRolRequest
	var rolResponse dto.RolResponse

	// Validamos el body
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error de datos "+err.Error())
		return
	}

	// Creamos el Rol
	rol := models.Rol{
		NombreRol:      request.NombreRol,
		DescripcionRol: request.DescripcionRol,
	}

	// Agregamos el Rol a la Base de Datos
	if err := configs.DB.Create(&rol).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo crear el Rol")
		return
	}

	// Response
	rolResponse = dto.RolResponse{
		ID:             rol.ID,
		CreatedAt:      rol.CreatedAt,
		UpdatedAt:      rol.UpdatedAt,
		NombreRol:      rol.NombreRol,
		DescripcionRol: rol.DescripcionRol,
	}

	// Mostramos mensaje de creación éxitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Rol creado con éxito",
	})

	// Mostramos el Response
	c.JSON(http.StatusCreated, gin.H{
		"rol": rolResponse,
	})
}

// Función para obtener todos los roles
func ObtenerRoles(c *gin.Context) {
	var roles []models.Rol
	var rolesResponse []dto.RolResponse

	// Buscar todos los roles en la base de datos
	if err := configs.DB.Where("deleted_at IS NULL").Find(&roles).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de obtener los Roles")
		return
	}

	if len(roles) == 0 {
		HandleError(c, nil, http.StatusNotFound, "Roles no encontrados")
		return
	}

	for _, rol := range roles {
		rolesResponse = append(rolesResponse, dto.RolResponse{
			ID:             rol.ID,
			CreatedAt:      rol.CreatedAt,
			UpdatedAt:      rol.UpdatedAt,
			NombreRol:      rol.NombreRol,
			DescripcionRol: rol.DescripcionRol,
		})
	}

	// Mostrar/enviar Roles
	c.JSON(http.StatusOK, gin.H{
		"roles": rolesResponse,
	})
}

// Función para obtener el Rol de acuerdo a su ID
func ObtenerRol(c *gin.Context) {
	var rol models.Rol
	var rolResponse dto.RolResponse

	idParamRol := c.Param("rolID")
	rolID, err := strconv.ParseUint(idParamRol, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol inválido")
		return
	}

	// Buscar Rol en la base de datos de acuerdo a su ID
	if err := configs.DB.Where("deleted_at IS NULL").Where("id = ?", rolID).First(&rol).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Rol no encontrsado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de buscar el Rol")
		return
	}

	rolResponse = dto.RolResponse{
		ID:             rol.ID,
		CreatedAt:      rol.CreatedAt,
		UpdatedAt:      rol.UpdatedAt,
		NombreRol:      rol.NombreRol,
		DescripcionRol: rol.DescripcionRol,
	}

	// Mostrar/enviar rol
	c.JSON(http.StatusOK, gin.H{
		"rol": rolResponse,
	})
}

// Función para actualizar datos de un Rol
func ActualizarRol(c *gin.Context) {
	var request dto.ActualizarRolRequest
	var rol models.Rol
	var rolResponse dto.RolResponse

	// Bind JSON del request a la estructura del DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error de Datos")
		return
	}

	idParamRol := c.Param("rolID")
	rolID, err := strconv.ParseUint(idParamRol, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol inválido")
		return
	}

	// Buscar rol en la base de datos de acuerdo a su ID
	if err := configs.DB.Where("deleted_at IS NULL").Where("id = ?", rolID).First(&rol).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Rol no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de obtener datos del Rol")
		return
	}

	// Actualizar datos del Rol
	rol.NombreRol = request.NombreRol
	rol.DescripcionRol = request.DescripcionRol

	// Guardar los cambios en la base de datos
	if err := configs.DB.Save(&rol).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se puedo actalizar los datos del Rol")
		return
	}

	rolResponse = dto.RolResponse{
		ID:             rol.ID,
		CreatedAt:      rol.CreatedAt,
		UpdatedAt:      rol.UpdatedAt,
		NombreRol:      rol.NombreRol,
		DescripcionRol: rol.DescripcionRol,
	}

	// Mostrar/enviar mensaje de actualización exitosa y el rolResponse
	c.JSON(http.StatusOK, gin.H{
		"message": "Rol actualizado con éxito",
		"rol":     rolResponse,
	})

}

// Función para elimar un rol lógicamente
func EliminarRol(c *gin.Context) {
	var rol models.Rol

	idParamRol := c.Param("rolID")
	rolID, err := strconv.ParseUint(idParamRol, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID Rol inválido")
		return
	}

	// Buscar el Rol en la base de datos
	if err := configs.DB.Where("deleted_at IS NULL").Where("id = ?", rolID).First(&rol).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Rol no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo obtener datos del Rol")
		return
	}

	// Verificar si el Rol ya se encuentra eliminado lógicamente
	/* if rol.DeletedAt != nil && !rol.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "Rol seleccionado ya se encuentra eliminado")
		return
	} */

	// Poner Fecha y hora de la eliminación lógica
	now := time.Now()
	rol.DeletedAt = &now

	// Guarda Fecha y hora de la eliminación lógica del Rol
	if err := configs.DB.Save(&rol).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar el Rol, intente nuevamente más tarde")
		return
	}

	// Mostrar/enviar mensaje de eliminación exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Rol eliminado con éxito",
	})
}
