package database

import (
	"fiber-project/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Carrega o arquivo .env em ambiente de desenvolvimento (não em produção)
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente padrão.")
		}
	}

	// Obtém o DSN do banco de dados da variável de ambiente
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("Erro: A variável de ambiente DB_DSN não está configurada")
	}

	// Inicializa a conexão com o banco de dados
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	DB = connection

	// Realiza a migração automática dos modelos
	migrate(connection)
}

// Função auxiliar para realizar a migração dos modelos
func migrate(connection *gorm.DB) {
	err := connection.AutoMigrate(
		&models.User{},
		&models.PasswordReset{},
	)
	if err != nil {
		log.Fatalf("Erro ao realizar migração automática: %v", err)
	}

	log.Println("Migração do banco de dados realizada com sucesso!")
}
