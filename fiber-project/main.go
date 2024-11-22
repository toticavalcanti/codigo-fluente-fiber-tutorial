package main

import (
	"fiber-project/database"
	"fiber-project/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	app := fiber.New()

	// Configuração de CORS dinâmica
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			// Permite todas as origens válidas (não vazias)
			return origin != ""
		},
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Adiciona o middleware de redirecionamento para reset de senha
	app.Use("/api/reset/*", func(c *fiber.Ctx) error {
		if c.Method() == "GET" {
			path := c.Path() // Pega o caminho completo incluindo /reset/ e o token
			frontendURL := os.Getenv("FRONTEND_URL")
			return c.Redirect(frontendURL+path, 301)
		}
		return c.Next()
	})

	// Configuração das rotas
	routes.Setup(app)

	// Configuração da porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Inicialização do servidor
	app.Listen(":" + port)
}
