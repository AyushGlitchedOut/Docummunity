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

func InitServer(port string, db *sql.DB, firebase *firebase.App) *gin.Engine {
	router := gin.Default()

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
	freeRoutes := router.Group("/")
	{

		freeRoutes.GET("/", HandleGET)
	}

	//Secure Routes
	authRequired := router.Group("/")
	authRequired.Use(auth.AuthMiddleware(firebaseAuth))
	{

	}

	return router

}
