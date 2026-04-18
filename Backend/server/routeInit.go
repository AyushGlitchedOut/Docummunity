package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/AyushGlitchedOut/Docummunity/server/consts"
	"github.com/AyushGlitchedOut/Docummunity/server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const ()

// TODO: Handle Timeouts by using a http.Server and attaching gin's Eouter to it, with configuration for large upload and all.
func InitServer(port string, db *sql.DB, firebase *firebase.App) *http.Server {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	//TEMPORARY CORS POLICY! REMOVE IN PRODUCTION
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization"},
	}))
	router.MaxMultipartMemory = consts.MaxPerRequestServerMemorySize

	httpServer := &http.Server{
		Addr:    port,
		Handler: router,
	}

	//Get Auth Function from Firebase App
	firebaseAuth, err := firebase.Auth(context.Background())
	if err != nil {
		log.Fatal("Error Configuring Firebase Admin SDK")
	}

	//Routing system
	serverRoutes := router.Group("/api")
	{
		//USER ROUTES
		routes.CreateUserRoutes(serverRoutes, firebaseAuth, db)

		//DATA Routes
		routes.CreateDataRoutes(serverRoutes, firebaseAuth, db)

		//Other routes
		routes.CreateOtherRoutes(serverRoutes, firebaseAuth)

	}

	return httpServer

}
