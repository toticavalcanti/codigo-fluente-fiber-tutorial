package services

import (
	"fiber-project/database"
	"fiber-project/models"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("Erro ao fazer parse do body em Forgot: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}

	database.DB.Create(&passwordReset)
	fmt.Printf("Token gerado para redefinição: %s\n", token)
	fmt.Printf("Email do usuário: %s\n", data["email"])

	// Configuração de autenticação SMTP para Gmail
	auth := smtp.PlainAuth("", os.Getenv("GMAIL_EMAIL"), os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com")

	to := []string{data["email"]}
	msg := []byte("To: " + data["email"] + "\r\n" +
		"Subject: Redefina sua senha\r\n" +
		"\r\n" +
		"Use o link para redefinir sua senha: " + os.Getenv("APP_URL") + "/reset/" + token + "\r\n")

	fmt.Printf("URL do reset no email: %s/reset/%s\n", os.Getenv("APP_URL"), token)

	// Envia o email usando o servidor SMTP do Gmail
	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("GMAIL_EMAIL"), to, msg)
	if err != nil {
		fmt.Printf("Erro ao enviar email: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error sending email",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Reset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("Erro ao fazer parse do body em Reset: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	fmt.Printf("Dados recebidos no Reset: %+v\n", data)

	// Valida se o token foi enviado
	if data["token"] == "" {
		fmt.Println("Token não fornecido na requisição")
		return c.Status(400).JSON(fiber.Map{
			"message": "Token is required!",
		})
	}

	fmt.Printf("Token recebido no backend: %s\n", data["token"])

	// Primeiro vamos listar todos os tokens no banco para debug
	var allTokens []models.PasswordReset
	if err := database.DB.Find(&allTokens).Error; err != nil {
		fmt.Printf("Erro ao buscar tokens: %v\n", err)
	}
	fmt.Printf("Tokens no banco: %+v\n", allTokens)

	// Busca o token específico
	var passwordReset = models.PasswordReset{}
	result := database.DB.Debug().Where("token = ?", data["token"]).Last(&passwordReset)

	if result.Error != nil {
		fmt.Printf("Erro ao buscar token específico: %v\n", result.Error)
		fmt.Printf("Token não encontrado: %s\n", data["token"])
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid token!",
			"error":   result.Error.Error(),
		})
	}

	// Se chegou aqui, encontrou o token, continua com o resto...
	if data["password"] != data["confirm_password"] {
		return c.Status(400).JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	updateResult := database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)
	if updateResult.Error != nil {
		fmt.Printf("Erro ao atualizar senha: %v\n", updateResult.Error)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error updating password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password successfully updated!",
	})
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
