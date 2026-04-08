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

	//FREE Routes
	freeRoutes := router.Group("/api")
	{
		freeRoutes.GET("/", HandlePING)
		freeRoutes.GET("/user/SEARCH", HandleUserSEARCH(db))
		freeRoutes.GET("/user/GET", HandleUserGET(db))
		freeRoutes.GET("/data/SEARCH", HandleDataSEARCH(db))
		freeRoutes.GET("/data/GET", HandleDataGET(db))
	}

	//User Routes
	userRoutes := router.Group("/api/user")
	userRoutes.Use(auth.AuthMiddleware(firebaseAuth))
	{
		userRoutes.GET("/ACCOUNT", HandleUserACCOUNT(db))
		userRoutes.POST("/CREATE", HandleUserCREATE(db))
		userRoutes.PATCH("/UPDATE", HandleUserUPDATE(db))
		userRoutes.DELETE("/DELETE", HandleUserDELETE(db))
		userRoutes.DELETE("/DELETE/FULL", HandleUserDELETEWithRecords(db))
	}

	//Data Routes
	dataRoutes := router.Group("/api/data")
	dataRoutes.Use(auth.AuthMiddleware(firebaseAuth))
	{
		dataRoutes.POST("/CREATE", HandleDataCREATE(db))
		dataRoutes.PATCH("/UPDATE", HandleDataUPDATE(db))
		dataRoutes.DELETE("/DELETE", HandleDataDELETE(db))
	}

	//Test Routes
	router.POST("/test", auth.AuthMiddleware(firebaseAuth), VerifyTest)

	return router

}
