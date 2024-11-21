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

// Função para enviar o email de redefinição de senha
func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	// Parseia o corpo da requisição para map
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Gera um token aleatório para redefinição de senha
	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{
		Email:     data["email"],
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Salva o token no banco de dados
	database.DB.Create(&passwordReset)

	// Configuração SMTP para enviar o email
	auth := smtp.PlainAuth("", os.Getenv("GMAIL_EMAIL"), os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com")

	// Destinatário e conteúdo do email
	to := []string{data["email"]}
	msg := []byte("To: " + data["email"] + "\r\n" +
		"Subject: Redefina sua senha\r\n" +
		"\r\n" +
		"Use o link para redefinir sua senha: " + os.Getenv("FRONTEND_URL") + "/reset/" + token + "\r\n")

	// Envia o email usando o servidor SMTP do Gmail
	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("GMAIL_EMAIL"), to, msg)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error sending email",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// Função para redefinir a senha
func Reset(c *fiber.Ctx) error {
	var data map[string]string

	// Parseia o corpo da requisição para map
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Verifica se as senhas coincidem
	if data["password"] != data["confirm_password"] {
		return c.Status(400).JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	// Busca o token no banco de dados
	var passwordReset = models.PasswordReset{}
	if err := database.DB.Where("token = ?", data["token"]).Last(&passwordReset); err.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid token!",
		})
	}

	// Gera o hash da nova senha
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// Atualiza a senha do usuário com base no email associado ao token
	database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	// Retorna sucesso
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// Função auxiliar para gerar tokens aleatórios
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
