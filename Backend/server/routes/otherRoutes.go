package routes

import (
	"firebase.google.com/go/auth"
	"github.com/AyushGlitchedOut/Docummunity/authUtils"
	"github.com/AyushGlitchedOut/Docummunity/server/handlers"
	"github.com/gin-gonic/gin"
)

func CreateOtherRoutes(group *gin.RouterGroup, firebaseAuth *auth.Client) {
	otherRoutes := group.Group("/")
	{
		//Auth Required
		protectedOtherRoutes := otherRoutes.Group("/")
		protectedOtherRoutes.Use(authUtils.AuthMiddleware(firebaseAuth))
		{
			protectedOtherRoutes.POST("/test", handlers.VerifyTest)
		}

		//Free Routes
		freeOtherRoutes := otherRoutes.Group("/")
		{
			freeOtherRoutes.GET("/", handlers.HandlePING)
		}

	}
}
