package utils

import (
	"fmt"
	"log"
	"time"
)

const (
	Info    = "INFO"
	Warning = "WARNING"
	Error   = "ERROR"
)

func HandleError(err error, customMessage string, level string) error {

	if err != nil {
		// Log the error with a timestamp
		logMessage := fmt.Sprintf(
			"\n%s [%s] - %s\n%v",
			time.Now().Format(time.RFC3339),
			level,
			customMessage,
			err,
		)
		log.Println(logMessage)

		// return error
		return fmt.Errorf(
			"%s: %w",
			customMessage,
			err,
		)

	}
	return nil
}

func HandleErrorJSON(err error, customMessage string, level string) map[string]interface{} {
	if err != nil {
		// Log the error with a timestamp
		logMessage := fmt.Sprintf(
			"\n%s [%s] - %s\n%v",
			time.Now().Format(time.RFC3339),
			level,
			customMessage,
			err,
		)
		log.Println(logMessage)

		return map[string]interface{}{
			"error":   customMessage,
			"message": err.Error(),
		}
	}
	return nil
}
