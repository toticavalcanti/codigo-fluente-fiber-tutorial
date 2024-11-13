package main

import (
	"fiber-project/database"
	"fiber-project/routes"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Conecta ao banco de dados
	database.Connect()
	app := fiber.New()

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			origin := os.Getenv("CORS_ORIGIN")
			if origin == "" {
				return "*"
			}
			return origin
		}(),
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	// Middleware para logar o cabeçalho "Origin"
	app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("Request Origin: %s\n", c.Get("Origin"))
		return c.Next()
	})

	// Configuração de rotas
	routes.Setup(app)

	// Configuração da porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Inicia o servidor
	app.Listen(":" + port)
}
