package utilities

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ParseToken(ctx *gin.Context) (string, error) {
	UID, exists := ctx.Get("tokenUID")
	if !exists {
		return "", fmt.Errorf("Error Parsing Token")
	}
	UIDstr, ok := UID.(string)
	if !ok {
		return "", fmt.Errorf("No UID found")
	}

	return UIDstr, nil
}
