package services

import (
	"fiber-project/database"
	"fiber-project/models"
	"math/rand"
	"net/smtp"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	// Parseia o corpo da requisição
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Gera um token de redefinição de senha
	token := RandStringRunes(12)

	// Cria um registro na tabela de resets
	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}

	database.DB.Create(&passwordReset)

	// Configuração de autenticação SMTP para Gmail
	auth := smtp.PlainAuth("", os.Getenv("GMAIL_USERNAME"), os.Getenv("GMAIL_PASSWORD"), "smtp.gmail.com")

	to := []string{data["email"]}
	msg := []byte("To: " + data["email"] + "\r\n" +
		"Subject: Redefina sua senha\r\n" +
		"\r\n" +
		"Use o link para redefinir sua senha: " + os.Getenv("APP_URL") + "/reset/" + token + "\r\n")

	// Envia o e-mail usando o servidor SMTP do Gmail
	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("GMAIL_USERNAME"), to, msg)
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

	// Parseia o corpo da requisição
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Verifica se as senhas coincidem
	if data["password"] != data["confirm_password"] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	// Verifica o token e busca o registro correspondente
	var passwordReset = models.PasswordReset{}
	if err := database.DB.Where("token = ?", data["token"]).Last(&passwordReset); err.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token!",
		})
	}

	// Gera o hash da nova senha
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// Atualiza a senha do usuário
	database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	// Retorna uma resposta de sucesso
	return c.JSON(fiber.Map{
		"message": "Password successfully reset",
	})
}

// Função auxiliar para gerar um token aleatório
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
