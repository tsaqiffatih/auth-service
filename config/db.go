package config

import (
	"log"
	"os"

	"github.com/tsaqiffatih/auth-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	db.AutoMigrate(&models.User{}, &models.Token{})

	return db
}

// Function Seeding data to database
func SeedDatabase(db *gorm.DB) {
	users := []models.User{
		{Name: "Budi Susanto", Email: "budiSusanto@mail.com"},
		{Name: "Sangkuriang Joko", Email: "sangkuriang@mail.com"},
	}

	for _, user := range users {
		result := db.Create(&user)
		if result.Error != nil {
			log.Println("Error seeding user: ", result.Error)
		} else {
			log.Println("User seeded: ", user)
		}
	}
}
