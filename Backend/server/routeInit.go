package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"github.com/AyushGlitchedOut/Docummunity/consts"
	"github.com/AyushGlitchedOut/Docummunity/server/handlers"
	"github.com/AyushGlitchedOut/Docummunity/server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitServer(port string, db *sql.DB, firebase *firebase.App) *http.Server {
	router := gin.New()
	//MaxSize
	router.Use(handlers.MaxSizeMiddleware(consts.MaxDocumentSize + 2<<20)) //Max document size since its the biggest you can upload, and 2mb more for other details in the body
	router.MaxMultipartMemory = consts.MaxPerRequestServerMemorySize
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

	httpServer := &http.Server{
		Addr:              port,
		Handler:           router,
		MaxHeaderBytes:    256 << 10,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       2 * time.Minute,
		WriteTimeout:      0,
		IdleTimeout:       1 * time.Minute,
	}

	//Get Auth Function from Firebase App
	firebaseAuth, err := firebase.Auth(context.Background())
	if err != nil {
		log.Fatal("Error Configuring Firebase Admin SDK")
	}

	//TODO: Add rate-limiting to routes
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
