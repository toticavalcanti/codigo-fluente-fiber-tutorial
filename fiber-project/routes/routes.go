package routes

import (
	"fiber-project/controllers"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Rota de reset com status code 307
	app.Get("/reset/:token", func(c *fiber.Ctx) error {
		token := c.Params("token")
		frontendURL := os.Getenv("FRONTEND_URL")
		return c.Redirect(frontendURL+"/#/reset/"+token, 307)
	})

	// Grupo de API
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/user", controllers.User)
	api.Post("/logout", controllers.Logout)
	api.Post("/forgot", controllers.Forgot)
	api.Post("/reset", controllers.Reset)

	// Health check route
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"service": "auth-api",
			"version": "1.0.0",
		})
	})
}
