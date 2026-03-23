package routes

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitServer(port string, db *sql.DB) *gin.Engine {
	router := gin.Default()

	//TEMPORARY CORS POLICY! REMOVE IN PRODUCTION
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization"},
	}))

	router.GET("/", HandleGET)
	//CREATE
	router.POST("/upload", HandleCREATE(db))
	//READ
	router.GET("/download/:ID", HandleRead(db))
	// UPDATE
	router.PUT("/update/:ID", HandleUpdate(db))
	// DELETE
	router.DELETE("/delete/:ID", HandleDelete(db))

	return router

}
