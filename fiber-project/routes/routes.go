package routes

import (
	"fiber-project/controllers"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// Rotas de autenticação
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/user", controllers.User)
	api.Post("/logout", controllers.Logout)
	api.Post("/forgot", controllers.Forgot)
	api.Post("/reset", controllers.Reset)

	// Rota de verificação de saúde do serviço
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"service": "auth-api",
			"version": "1.0.0",
		})
	})

	// Rota para redirecionar o token para o frontend
	app.Get("/reset/:token", func(c *fiber.Ctx) error {
		token := c.Params("token")

		// Redireciona para o formulário de redefinição de senha no frontend
		return c.Redirect(os.Getenv("APP_URL")+"/reset/"+token, 302)
	})
}
