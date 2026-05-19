package authUtils

import (
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// Authorization Middleware attached to all Routes that require verification.
// Passes the UID obtained from JWT as key/value pair in ctx with key "tokenUID"
func AuthMiddleware(firebaseAuth *auth.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Get JWT header
		authHeader := ctx.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		//verify JWT
		user, err := firebaseAuth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Access Denied",
			})
			ctx.Abort()
			return
		}

		//DEBUG: REMOVE IN PRODUCTION
		log.Println(token)

		//Set the tokenUID in context
		ctx.Set("tokenUID", user.UID)
		ctx.Next()
	}
}
