package services

import (
	"fiber-project/database"
	"fiber-project/models"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Geração de um token aleatório
	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{
		Email:     data["email"],
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Expira em 1 hora
	}

	// Salva no banco de dados
	if err := database.DB.Create(&passwordReset).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating password reset entry",
		})
	}

	// Configuração do e-mail
	auth := smtp.PlainAuth("", os.Getenv("GMAIL_USERNAME"), os.Getenv("GMAIL_PASSWORD"), "smtp.gmail.com")
	to := []string{data["email"]}
	msg := []byte("To: " + data["email"] + "\r\n" +
		"Subject: Password Reset Request\r\n" +
		"\r\n" +
		"Click the link below to reset your password:\r\n" +
		os.Getenv("APP_URL") + "/reset/" + token + "\r\n")

	// Envia o e-mail
	if err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("GMAIL_USERNAME"), to, msg); err != nil {
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

	// Verifica se as senhas coincidem
	if data["password"] != data["confirm_password"] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Passwords do not match",
		})
	}

	// Verifica a validade do token
	var passwordReset models.PasswordReset
	if err := database.DB.Where("token = ? AND expires_at > ?", data["token"], time.Now()).Last(&passwordReset).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}

	// Gera o hash da nova senha
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error processing password",
		})
	}

	// Atualiza a senha do usuário
	if err := database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating password",
		})
	}

	// Remove o token usado
	database.DB.Delete(&passwordReset)

	return c.JSON(fiber.Map{
		"message": "Password successfully reset",
	})
}

// Função auxiliar para gerar tokens aleatórios
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
