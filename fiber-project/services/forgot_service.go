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

// formatosData contém os possíveis formatos de data que podem vir do banco
var formatosData = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.999Z",
	"2006-01-02T15:04:05-07:00",
	time.RFC3339,
	time.RFC3339Nano,
}

// parseExpirationDate tenta parsear a data usando diferentes formatos
func parseExpirationDate(dateStr string) (time.Time, error) {
	fmt.Printf("Tentando parsear data: %s\n", dateStr)

	for _, formato := range formatosData {
		parsedTime, err := time.Parse(formato, dateStr)
		if err == nil {
			fmt.Printf("Sucesso usando formato: %s\n", formato)
			// Se necessário, converter para o timezone local
			localTime := parsedTime.In(time.Local)
			return localTime, nil
		}
		fmt.Printf("Falha com formato %s: %v\n", formato, err)
	}

	return time.Time{}, fmt.Errorf("nenhum formato conhecido funcionou para a data '%s'", dateStr)
}

// RandStringRunes gera uma string aleatória para o token
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Forgot lida com a solicitação de redefinição de senha
func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("Erro ao fazer parse do body em Forgot: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Verifica se o email foi fornecido
	if data["email"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email is required",
		})
	}

	// Gera um token aleatório
	token := RandStringRunes(12)

	// Cria o registro de redefinição de senha
	passwordReset := models.PasswordReset{
		Email:     data["email"],
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	fmt.Printf("Token criado com expiração para: %v\n", passwordReset.ExpiresAt)

	// Salva o token no banco de dados
	if err := database.DB.Create(&passwordReset).Error; err != nil {
		fmt.Printf("Erro ao salvar token: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error saving token",
		})
	}

	// Configuração do email
	auth := smtp.PlainAuth("", os.Getenv("GMAIL_EMAIL"), os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com")

	to := []string{data["email"]}
	msg := []byte("To: " + data["email"] + "\r\n" +
		"Subject: Redefina sua senha\r\n" +
		"\r\n" +
		"Use o link para redefinir sua senha: " + os.Getenv("APP_URL") + "/reset/" + token + "\r\n")

	// Envia o email
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

// Reset lida com a atualização da senha
func Reset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("Erro ao fazer parse do body em Reset: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	fmt.Printf("Dados recebidos no Reset: %+v\n", data)

	// Validações básicas
	if data["token"] == "" {
		fmt.Println("Token não fornecido na requisição")
		return c.Status(400).JSON(fiber.Map{
			"message": "Token is required!",
		})
	}

	if data["password"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password is required!",
		})
	}

	if data["password"] != data["confirm_password"] {
		fmt.Println("Senhas não coincidem")
		return c.Status(400).JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	// Busca informações do token
	var resetInfo struct {
		Email     string
		Token     string
		ExpiresAt string
	}

	result := database.DB.Model(&models.PasswordReset{}).
		Select("email, token, expires_at").
		Where("token = ?", data["token"]).
		First(&resetInfo)

	if result.Error != nil {
		fmt.Printf("Erro ao buscar token no banco: %v\n", result.Error)
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid token!",
		})
	}

	fmt.Printf("Data de expiração (string): %s\n", resetInfo.ExpiresAt)

	// Parse da data de expiração
	expiresAt, err := parseExpirationDate(resetInfo.ExpiresAt)
	if err != nil {
		fmt.Printf("Erro ao parsear data de expiração: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error processing token expiration",
		})
	}

	fmt.Printf("Data de expiração (time.Time): %v\n", expiresAt)

	// Verifica se o token expirou
	if time.Now().After(expiresAt) {
		fmt.Printf("Token expirado. Expiração: %v, Agora: %v\n", expiresAt, time.Now())
		return c.Status(400).JSON(fiber.Map{
			"message": "Token has expired",
		})
	}

	// Gera o hash da nova senha
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		fmt.Printf("Erro ao gerar hash da senha: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error generating password hash",
		})
	}

	// Atualiza a senha do usuário
	updateResult := database.DB.Model(&models.User{}).
		Where("email = ?", resetInfo.Email).
		Update("password", password)

	if updateResult.Error != nil {
		fmt.Printf("Erro ao atualizar senha: %v\n", updateResult.Error)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error updating password",
		})
	}

	// Remove o token usado
	database.DB.Where("token = ?", data["token"]).Delete(&models.PasswordReset{})

	fmt.Printf("Senha atualizada com sucesso para o email: %s\n", resetInfo.Email)

	return c.JSON(fiber.Map{
		"message": "Password successfully updated!",
	})
}
