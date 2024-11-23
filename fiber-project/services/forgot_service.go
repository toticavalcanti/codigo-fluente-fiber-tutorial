package services

import (
	"fiber-project/database"
	"fiber-project/models"
	"fmt"
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
		fmt.Printf("Erro ao fazer parse do body em Forgot: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{
		Email:     data["email"],
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	fmt.Printf("Token criado com expiração para: %v\n", passwordReset.ExpiresAt)

	if err := database.DB.Create(&passwordReset).Error; err != nil {
		fmt.Printf("Erro ao salvar token: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error saving token",
		})
	}

	// Configuração de autenticação SMTP para Gmail
	auth := smtp.PlainAuth("", os.Getenv("GMAIL_EMAIL"), os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com")

	to := []string{data["email"]}
	msg := []byte("To: " + data["email"] + "\r\n" +
		"Subject: Redefina sua senha\r\n" +
		"\r\n" +
		"Use o link para redefinir sua senha: " + os.Getenv("APP_URL") + "/reset/" + token + "\r\n")

	// Envia o email usando o servidor SMTP do Gmail
	if err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("GMAIL_EMAIL"), to, msg); err != nil {
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

	// Valida se o token foi enviado no corpo da requisição
	if data["token"] == "" {
		fmt.Println("Token não fornecido na requisição")
		return c.Status(400).JSON(fiber.Map{
			"message": "Token is required!",
		})
	}

	fmt.Printf("Token recebido no backend: %s\n", data["token"])

	// Verifica se a senha e confirmação coincidem
	if data["password"] != data["confirm_password"] {
		fmt.Println("Senhas não coincidem")
		return c.Status(400).JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	// Busca o registro do token no banco de dados
	var passwordReset models.PasswordReset
	result := database.DB.Where("token = ?", data["token"]).First(&passwordReset)

	if result.Error != nil {
		fmt.Printf("Erro ao buscar token no banco: %v\n", result.Error)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid token!",
		})
	}

	// Logs para verificar a expiração
	fmt.Printf("Token encontrado com expiração: %v\n", passwordReset.ExpiresAt)
	fmt.Printf("Hora atual: %v\n", time.Now())
	fmt.Printf("Token expirado? %v\n", time.Now().After(passwordReset.ExpiresAt))

	// Verifica se o token expirou
	if time.Now().After(passwordReset.ExpiresAt) {
		fmt.Printf("Token expirado. Expiração: %v, Agora: %v\n",
			passwordReset.ExpiresAt, time.Now())
		return c.Status(400).JSON(fiber.Map{
			"message": "Token has expired",
		})
	}

	// Atualiza a senha do usuário
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		fmt.Printf("Erro ao gerar hash da senha: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error generating password hash",
		})
	}

	updateResult := database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)
	if updateResult.Error != nil {
		fmt.Printf("Erro ao atualizar senha: %v\n", updateResult.Error)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error updating password",
		})
	}

	fmt.Printf("Senha atualizada com sucesso para o email: %s\n", passwordReset.Email)

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
