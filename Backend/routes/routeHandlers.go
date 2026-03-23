package routes

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
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
			ctx.Status(http.StatusBadRequest)
			return
		}

		path := "./uploads/" + file.Filename

		err = ctx.SaveUploadedFile(file, path)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		data = dbUtils.TestData{
			TIMEID:   int(time.Now().Unix()),
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

		results, err := dbUtils.ReadData(ID, db)

		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusNotFound)
			return
		}

		fileName := strings.Split(results[0].FILEPATH, "/")

		ctx.FileAttachment(results[0].FILEPATH, fileName[len(fileName)-1])
	}
}

func HandleUpdate(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Param("ID")
		newFile, err := ctx.FormFile("file")
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusBadRequest)
			return
		}

		oldData, err := dbUtils.ReadData(ID, db)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		err = dbUtils.DeleteData(ID, db)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		err = ctx.SaveUploadedFile(newFile, oldData[0].FILEPATH)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		newData := dbUtils.TestData{TIMEID: int(time.Now().Unix()), NAME: oldData[0].NAME, FILEPATH: oldData[0].FILEPATH}

		err = dbUtils.InsertData(newData, db)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}

}

func HandleDelete(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Param("ID")
		err := dbUtils.DeleteData(ID, db)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.Status(http.StatusOK)
	}
}
