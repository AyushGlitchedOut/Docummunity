package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("-----------------DOCUMMUNITY BACKEND----------------------------")
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

	//CREATE
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

	//READ
	router.GET("/download/:ID", func(ctx *gin.Context) {
		ID := ctx.Param("ID")
		fileName := FindFileFromDirectory(ID)
		if fileName == "" {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.FileAttachment("./uploads/"+fileName, fileName)
	})

	// UPDATE

	// DELETE
	router.DELETE("/delete/:ID", func(ctx *gin.Context) {
		ID := ctx.Param("ID")
		file := "./uploads/" + FindFileFromDirectory(ID)
		err := os.Remove(file)
		if err != nil {
			if os.IsNotExist(err) {
				ctx.Status(http.StatusNotFound)
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.Status(http.StatusOK)

	})

	if err := router.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
	//
	//
	//
	//
	//
	//DONT PUT ANYTHING HERE OR BELOW HERE
	//
	//
	//

}

func FindFileFromDirectory(ID string) string {

	files, err := os.ReadDir("./uploads")

	if err != nil {
		log.Println("Something Went Wrong")
		return ""
	}
	for _, file := range files {
		fileParts := strings.Split(file.Name(), ".")
		if len(fileParts) <= 1 {
			filename := strings.Join(fileParts, "")
			if filename == ID {
				return file.Name()
			}
		}
		var fileNameArray []string
		fileNameArray = append(fileNameArray, fileParts[:len(fileParts)-1]...)
		filename := strings.Join(fileNameArray, ".")
		if filename == ID {
			return file.Name()
		}

	}
	return ""
}
