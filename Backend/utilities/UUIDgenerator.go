package utilities

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	ID, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}
	return ID.String()
}
