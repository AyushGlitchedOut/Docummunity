package server

import (
	"context"
	"database/sql"
	"log"

	firebase "firebase.google.com/go"
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
	authRequired.Use(authMiddleware(firebaseAuth))
	{
		//CREATE
		authRequired.POST("/upload", HandleCREATE(db))
		//READ
		authRequired.GET("/download/:ID", HandleRead(db))
		// UPDATE
		authRequired.PUT("/update/:ID", HandleUpdate(db))
		// DELETE
		authRequired.DELETE("/delete/:ID", HandleDelete(db))
		//Verification Prototype
		authRequired.GET("/verify", verifyLogin(firebaseAuth))
	}

	return router

}
