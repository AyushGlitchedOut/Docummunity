package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/utilities"
	"github.com/gin-gonic/gin"
)

func HandleGET(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OHHHH YEAHHHHHHHHH",
	})
}

func HandleCREATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data dbUtils.TestData
		ID := ctx.PostForm("ID")
		file, err := ctx.FormFile("file")
		if err != nil {
			log.Println(err)
			return
		}

		path := "./uploads/" + file.Filename

		err = ctx.SaveUploadedFile(file, path)
		if err != nil {
			log.Println(err)
			return
		}

		data = dbUtils.TestData{
			ID:       int(time.Now().Unix()),
			NAME:     ID,
			FILEPATH: path,
		}

		err = dbUtils.InsertData(data, db)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func HandleRead(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Param("ID")
		fileName := utilities.FindFileFromDirectory(ID)
		if fileName == "" {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.FileAttachment("./uploads/"+fileName, fileName)
	}
}

func HandleUpdate(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) { //Delete the original file first
		ID := ctx.Param("ID")
		file := "./uploads/" + utilities.FindFileFromDirectory(ID)
		log.Println(file)
		err := os.Remove(file)
		if err != nil {
			if os.IsNotExist(err) {
				ctx.Status(http.StatusNotFound)
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		newFile, err := ctx.FormFile("file")
		fileEXT := strings.Split(newFile.Filename, ".")

		if len(fileEXT) <= 1 {
			ctx.SaveUploadedFile(newFile, "./uploads/"+ID)
			ctx.Status(http.StatusOK)
			return
		}

		ctx.SaveUploadedFile(newFile, "./uploads/"+ID+"."+fileEXT[len(fileEXT)-1])

		ctx.Status(http.StatusOK)
	}

}

func HandleDelete(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Param("ID")
		file := "./uploads/" + utilities.FindFileFromDirectory(ID)
		err := os.Remove(file)
		if err != nil {
			if os.IsNotExist(err) {
				ctx.Status(http.StatusNotFound)
				return
			}
			ctx.Status(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		ctx.Status(http.StatusOK)
	}
}
