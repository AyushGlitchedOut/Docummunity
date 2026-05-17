package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Very Simple Middlware used on the base App to limit the maxSize of a Request
func MaxSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)
		ctx.Next()
	}
}
