package server

import (
	"context"
	"database/sql"
	"log"

	firebase "firebase.google.com/go"
	"github.com/AyushGlitchedOut/Docummunity/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	maxPerRequestServerMemorySize = 64 << 20

	maxPictureSize  = 2 << 20  //2mb
	maxDocumentSize = 40 << 20 //40mb
)

// TODO: Handle Timeouts by using a http.Server and attaching gin's Eouter to it, with configuration for large upload and all.
func InitServer(port string, db *sql.DB, firebase *firebase.App) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.MaxMultipartMemory = maxPerRequestServerMemorySize

	//Get Auth Function from Firebase App
	firebaseAuth, err := firebase.Auth(context.Background())
	if err != nil {
		log.Fatal("Error Configuring Firebase Admin SDK")
	}

	//TEMPORARY CORS POLICY! REMOVE IN PRODUCTION
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization"},
	}))

	//Routing system
	serverRoutes := router.Group("/api")
	{
		//USER ROUTES
		userRoutes := serverRoutes.Group("/user")
		{
			//Auth Required
			protectedUserRoutes := userRoutes.Group("/")
			protectedUserRoutes.Use(auth.AuthMiddleware(firebaseAuth))
			{
				protectedUserRoutes.GET("/ACCOUNT", HandleUserACCOUNT(db))
				protectedUserRoutes.POST("/CREATE", HandleUserCREATE(db))
				protectedUserRoutes.PATCH("/UPDATE", HandleUserUPDATE(db))
				protectedUserRoutes.DELETE("/DELETE", HandleUserDELETE(db))
				protectedUserRoutes.DELETE("/DELETE/FULL", HandleUserDELETEWithRecords(db))
			}

			//Free Routes
			freeUserRoutes := userRoutes.Group("/")
			{
				freeUserRoutes.GET("/SEARCH/:query", HandleUserSEARCH(db))
				freeUserRoutes.GET("/GET/:uid", HandleUserGET(db))
				freeUserRoutes.GET("/RECORDS/:uid", HandleUserRecordsGET(db))
				freeUserRoutes.GET("/PROFILE_PIC/:filename", HostUserPROFILE_PIC())
			}

		}

		//DATA Routes
		dataRoutes := serverRoutes.Group("/data")
		{
			//Auth reuired
			protectedDataRoutes := dataRoutes.Group("/")
			protectedDataRoutes.Use(auth.AuthMiddleware(firebaseAuth))
			{
				protectedDataRoutes.POST("/CREATE", HandleDataCREATE(db))
				protectedDataRoutes.PATCH("/UPDATE", HandleDataUPDATE(db))
				protectedDataRoutes.DELETE("/DELETE", HandleDataDELETE(db))
			}

			//Free routes
			freeDataRoutes := dataRoutes.Group("/")
			{
				freeDataRoutes.GET("/SEARCH/:query", HandleDataSEARCH(db))
				freeDataRoutes.GET("/GET/:uuid", HandleDataGET(db))
				freeDataRoutes.GET("/PREVIEW/:filename", HostDataPreview())
				freeDataRoutes.GET("/FILE/:filename", HostDataFiles())
			}
		}

		//Other routes
		otherRoutes := serverRoutes.Group("/")
		{
			//Auth Required
			protectedOtherRoutes := otherRoutes.Group("/")
			protectedOtherRoutes.Use(auth.AuthMiddleware(firebaseAuth))
			{
				protectedOtherRoutes.POST("/test", VerifyTest)
			}

			//Free Routes
			freeOtherRoutes := otherRoutes.Group("/")
			{
				freeOtherRoutes.GET("/", HandlePING)
			}

		}

	}

	return router

}
