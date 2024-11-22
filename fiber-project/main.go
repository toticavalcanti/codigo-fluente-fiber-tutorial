package main

import (
	"fiber-project/database"
	"fiber-project/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Conectar ao banco de dados
	database.Connect()

	// Inicializar o Fiber
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
