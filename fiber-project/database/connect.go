package database

import (
	"fmt"

	"fiber-project/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() {
	var dsn = "toticavalcanti:mysql1234@/fluent_admin?charset=utf8mb4&parseTime=True&loc=Local"
	var v = "Não conseguiu conectar ao banco de dados"
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(v)
	}

	connection.AutoMigrate(&models.User{})
	fmt.Println("Conexão OK!")
}
