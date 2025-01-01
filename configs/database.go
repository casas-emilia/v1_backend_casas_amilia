package configs

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	// Obtener las variables de entorno para la base de datos
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Verificar si todas las variables de entorno están configuradas
	requiredEnv := map[string]string{
		"DB_USER":     dbUser,
		"DB_PASSWORD": dbPassword,
		"DB_HOST":     dbHost,
		"DB_PORT":     dbPort,
		"DB_NAME":     dbName,
	}

	for key, value := range requiredEnv {
		if value == "" {
			log.Fatalf("Falta configurar la variable de entorno: %s", key)
		}
	}

	// Crear el DSN (Data Source Name) utilizando las variables de entorno
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Conectar a la base de datos
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Validar la conexión
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error al obtener la conexión de la base de datos: %v", err)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos exitosa")
}
