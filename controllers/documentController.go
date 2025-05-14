package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"file-signer/config"
	"file-signer/models"
	"file-signer/utils"
)

func SignFileHandler(uploadDir string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, savePath, ext, id, fileName, err := UploadFile(c, uploadDir)
		if err != nil {
			return utils.HandleError(err, "Failed to upload file", utils.Error)
		}

		signature, err := SignAndEmbed(savePath, id)
		if err != nil {
			return utils.HandleError(err, "Failed to sign and embed", utils.Error)
		}

		if strings.HasSuffix(strings.ToLower(ext), ".jpg") || strings.HasSuffix(strings.ToLower(ext), ".jpeg") || strings.HasSuffix(strings.ToLower(ext), ".png") {
			err := utils.AddIDToImage(savePath, id, signature)
			if err != nil {
				return utils.HandleError(err, "Failed to add ID to image", utils.Error)
			}
		}

		if strings.HasSuffix(strings.ToLower(ext), ".pdf") {
			err := utils.AddIDToPDF(savePath, id, signature)
			if err != nil {
				return utils.HandleError(err, "Failed to add ID to PDF", utils.Error)
			}
		}

		doc := models.Document{
			ID:        id,
			FileName:  fileName,
			FilePath:  savePath,
			Signature: signature,
			FileType:  strings.ToLower(ext),
			CreatedAt: time.Now(),
		}

		config.DB.Create(&doc)

		return c.JSON(fiber.Map{
			"message":   "File uploaded and signed successfully",
			"id":        doc.ID,
			"signature": doc.Signature,
		})
	}
}

func VerifyFileByIdHandler(c *fiber.Ctx) error {

	id := c.FormValue("id")
	if id == "" {
		return utils.HandleError(fmt.Errorf("ID is required"), "ID is required", utils.Error)
	}

	docInfo, err := VerifyID(id)
	if err != nil {
		return utils.HandleError(err, "Failed to find document by ID", utils.Error)
	}

	return c.JSON(fiber.Map{
		"filePath":      docInfo.FilePath,
		"fileType":      docInfo.FileType,
		"fileSignature": docInfo.FileSignature,
		"createdAt":     docInfo.CreatedAt,
	})
}

func VerifyFileByUploadHandler(tempDir string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uploadedFile, savePath, _, _, _, err := UploadFile(c, tempDir)
		if err != nil {
			return utils.HandleError(err, "Failed to upload file", utils.Error)
		}

		defer uploadedFile.Close()

		_, signatureID, err := VerifyFile(savePath)
		if err != nil {
			RemoveFile(savePath)
			return utils.HandleError(err, "Failed to verify signature", utils.Error)

		}

		err = RemoveFile(savePath)
		if err != nil {
			return utils.HandleError(err, "Failed to remove temp file", utils.Error)
		}

		docInfo, err := VerifyID(signatureID)
		if err != nil {
			return utils.HandleError(err, "Failed to find document by ID", utils.Error)
		}

		return c.JSON(fiber.Map{
			"filePath":      docInfo.FilePath,
			"fileType":      docInfo.FileType,
			"fileSignature": docInfo.FileSignature,
			"createdAt":     docInfo.CreatedAt,
		})

	}
}
