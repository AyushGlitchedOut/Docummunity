package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("DOCUMMUNITY BACKEND")
	port := ":8080"

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

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "OHHHH YEAHHHHHHHHH",
		})
	})

	router.POST("/upload", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
		}
		ID := ctx.PostForm("ID")
		fileEXT := strings.Split(file.Filename, ".")

		if len(fileEXT) <= 1 {
			ctx.SaveUploadedFile(file, "./uploads/"+ID)
			ctx.JSON(http.StatusOK, nil)
			return
		}

		ctx.SaveUploadedFile(file, "./uploads/"+ID+"."+fileEXT[len(fileEXT)-1])

		ctx.JSON(http.StatusOK, nil)
	})

	if err := router.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
