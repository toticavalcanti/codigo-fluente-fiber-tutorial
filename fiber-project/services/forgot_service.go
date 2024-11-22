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

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}

	database.DB.Create(&passwordReset)

	// Configuração de autenticação SMTP para Gmail
	auth := smtp.PlainAuth("", os.Getenv("GMAIL_EMAIL"), os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com")

	to := []string{data["email"]}
	msg := []byte("To: " + data["email"] + "\r\n" +
		"Subject: Redefina sua senha\r\n" +
		"\r\n" +
		"Use o link para redefinir sua senha: " + os.Getenv("FRONTEND_URL") + "/reset/" + token + "\r\n")

	// Envia o email usando o servidor SMTP do Gmail
	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("GMAIL_EMAIL"), to, msg)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Reset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["confirm_password"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	var passwordReset = models.PasswordReset{}
	if err := database.DB.Where("token = ?", data["token"]).Last(&passwordReset); err.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid token!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	return c.JSON(fiber.Map{
		"message": "success",
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
