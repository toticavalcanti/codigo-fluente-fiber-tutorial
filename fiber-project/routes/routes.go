package routes

import (
	"fiber-project/controllers"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// Rotas principais da API
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/user", controllers.User)
	api.Post("/logout", controllers.Logout)
	api.Post("/forgot", controllers.Forgot)
	api.Post("/reset", controllers.Reset)

	// Rota de verificação de saúde
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"service": "auth-api",
			"version": "1.0.0",
		})
	})

	// Rota para redirecionar /reset/:token para o frontend
	app.Get("/reset/:token", func(c *fiber.Ctx) error {
		token := c.Params("token")

		// Obtém o URL do frontend da variável de ambiente FRONTEND_URL
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			return c.Status(500).SendString("FRONTEND_URL não está configurado nas variáveis de ambiente")
		}

		// Redireciona para o frontend com o token
		return c.Redirect(frontendURL+"/reset/"+token, 302)
	})
}
