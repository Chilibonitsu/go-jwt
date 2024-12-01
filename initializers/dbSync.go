package initializers

import (
	"fmt"
	"go-jwt/models"
	"log"
)

func SyncDB() {
	if DB == nil {
		log.Fatal("Database connection is nil. Please check your connection.")
	}
	DB.AutoMigrate(&models.User{})
	var user models.User

	result := DB.First(&user)

	if result.Error != nil {
		fmt.Println("Error fetching user:", result.Error)
		return
	}

	fmt.Println("User data:", user)
}
