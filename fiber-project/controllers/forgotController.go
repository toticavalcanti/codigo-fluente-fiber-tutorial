package controllers

import (
	"fiber-project/database"
	"fiber-project/models"
	"math/rand"
	"net/smtp"

	"github.com/gofiber/fiber/v2"
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

	from := "fluentcode@exemple.com"

	to := []string{
		data["email"],
	}

	url := "http://localhost:3000/reset/" + token

	message := []byte("Clique <a href=\"" + url + "\">aqui</a> para redefinir sua senha!")

	err := smtp.SendMail("localhost:1025", nil, from, to, message)

	if err != nil {
		return err
	}

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
