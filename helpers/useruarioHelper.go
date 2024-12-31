package helpers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Función para validar el ID de usuario en el contexto, permitiendo acceso a superadministradores
func ValidarUsuarioID(c *gin.Context) (uint, error) {

	// Obtener usuarioID y roles del contexto
	userID, exists := c.Get("usuarioID")
	if !exists {
		return 0, fmt.Errorf("no se pudo obtener el ID de usuario")
	}

	roles, exists := c.Get("roles")
	if !exists {
		return 0, fmt.Errorf("no se pudo obtener los roles del usuario")
	}

	// Convertir usuarioID a uint
	usuarioID, ok := userID.(uint)
	if !ok {
		return 0, fmt.Errorf("ID de usuario inválido")
	}

	// Verificar si el usuario tiene el rol de superadministrador
	rolesSlice, ok := roles.([]string)
	if !ok {
		return 0, fmt.Errorf("roles del usuario inválidos")
	}
	for _, rol := range rolesSlice {
		if rol == "super_administrador" {
			return usuarioID, nil // Permitir acceso sin validar ID si es superadministrador
		}
	}

	// Validar que el usuario solo acceda a sus propios recursos si no es superadmin
	paramUsuarioID, err := strconv.ParseUint(c.Param("usuarioID"), 10, 64)
	if err != nil || uint(paramUsuarioID) != usuarioID {
		return 0, fmt.Errorf("usuario no autorizado")
	}

	return usuarioID, nil
}
