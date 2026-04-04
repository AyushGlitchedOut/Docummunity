package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGET(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OHHHH YEAHHHHHHHHH",
	})
}
