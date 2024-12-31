package main

import (
	"log"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/models"
)

func init() {
	configs.ConnectToDB()
}

func main() {
	log.Println("Iniciando migraciones...")

	err := configs.DB.AutoMigrate(
		&models.Caracteristica{},
		&models.Categoria{},
		&models.Contacto{},
		&models.Credencial{},
		&models.Empresa{},
		&models.Estilo{},
		&models.Imagen_noticia{},
		&models.Imagen_prefabricada{},
		&models.Incluye{},
		&models.Noticia{},
		&models.Portada{},
		&models.Precio{},
		&models.Prefabricada{},
		&models.Red{},
		&models.Rol{},
		&models.Rol_usuario{},
		&models.Servicio{},
		&models.Tipo{},
		&models.Tipo_categoria{},
		&models.Usuario{},
		&models.Recuperacion{},
	)
	if err != nil {
		log.Fatalf("Error durante la migración: %v", err)
	}

	log.Println("Migraciones completadas exitosamente")

	// Seed initial data
	seedInitialData()
}

func seedInitialData() {
	// Create super_administrador role
	superAdminRole := models.Rol{
		NombreRol:      "super_administrador",
		DescripcionRol: "Super administrador del sistema",
	}
	result := configs.DB.Create(&superAdminRole)
	if result.Error != nil {
		log.Fatalf("Error al crear el rol super_administrador: %v", result.Error)
	}

	// Crear una empresa
	empresa := models.Empresa{
		NombreEmpresa:      "Empresa Demo",
		DescripcionEmpresa: "Descripción de la empresa demo",
		HistoriaEmpresa:    "Historia de la empresa demo",
		MisionEmpresa:      "Misión de la empresa demo",
		VisionEmpresa:      "Visión de la empresa demo",
		UbicacionEmpresa:   "Ubicación demo",
		CelularEmpresa:     "123456789",
		EmailEmpresa:       "empresa@demo.com",
	}
	result = configs.DB.Create(&empresa)
	if result.Error != nil {
		log.Fatalf("Error al crear la empresa inicial: %v", result.Error)
	}

	// Create super_administrador user
	superAdminUser := models.Usuario{
		PrimerNombre:   "Admin",
		PrimerApellido: "Demo",
		EmpresaID:      empresa.ID,
		Credencial: &models.Credencial{
			Email:    "7.cristian.u@gmail.com",
			Password: "password", // Reemplazar con la contraseña hasheada
		},
	}
	result = configs.DB.Create(&superAdminUser)
	if result.Error != nil {
		log.Fatalf("Error al crear el usuario super_administrador: %v", result.Error)
	}

	// Assign role to user
	rolUsuario := models.Rol_usuario{
		UsuarioID: superAdminUser.ID,
		RolID:     superAdminRole.ID,
	}
	result = configs.DB.Create(&rolUsuario)
	if result.Error != nil {
		log.Fatalf("Error al asignar el rol super_administrador al usuario: %v", result.Error)
	}

	log.Println("Datos iniciales creados exitosamente")
}