package controllers

import (
	"file-signer/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func UploadFile(c *fiber.Ctx, uploadDir string) (*os.File, string, string, string, string, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, "", "", "", "", utils.HandleError(err, "Error loading .env file", utils.Error)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, "", "", "", "", utils.HandleError(err, "Invalid file", utils.Error)
	}

	fileName := fileHeader.Filename
	id := uuid.New().String()

	if uploadDir == "" {
		uploadDir = os.Getenv("UPLOAD_DIR")
		if uploadDir == "" {
			return nil, "", "", "", "", utils.HandleError(err, "UPLOAD_DIR not set", utils.Error)
		}
	}

	ext := filepath.Ext(fileHeader.Filename)
	savePath := fmt.Sprintf("%s/%s%s", uploadDir, id, ext)

	if err := c.SaveFile(fileHeader, savePath); err != nil {
		return nil, "", "", "", "", utils.HandleError(err, "Failed to save file", utils.Error)
	}

	uploadedFile, err := os.Open(savePath)
	if err != nil {
		return nil, "", "", "", "", utils.HandleError(err, "Failed to open file", utils.Error)
	}

	return uploadedFile, savePath, ext, id, fileName, nil
}

func RemoveFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return utils.HandleError(err, "Failed to remove file", utils.Error)
	}
	return nil
}
