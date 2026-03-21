package routes

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/AyushGlitchedOut/Docummunity/utilities"
	"github.com/gin-gonic/gin"
)

func HandleGET(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OHHHH YEAHHHHHHHHH",
	})
}

func HandleCREATE(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ID := ctx.PostForm("ID")
	fileEXT := strings.Split(file.Filename, ".")

	if len(fileEXT) <= 1 {
		ctx.SaveUploadedFile(file, "./uploads/"+ID)
		ctx.Status(http.StatusOK)
		return
	}

	ctx.SaveUploadedFile(file, "./uploads/"+ID+"."+fileEXT[len(fileEXT)-1])

	ctx.Status(http.StatusOK)
}

func HandleRead(ctx *gin.Context) {
	ID := ctx.Param("ID")
	fileName := utilities.FindFileFromDirectory(ID)
	if fileName == "" {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.FileAttachment("./uploads/"+fileName, fileName)
}

func HandleUpdate(ctx *gin.Context) {

	//Delete the original file first
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

func HandleDelete(ctx *gin.Context) {
	ID := ctx.Param("ID")
	file := "./uploads/" + utilities.FindFileFromDirectory(ID)
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
}
