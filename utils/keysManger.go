package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func loadPublicKey() (*rsa.PublicKey, error) {

	keyPath := os.Getenv("PUBLIC_KEY_PATH")
	keyBytes, err := os.ReadFile(keyPath)

	if err != nil {
		// return nil, fmt.Errorf("failed to read public key: %w", err)
		return nil, HandleError(err, "Failed to read public key", Error)
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		// return nil, fmt.Errorf("invalid public key format")
		return nil, HandleError(err, "Invalid public key format", Error)
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// return nil, fmt.Errorf("failed to parse public key: %w", err)
		return nil, HandleError(err, "Failed to parse public key", Error)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		// return nil, fmt.Errorf("public key is not of type rsa.PublicKey")
		return nil, HandleError(err, "Public key is not of type rsa.PublicKey", Error)
	}
	return rsaPublicKey, nil
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	keyPath := os.Getenv("PRIVATE_KEY_PATH")
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		// return nil, fmt.Errorf("failed to read private key: %w", err)
		return nil, HandleError(err, "Failed to read private key", Error)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		// return nil, fmt.Errorf("invalid private key format")
		return nil, HandleError(err, "Invalid private key format", Error)
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// return nil, fmt.Errorf("failed to parse private key: %w", err)
		return nil, HandleError(err, "Failed to parse private key", Error)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		// return nil, fmt.Errorf("private key is not of type rsa.PrivateKey")
		return nil, HandleError(err, "Private key is not of type rsa.PrivateKey", Error)
	}

	return rsaPrivateKey, nil
}
