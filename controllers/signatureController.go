package controllers

import (
	"file-signer/utils"
	"fmt"
	"os"
	"os/exec"

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

func SignAndEmbed(filePath string, id string) (string, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", utils.HandleError(err, "Failed to read file", utils.Error)
	}

	signature, err := signFile(data)
	if err != nil {
		return "", utils.HandleError(err, "Failed to sign file", utils.Error)
	}

	err = embedSignatureFile(filePath, id, signature)

	if err != nil {
		return "", utils.HandleError(err, "Failed to embed signature", utils.Error)
	}

	return signature, nil
}

func signFile(data []byte) (string, error) {

	privateKey, _ := PrivateKey()

	hashed := sha256.Sum256(data)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, cryptoHash(), hashed[:])
	if err != nil {
		return "", utils.HandleError(err, "Failed to sign", utils.Error)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func cryptoHash() crypto.Hash {
	return crypto.SHA256
}

func embedSignatureFile(filePath, id, signature string) error {

	comment := fmt.Sprintf("ID:%s;SIG:%s", id, signature)

	cmd := exec.Command("exiftool", "-all=", "-overwrite_original", filePath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		if err := exec.Command("rm", filePath).Run(); err != nil {
			return utils.HandleError(err, "Failed to remove file", utils.Error)
		}
		return utils.HandleError(err, "Failed to clean up Exif", utils.Error)
	}

	cmd = exec.Command("exiftool", "-overwrite_original", "-UserComment="+comment, filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return utils.HandleError(err, fmt.Sprintf("Failed to write Exif: %s", string(output)), utils.Error)
	}

	return nil
}
