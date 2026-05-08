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

// DEV NOTE: I wrote this function initially without the deletion of the old file. I wrote the old-file deletion code in the commit after 4 May, and wrote the whole thing myself (see the changes between commits of 4 May and after that.). This is one of the more complex logics I have designed myself, and am proud that I wrote the entire thing without using any AI.
func HandleDataUPDATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Update Info store
		updatedRecordInfo := &dbUtils.DataInfoUpdate{}

		//Get Data UUID
		uuid := ctx.Param("uuid")
		if uuid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No uuid found for the record",
			})
			return
		}

		//Get User's UID
		CreatorUID, err := authUtils.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//Get Details
		updatedRecordInfo.NAME = ctx.PostForm("NAME")
		updatedRecordInfo.DESCRIPTION = ctx.PostForm("DESCRIPTION")
		emptyPreview := ctx.PostForm("emptyPreview")

		if updatedRecordInfo.NAME == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Name Provided",
			})
			return
		}
		if emptyPreview == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Please provide the emptyPreview argument",
			})
			return
		}

		//Get Old Record Info
		oldRecordInfo, err := dbUtils.GetRecord(ctx, uuid, db)
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

		//Get Preview-Image
		deletedOldPreviewIMGPath := ""
		newPreviewIMGPath := ""
		newPreviewIMG, err := ctx.FormFile("PREVIEW")
		if emptyPreview == "true" {
			updatedRecordInfo.PREVIEW_IMG_PATH = ""
			if oldRecordInfo.PREVIEW_IMG_PATH != "" {

				deletedOldPreviewIMGPathFileName := filepath.Base(oldRecordInfo.PREVIEW_IMG_PATH)
				deletedOldPreviewIMGPathLocation := filepath.Dir(oldRecordInfo.PREVIEW_IMG_PATH)
				deletedOldPreviewIMGPath = filepath.Join(deletedOldPreviewIMGPathLocation, "__DELETED__"+deletedOldPreviewIMGPathFileName)
				err = os.Rename(oldRecordInfo.PREVIEW_IMG_PATH, deletedOldPreviewIMGPath)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
			}
		}
		if emptyPreview != "true" {
			if err != nil {
				updatedRecordInfo.PREVIEW_IMG_PATH = oldRecordInfo.PREVIEW_IMG_PATH
				if strings.Contains(err.Error(), "request body too large") {
					ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "File too Large",
					})
					return
				}
			} else {
				if newPreviewIMG.Size > consts.MaxPictureSize {
					ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "Preview Picture should be less than " + strconv.Itoa(consts.MaxPictureSize>>20) + "mb!",
					})
					return
				}

				//Mark the old image as deleted
				if oldRecordInfo.PREVIEW_IMG_PATH != "" {

					deletedOldPreviewIMGPathFileName := filepath.Base(oldRecordInfo.PREVIEW_IMG_PATH)
					deletedOldPreviewIMGPathLocation := filepath.Dir(oldRecordInfo.PREVIEW_IMG_PATH)
					deletedOldPreviewIMGPath = filepath.Join(deletedOldPreviewIMGPathLocation, "__DELETED__"+deletedOldPreviewIMGPathFileName)

					err = os.Rename(oldRecordInfo.PREVIEW_IMG_PATH, deletedOldPreviewIMGPath)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}
				}

				newPreviewIMGPath = consts.PreviewImgDirectory + uuid + filepath.Ext(newPreviewIMG.Filename)
				err = ctx.SaveUploadedFile(newPreviewIMG, newPreviewIMGPath)
				if err != nil {
					if deletedOldPreviewIMGPath != "" {
						FSerr := os.Rename(deletedOldPreviewIMGPath, oldRecordInfo.PREVIEW_IMG_PATH)
						if FSerr != nil {
							ctx.JSON(http.StatusInternalServerError, gin.H{
								"error": FSerr.Error(),
							})
							return
						}
					}
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
				updatedRecordInfo.PREVIEW_IMG_PATH = newPreviewIMGPath
			}
		}

		//Update the actual record
		err = dbUtils.UpdateRecord(ctx, uuid, updatedRecordInfo, CreatorUID, db)
		if err != nil {
			var FSerr error
			if newPreviewIMGPath != "" {
				FSerr = os.Remove(newPreviewIMGPath)
				if FSerr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": FSerr.Error(),
					})
					return
				}
			}
			if deletedOldPreviewIMGPath != "" {
				FSerr = os.Rename(deletedOldPreviewIMGPath, oldRecordInfo.PREVIEW_IMG_PATH)
				if FSerr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": FSerr.Error(),
					})
					return
				}
			}
			if strings.Contains(err.Error(), "No Record found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No Record Found for the UUID",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		//delete the file finally
		if deletedOldPreviewIMGPath != "" {
			err = os.Remove(deletedOldPreviewIMGPath)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "Updated Record Successfully",
					"warning": "Record Updated, but old Preview file still exists. Request manual deletion",
				})
				return
			}

		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Updated Record Successfully",
		})

	}
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
	return func(ctx *gin.Context) {
		searchResults := []*dbUtils.DATA{}

		//Get Query
		query := ctx.Param("query")
		if query == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Query Found",
			})
			return
		}

		//Get useDescription parameter
		useDescription, err := strconv.ParseBool(ctx.DefaultQuery("useDescription", "false"))
		if err != nil {
			useDescription = false
		}

		//Search From db
		results, err := dbUtils.SearchRecord(ctx, strings.Split(query, " "), db, useDescription)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if len(results) < 1 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No Records Found",
			})
			return
		}
		searchResults = results

		ctx.JSON(http.StatusOK, gin.H{
			"message": searchResults,
		})

	}
}
