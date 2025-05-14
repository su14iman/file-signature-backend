// @title file-signer
// @version 1.0
// @description File signer and verifier API
// @host localhost:4000
// @BasePath /api
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"

	"github.com/joho/godotenv"

	"file-signer/config"
	"file-signer/models"
	"file-signer/routes"
	"file-signer/utils"
)

func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		utils.HandleError(err, "Failed to load .env file", utils.Error)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		utils.HandleError(
			fmt.Errorf("APP_PORT not set"),
			"APP_PORT not set, using default port 3000",
			utils.Error,
		)
		port = "3000"
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		utils.HandleError(
			fmt.Errorf("FRONTEND_ORIGIN not set"),
			"FRONTEND_ORIGIN not set, using default http://localhost:3000",
			utils.Error,
		)
		frontendOrigin = "http://localhost:3000"
	}

	// Start Fiber app
	app := fiber.New()

	// Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendOrigin,
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Set up static file serving
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		utils.HandleError(
			fmt.Errorf("UPLOAD_DIR not set"),
			"UPLOAD_DIR not set, using default ./uploads",
			utils.Error,
		)
	}
	app.Static("/uploads", uploadDir)

	tempDir := os.Getenv("TEMP_DIR")
	if tempDir == "" {
		utils.HandleError(
			fmt.Errorf("TEMP_DIR not set"),
			"TEMP_DIR not set, using default /tmp/file-signer",
			utils.Error,
		)
		tempDir = "/tmp/file-signer"
	}

	config.ConnectDatabase()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Document{},
	)
	config.CreateSuperAdminIfNotExists()

	if os.Getenv("ENABLE_SWAGGER") == "true" {
		app.Get("/swagger/*", swagger.HandlerDefault)
		log.Println("📚 Swagger enabled at /swagger/index.html")
	}

	routes.SetupRoutes(app, uploadDir, tempDir)

	log.Fatal(app.Listen(":" + port))

}
