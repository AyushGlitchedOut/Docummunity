package routes

import (
	"database/sql"

	"firebase.google.com/go/auth"
	"github.com/AyushGlitchedOut/Docummunity/authUtils"
	"github.com/AyushGlitchedOut/Docummunity/server/handlers"
	"github.com/gin-gonic/gin"
)

func CreateDataRoutes(group *gin.RouterGroup, firebaseAuth *auth.Client, db *sql.DB) {
	dataRoutes := group.Group("/data")
	{
		//Auth reuired
		protectedDataRoutes := dataRoutes.Group("/")
		protectedDataRoutes.Use(authUtils.AuthMiddleware(firebaseAuth))
		{
			protectedDataRoutes.POST("/CREATE", handlers.HandleDataCREATE(db))
			protectedDataRoutes.PATCH("/UPDATE", handlers.HandleDataUPDATE(db))
			protectedDataRoutes.DELETE("/DELETE/:uuid", handlers.HandleDataDELETE(db))
		}

		//Free routes
		freeDataRoutes := dataRoutes.Group("/")
		{
			freeDataRoutes.GET("/SEARCH/:query", handlers.HandleDataSEARCH(db))
			freeDataRoutes.GET("/GET/:uuid", handlers.HandleDataGET(db))
			freeDataRoutes.GET("/PREVIEW/:filename", handlers.HostDataPreview())
			freeDataRoutes.GET("/FILE/:filename", handlers.HostDataFiles())
		}
	}
}
