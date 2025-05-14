package routes

import (
	"file-signer/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, uploadDir string, tempDir string) {
	api := app.Group("/api")

	api.Post("/admin/login", controllers.LoginAdmin)

	api.Post("/sign", controllers.SignFileHandler(uploadDir))
	api.Post("/verify/id", controllers.VerifyFileByIdHandler)
	api.Post("/verify/file", controllers.VerifyFileByUploadHandler(tempDir))

}
