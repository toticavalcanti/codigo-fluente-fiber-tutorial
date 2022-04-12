package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var dsn = "root:mysql1234@/fluent_admin?charset=utf8mb4&parseTime=True&loc=Local"
	var v = "Não conseguiu conectar ao banco de dados"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(v)
	}
	fmt.Println("Conexão OK!")
	fmt.Println(db)

	app := fiber.New()

	app.Get("/", home)

	app.Listen(":3000")
}

func home(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
