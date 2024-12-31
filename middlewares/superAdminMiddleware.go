package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el rol almacenado en el contexto
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Rol no encontrado en el contexto"})
			c.Abort()
			return
		}

		// Verificar si el rol incluye "super_administrador"
		rolesList, ok := roles.([]string) // asumiendo que los roles están almacenados como slice de strings
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Formato de roles inválido"})
			c.Abort()
			return
		}

		isSuperAdmin := false
		for _, role := range rolesList {
			if role == "super_administrador" {
				isSuperAdmin = true
				break
			}
		}

		// Si no tiene el rol "super_administrador", bloquea el acceso
		if !isSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: Se requiere rol de super_administrador"})
			c.Abort()
			return
		}

		// Continuar si tiene el rol adecuado
		c.Next()
	}
}
