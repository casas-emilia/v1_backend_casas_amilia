package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// EnviarEmailRecuperacion envía un email de recuperación de contraseña
func EnviarEmailRecuperacion(email, link string) error {
	// Configuración del servidor SMTP
	smtpHost := "smtp.gmail.com" // Cambia esto por el host de tu proveedor SMTP (por ejemplo, smtp.gmail.com)
	smtpPort := "587"            // Puerto estándar para conexiones TLS

	// Credenciales del remitente desde variables de entorno
	from := os.Getenv("EMAIL_ADDRESS")      // Dirección de email (en variable de entorno)
	password := os.Getenv("EMAIL_PASSWORD") // Contraseña o token de aplicación (en variable de entorno)

	// Verificar que las variables de entorno estén configuradas
	if from == "" || password == "" {
		return fmt.Errorf("las variables de entorno EMAIL_ADDRESS y EMAIL_PASSWORD no están configuradas")
	}

	// Construir el mensaje del email
	subject := "Recuperación de contraseña"
	body := fmt.Sprintf("Hola,\n\nHaz clic en el siguiente enlace para recuperar tu contraseña:\n\n%s\n\nSi no solicitaste esto, ignora este mensaje.", link)
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, email, subject, body)

	// Dirección del servidor SMTP
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Enviar el email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(message))
	if err != nil {
		log.Printf("Error al enviar el email: %v", err)
		return err
	}

	log.Printf("Email enviado a %s con éxito.", email)
	return nil
}
