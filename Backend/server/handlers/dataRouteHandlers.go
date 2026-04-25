package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AyushGlitchedOut/Docummunity/authUtils"
	"github.com/AyushGlitchedOut/Docummunity/consts"
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

func HostDataFiles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileName := ctx.Param("filename")
		if fileName == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No File Found for the Record",
			})
			return
		}

		verifiedFileName := filepath.Base(fileName)

		filePath := filepath.Join(consts.FileDirectory, verifiedFileName)

		ctx.File(filePath)
	}
}

func HostDataPreview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileName := ctx.Param("filename")
		if fileName == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Preview Found for the Record",
			})
			return
		}

		verifiedFileName := filepath.Base(fileName)

		filePath := filepath.Join(consts.PreviewImgDirectory, verifiedFileName)

		ctx.File(filePath)

	}
}

func HandleDataGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := &dbUtils.DATA{}

		//Get Data UUID
		uuid := ctx.Param("uuid")
		if uuid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No uuid found for the record",
			})
			return
		}

		results, err := dbUtils.GetRecord(ctx, uuid, db)
		if err != nil {
			if strings.Contains(err.Error(), "No Rows Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "The record doesn't exist",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		data = results

		ctx.JSON(http.StatusOK, gin.H{
			"message": data,
		})

	}
}

func HandleDataCREATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dataRecord := &dbUtils.DATA{}

		//See if file is there or not
		document, err := ctx.FormFile("FILE")
		if err != nil {
			if strings.Contains(err.Error(), "request body too large") {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "File too Large",
				})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "File Not Found!",
			})
			return
		}

		//Time-UUID
		uuid := utilities.GenerateUUID()
		dataRecord.UUID = uuid

		//Obtain Creator UID From JWT
		CreatorUID, err := authUtils.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		dataRecord.CREATOR_ID = CreatorUID

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
			if previewIMG.Size > consts.MaxPictureSize {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "Preview Picture should be less than " + strconv.Itoa(consts.MaxPictureSize>>20) + "mb!",
				})
				return
			}
			previewIMGPath = consts.PreviewImgDirectory + uuid + filepath.Ext(previewIMG.Filename)
			err = ctx.SaveUploadedFile(previewIMG, previewIMGPath)
			if err == nil {
				dataRecord.PREVIEW_IMG_PATH = previewIMGPath
			}
		}

		//Save File
		documentPath := consts.FileDirectory + uuid + filepath.Ext(document.Filename)
		if document.Size > consts.MaxDocumentSize {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "File must be lesser than " + strconv.Itoa(consts.MaxDocumentSize>>20) + "mb!",
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
	return func(ctx *gin.Context) {

		//Get Data UUID
		uuid := ctx.Param("uuid")
		if uuid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No uuid found for the record",
			})
			return
		}

		//Obtain Creator UID From JWT
		CreatorUID, err := authUtils.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = dbUtils.DeleteRecord(ctx, uuid, CreatorUID, db)

		if err != nil {
			if strings.Contains(err.Error(), "No Rows Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No Record Found to delete!",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Record Successfully Deleted",
		})
	}
}

func HandleDataSEARCH(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
