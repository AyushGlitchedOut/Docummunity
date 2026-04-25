package routes

import (
	"database/sql"

	"firebase.google.com/go/auth"
	"github.com/AyushGlitchedOut/Docummunity/authUtils"
	"github.com/AyushGlitchedOut/Docummunity/server/handlers"
	"github.com/gin-gonic/gin"
)

func CreateUserRoutes(group *gin.RouterGroup, firebaseAuth *auth.Client, db *sql.DB) {
	userRoutes := group.Group("/user")
	{
		//Auth Required
		protectedUserRoutes := userRoutes.Group("/")
		protectedUserRoutes.Use(authUtils.AuthMiddleware(firebaseAuth))
		{
			protectedUserRoutes.GET("/ACCOUNT", handlers.HandleUserACCOUNT(db))
			protectedUserRoutes.POST("/CREATE", handlers.HandleUserCREATE(db))
			protectedUserRoutes.PATCH("/UPDATE", handlers.HandleUserUPDATE(db))
			protectedUserRoutes.DELETE("/DELETE", handlers.HandleUserDELETE(db, true))
			protectedUserRoutes.DELETE("/DELETE/FULL", handlers.HandleUserDELETE(db, false))
		}

		//Free Routes
		freeUserRoutes := userRoutes.Group("/")
		{
			freeUserRoutes.GET("/SEARCH/:query", handlers.HandleUserSEARCH(db))
			freeUserRoutes.GET("/GET/:uid", handlers.HandleUserGET(db))
			freeUserRoutes.GET("/RECORDS/:uid", handlers.HandleUserRecordsGET(db))
			freeUserRoutes.GET("/PROFILE_PIC/:filename", handlers.HostUserPROFILE_PIC())
		}

	}
}
