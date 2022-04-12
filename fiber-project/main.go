package main

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var dsn = "root:mysql0401@/fluent_admin?charset=utf8mb4&parseTime=True&loc=Local"
	var v = "Não conseguiu conectar ao banco de dados"
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(v)
	}
	fmt.Println("Conexão OK!")

	d, e := divide(2, 0)
	fmt.Println(d, e)

	app := fiber.New()

	app.Get("/", home)

	app.Listen(":3000")
}

func home(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func divide(a int, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("you cannot divide by zero")
	}
	return a / b, nil
}
