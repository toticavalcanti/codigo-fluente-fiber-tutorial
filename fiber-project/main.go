package main

import (
	"fiber-project/database"
	"fiber-project/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Conectar ao banco de dados
	database.Connect()

	// Inicializar o Fiber
	app := fiber.New()

	// Obter a variável de ambiente APP_URL
	frontendURL := os.Getenv("APP_URL")
	if frontendURL == "" {
		log.Fatal("APP_URL environment variable is required")
	} else {
		log.Printf("APP_URL is set to: %s", frontendURL)
	}

	// Configuração de CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
	}))

	// Configurar rotas
	routes.Setup(app)

	// Definir a porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT is not set. Using default port: %s", port)
	}

	// Iniciar o servidor
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
