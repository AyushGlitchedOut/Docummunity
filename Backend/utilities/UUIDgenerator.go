package utilities

import (
	"log"

	"github.com/google/uuid"
)

// Simple Function to generate a Time-based UUID v7
func GenerateUUID() string {
	ID, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
	}
	return ID.String()
}
