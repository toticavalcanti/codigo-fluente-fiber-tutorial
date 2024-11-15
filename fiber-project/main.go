// main.go
package main

import (
	"fiber-project/database"
	"fiber-project/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Inicializar conexão com banco de dados
	database.Connect()

	app := fiber.New(fiber.Config{
		// Aumentar o limite do body parser para lidar com uploads maiores
		BodyLimit: 10 * 1024 * 1024,
	})

	// Logger middleware para debugging
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	// Configuração CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("FRONTEND_URL"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
	}))

	// Configuração das rotas
	routes.Setup(app)

	// Porta do servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Logging da URL do frontend permitida
	log.Printf("Allowed Frontend URL: %s\n", os.Getenv("FRONTEND_URL"))

	// Iniciar servidor
	log.Printf("Server starting on port %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
