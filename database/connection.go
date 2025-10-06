package database

import (
	"fmt"
	"log"
	"os"

	"github.com/henriquefac/backend_go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=cerberus_db port=%s sslmode=disable",
		host, user, pass, port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	database.AutoMigrate(
		&models.User{},
		&models.Friendship{},
		&models.Achievement{},
		&models.MissingAnimal{},
		&models.AnimalSpottedRegister{},
		&models.AnimalReturnedRegister{},
	)

	DB = database
	log.Println("Conecato ao banco e migração concluída")
}
