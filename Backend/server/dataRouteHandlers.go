package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/utilities"
	"github.com/gin-gonic/gin"
)

func HandlePING(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Server Active!",
	})
}

// Test functions
func VerifyTest(ctx *gin.Context) {
	token, _ := ctx.Get("tokenUID")
	fmt.Println("Token: ", token)
}

// Data Functions
func HandleDataGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Fetch Data",
		})
	}
}

func HandleDataCREATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dataRecord := &dbUtils.DATA{}

		//Time-UUID
		uuid := utilities.GenerateUUID()
		dataRecord.UUID = uuid

		//CreatorID
		creatorID, err := utilities.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		dataRecord.CREATOR_ID = creatorID

		//Name and Description
		dataRecord.NAME = ctx.PostForm("NAME")
		dataRecord.DESCRIPTION = ctx.PostForm("DESCRIPTION")
		if dataRecord.NAME == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Name Provided",
			})
			return
		}

		//Save Preview Img
		previewIMGPath := ""
		previewIMG, err := ctx.FormFile("PREVIEW")
		if err == nil {
			if previewIMG.Size > maxPictureSize {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "Preview Picture should be less than " + strconv.Itoa(maxPictureSize>>20) + "mb!",
				})
				return
			}
			previewIMGPath = utilities.PreviewImgDirectory + uuid + filepath.Ext(previewIMG.Filename)
			err = ctx.SaveUploadedFile(previewIMG, previewIMGPath)
			if err == nil {
				dataRecord.PREVIEW_IMG_PATH = previewIMGPath
			}
		}

		//Save File
		document, err := ctx.FormFile("FILE")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "File Not Found!",
			})
			return
		}
		documentPath := utilities.FileDirectory + uuid + filepath.Ext(document.Filename)
		if document.Size > maxDocumentSize {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "File must be lesser than " + strconv.Itoa(maxDocumentSize>>20) + "mb!",
			})
			return
		}
		err = ctx.SaveUploadedFile(document, documentPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error!",
			})
			return
		}
		dataRecord.FILEPATH = documentPath

		//Create The Record
		err = dbUtils.CreateRecord(ctx, dataRecord, db)
		if err != nil {
			FSerr := os.Remove(documentPath)
			FSerr = os.Remove(previewIMGPath)
			if FSerr != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Many Internal Server Errors",
				})
				return
			}
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "The Supposed Creator Doesnt exist",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Created Record",
		})

	}
}

func HandleDataUPDATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func HandleDataDELETE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func HandleDataSEARCH(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
