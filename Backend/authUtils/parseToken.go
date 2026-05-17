package authUtils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Obtains the tokenUID stored by authMiddleware in the context, and returns the UID (if found) as string
func ParseToken(ctx *gin.Context) (string, error) {
	UID, exists := ctx.Get("tokenUID")

	//Check if the tokenUID field exists or not
	if !exists {
		return "", fmt.Errorf("No UID Found")
	}

	//Convert tokenUID from "any" to "string"
	UIDstr, ok := UID.(string)
	if !ok {
		return "", fmt.Errorf("Error Parsing Token")
	}

	return UIDstr, nil
}
