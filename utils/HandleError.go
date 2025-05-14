package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	Info    = "INFO"
	Warning = "WARNING"
	Error   = "ERROR"
)

func HandleError(err error, customMessage string, level string) error {

	if err != nil {

		logMessage := fmt.Sprintf(
			"\n%s [%s] - %s\n%v",
			time.Now().Format(time.RFC3339),
			level,
			customMessage,
			err,
		)

		if os.Getenv("ENABLE_LOGGING") == "true" {
			log.Println(logMessage)
		}

		// return error
		return fmt.Errorf(
			"%s: %w",
			customMessage,
			err,
		)

	}
	return nil
}

func MainErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error":  err.Error(),
		"status": code,
	})
}
