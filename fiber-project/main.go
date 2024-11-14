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
			return origin != ""
		},
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Middleware para servir arquivos estáticos
	app.Static("/", "./frontend/react-auth/build")

	// Middleware de fallback para React Router
	app.Use(func(c *fiber.Ctx) error {
		if err := c.SendFile("./frontend/react-auth/build/index.html"); err != nil {
			return c.Status(404).SendString("Page not found")
		}
		return nil
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
