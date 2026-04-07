package auth

import (
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(firebaseAuth *auth.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := firebaseAuth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		fmt.Println(token)
		ctx.Set("tokenUID", user.UID)

		ctx.Next()
	}
}
