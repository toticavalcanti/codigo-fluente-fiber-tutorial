package controllers

import (
	"fiber-project/services"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	return services.Register(c)
}

func Login(c *fiber.Ctx) error {
	return services.Login(c)
}

func User(c *fiber.Ctx) error {
	return services.User(c)
}

func Logout(c *fiber.Ctx) error {
	return services.Logout(c)
}
