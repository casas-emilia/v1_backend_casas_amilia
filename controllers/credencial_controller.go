package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"
	"v1_prefabricadas/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función para crear Credenciales de un Usuario
func CrearCredencial(c *gin.Context) {
	var request dto.CrearCredencial
	var credencial models.Credencial
	var credencialResponse dto.CredencialResponse

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos "+err.Error())
		return
	}

	// Verificar si el usuario ya tiene credenciales
	if err := configs.DB.Where("usuario_id = ?", usuarioID).First(&credencial).Error; err == nil {
		// Credenciales ya existen
		HandleError(c, nil, http.StatusConflict, "El Usuario ya cuenta con credenciales de acceso")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Error inesperado en la consulta
		HandleError(c, err, http.StatusInternalServerError, "Error al verificar credenciales existentes")
		return
	}

	// Iniciar una transacción para crear las credenciales
	tx := configs.DB.Begin()

	// Hashear la contraseña
	hashedpassword, err := services.HashPassword(request.Password)
	if err != nil {
		tx.Rollback()
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo hashear la contraseña")
		return
	}

	credencial.Email = request.Email
	credencial.Password = hashedpassword
	credencial.UsuarioID = uint(usuarioID)

	if err := tx.Create(&credencial).Error; err != nil {
		tx.Rollback()
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo crear las credenciales, intente con otro Email")
		return
	}

	tx.Commit()

	credencialResponse = dto.CredencialResponse{
		ID:        credencial.ID,
		CreatedAt: credencial.CreatedAt,
		UpdatedAt: credencial.UpdatedAt,
		Email:     credencial.Email,
		UsuarioID: credencial.UsuarioID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Credenciales de acceso creadas con éxito",
		"credencial": credencialResponse,
	})
}

// Función para obtener email de credenciales
func ObtenerCredenciales(c *gin.Context) {
	var credencialesResponse dto.CredencialResponse
	var credenciales models.Credencial

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	// Buscar credencial
	if err := configs.DB.Where("usuario_id = ?", usuarioID).First(&credenciales).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Credenciales no encontradas")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener credenciales")
		return
	}

	credencialesResponse = dto.CredencialResponse{
		ID:        credenciales.ID,
		CreatedAt: credenciales.CreatedAt,
		UpdatedAt: credenciales.UpdatedAt,
		Email:     credenciales.Email,
		UsuarioID: credenciales.UsuarioID,
	}

	// Enviar mensaje de éxito
	c.JSON(http.StatusOK, gin.H{
		"credenciales": credencialesResponse,
	})

}

// Función para actualizar Credenciales
func ActualizarCredenciales(c *gin.Context) {
	var request dto.ActualizarCredencialRequest
	var credencial models.Credencial
	var credencialResponse dto.CredencialResponse

	idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	}

	idParamCredencial := c.Param("credencialID")
	credencialID, err := strconv.ParseUint(idParamCredencial, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Credencial inválido")
		return
	}

	// Bind JSON del request a la estructura DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos "+err.Error())
		return
	}

	// Iniciar transacción
	tx := configs.DB.Begin()

	// Buscar credencial
	if err := tx.Where("usuario_id = ? AND id = ?", usuarioID, credencialID).First(&credencial).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Credenciales no encontradas")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener credenciales")
		return
	}

	// Validar y actualizar el email
	if request.Email != "" {
		credencial.Email = request.Email
	} else {
		tx.Rollback()
		HandleError(c, nil, http.StatusBadRequest, "El campo Email no debe estar vacío")
		return
	}

	// Validar y actualizar la contraseña solo si se proporciona
	if request.Password != "" {
		hashedPassword, err := services.HashPassword(request.Password)
		if err != nil {
			tx.Rollback()
			HandleError(c, err, http.StatusInternalServerError, "Error al hashear la nueva contraseña")
			return
		}
		credencial.Password = hashedPassword
	}

	// Guardar actualización en la base de datos
	if err := tx.Save(&credencial).Error; err != nil {
		tx.Rollback()
		HandleError(c, err, http.StatusInternalServerError, "Error al actualizar credenciales")
		return
	}

	// Confirmar la transacción
	tx.Commit()

	// Preparar la respuesta
	credencialResponse = dto.CredencialResponse{
		ID:        credencial.ID,
		CreatedAt: credencial.CreatedAt,
		UpdatedAt: credencial.UpdatedAt,
		Email:     credencial.Email,
		UsuarioID: credencial.UsuarioID,
	}

	// Enviar mensaje de éxito
	c.JSON(http.StatusOK, gin.H{
		"message":      "Datos de acceso actualizados exitosamente",
		"credenciales": credencialResponse,
	})
}
