package utilities

import (
	"fmt"

	"github.com/google/uuid"
)

// Simple Function to generate a Time-based UUID v7
func GenerateUUID() string {
	ID, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}
	return ID.String()
}
