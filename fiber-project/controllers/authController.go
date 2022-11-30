package controllers

import (
	"fiber-project/database"
	"fiber-project/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
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

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	database.DB.Create(&user)
	return c.JSON(user)
}
