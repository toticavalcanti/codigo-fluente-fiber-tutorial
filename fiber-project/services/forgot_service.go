package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fiber-project/database"
	"fiber-project/models"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Template HTML para o email
const emailTemplate = `
<!DOCTYPE html>
<html>
<body>
    <h2>Redefinição de Senha</h2>
    <p>Você solicitou a redefinição de sua senha.</p>
    <p>Clique no link abaixo para redefinir sua senha:</p>
    <a href="{{.ResetLink}}">Redefinir Senha</a>
    <p>Este link expira em 1 hora.</p>
    <p>Se você não solicitou esta redefinição, ignore este email.</p>
</body>
</html>
`

func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Verificar se o email existe
	var user models.User
	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		// Não revelamos se o email existe ou não por segurança
		return c.JSON(fiber.Map{
			"message": "If the email exists, you will receive reset instructions",
		})
	}

	// Gerar token seguro usando crypto/rand
	token := generateSecureToken(32)

	// Criar registro de reset com expiração
	passwordReset := models.PasswordReset{
		Email:     data["email"],
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := database.DB.Create(&passwordReset).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error processing request",
		})
	}

	// Configurar autenticação Gmail
	auth := smtp.PlainAuth("",
		os.Getenv("GMAIL_EMAIL"),
		os.Getenv("GMAIL_APP_PASSWORD"),
		"smtp.gmail.com",
	)

	// Preparar o template do email
	resetLink := fmt.Sprintf("%s/reset/%s", os.Getenv("APP_URL"), token)
	emailBody, err := parseEmailTemplate(emailTemplate, resetLink)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error preparing email",
		})
	}

	// Montar a mensagem de email
	msg := buildEmailMessage(data["email"], "Redefinição de Senha", emailBody)

	// Enviar email
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("GMAIL_EMAIL"),
		[]string{data["email"]},
		msg,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error sending email",
		})
	}

	return c.JSON(fiber.Map{
		"message": "If the email exists, you will receive reset instructions",
	})
}

func Reset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if data["password"] != data["confirm_password"] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Passwords do not match",
		})
	}

	// Verificar token e expiração
	var passwordReset models.PasswordReset
	if err := database.DB.Where("token = ? AND expires_at > ?",
		data["token"],
		time.Now(),
	).Last(&passwordReset).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}

	// Gerar novo hash de senha
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error processing password",
		})
	}

	// Atualizar senha
	if err := database.DB.Model(&models.User{}).
		Where("email = ?", passwordReset.Email).
		Update("password", password).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating password",
		})
	}

	// Invalidar token usado
	database.DB.Delete(&passwordReset)

	return c.JSON(fiber.Map{
		"message": "Password successfully reset",
	})
}

// Funções auxiliares

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func parseEmailTemplate(tmpl string, resetLink string) (string, error) {
	t, err := template.New("reset_email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	data := struct {
		ResetLink string
	}{
		ResetLink: resetLink,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func buildEmailMessage(to, subject, body string) []byte {
	return []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s", to, subject, body))
}
