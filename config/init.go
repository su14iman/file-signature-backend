package config

import (
	"file-signer/models"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func CreateSuperAdminIfNotExists() {
	email := os.Getenv("SUPERADMIN_EMAIL")
	password := os.Getenv("SUPERADMIN_PASSWORD")

	if email == "" || password == "" {
		fmt.Println("⚠️  SUPERADMIN credentials not found in .env")
		return
	}

	var existing models.User
	if err := DB.Where("email = ?", email).First(&existing).Error; err == nil {
		fmt.Println("✅ SuperAdmin already exists")
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	super := models.User{
		Name:     "Super Admin",
		Email:    email,
		Password: string(hashed),
		Role:     "superadmin",
	}

	DB.Create(&super)
	fmt.Println("✅ SuperAdmin created from env")
}
