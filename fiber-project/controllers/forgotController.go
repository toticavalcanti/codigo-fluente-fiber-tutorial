package controllers

import (
	"fiber-project/services"

	"github.com/gofiber/fiber/v2"
)

func Forgot(c *fiber.Ctx) error {
	return services.Forgot(c)
}

func Reset(c *fiber.Ctx) error {
	return services.Reset(c)
}
