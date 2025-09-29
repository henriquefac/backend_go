package database

import (
	"github.com/henriquefac/backend_go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=seu_usuario password=sua_senha dbname=albums_db port=5432 sslmode=disable"
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
}
