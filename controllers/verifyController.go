package controllers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"file-signer/config"
	"file-signer/models"
	"file-signer/utils"

	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

func VerifyFile(filePath string) (bool, string, error) {

	cmd := exec.Command("exiftool", "-UserComment", filePath)
	output, err := cmd.Output()
	if err != nil {
		return false, "", utils.HandleError(err, fmt.Sprintf("Failed to extract metadata: %s", string(output)), utils.Error)
	}

	text := string(bytes.TrimSpace(output))

	if !strings.Contains(text, "ID:") || !strings.Contains(text, "SIG:") {
		return false, "", utils.HandleError(err, "Invalid signature format in file metadata", utils.Error)
	}

	parts := strings.Split(text, ";")

	id := strings.Split(
		strings.TrimPrefix(strings.TrimSpace(parts[0]), "ID:"),
		":",
	)[2]
	sigEncoded := strings.TrimPrefix(strings.TrimSpace(parts[1]), "SIG:")
	sigBytes, err := base64.StdEncoding.DecodeString(sigEncoded)

	if err != nil {
		return false, "", utils.HandleError(err, "Failed to decode signature", utils.Error)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, "", utils.HandleError(err, "Failed to read file", utils.Error)
	}

	hashed := sha256.Sum256(data)

	publicKey, err := PublicKey()
	if err != nil {
		return false, "", utils.HandleError(err, "Failed to load public key", utils.Error)
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], sigBytes)
	if err != nil {
		return false, id, nil
	}

	return true, id, nil
}

func VerifyID(id string) (models.DocumentInfo, error) {

	var doc models.Document
	result := config.DB.First(&doc, "id = ?", id)
	if result.Error != nil {
		return models.DocumentInfo{}, utils.HandleError(result.Error, fmt.Sprintf("Document with ID %s not found in the database", id), utils.Error)
	}

	filePath := doc.FilePath
	if filePath == "" {
		return models.DocumentInfo{}, utils.HandleError(fmt.Errorf("file path not found for document ID %s", id), "File path not found", utils.Error)
	}

	fileType := doc.FileType
	if fileType == "" {
		return models.DocumentInfo{}, utils.HandleError(fmt.Errorf("file type not found for document ID %s", id), "File type not found", utils.Error)
	}

	fileSignature := doc.Signature
	if fileSignature == "" {
		return models.DocumentInfo{}, utils.HandleError(fmt.Errorf("file signature not found for document ID %s", id), "File signature not found", utils.Error)
	}

	createAt := doc.CreatedAt
	if createAt.IsZero() {
		return models.DocumentInfo{}, utils.HandleError(fmt.Errorf("created at not found for document ID %s", id), "Created at not found", utils.Error)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return models.DocumentInfo{}, utils.HandleError(err, fmt.Sprintf("File does not exist at path %s", filePath), utils.Error)
	}

	return models.DocumentInfo{
		FilePath:      filePath,
		FileType:      fileType,
		FileSignature: fileSignature,
		CreatedAt:     createAt,
	}, nil
}
