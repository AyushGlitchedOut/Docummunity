package handlers

import (
	"database/sql"
	"log"
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

// DEBUG: Test functions
func VerifyTest(ctx *gin.Context) {
	token, _ := ctx.Get("tokenUID")
	log.Println("Token: ", token)
}

// Test If server is active
func HandlePING(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Server Active!",
	})
}

// Data Functions

// Handler to statically Host Document Files for Records
func HostDataFiles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get the filename from query
		fileName := ctx.Param("filename")

		if fileName == "" {
			//400, If the file isnt found
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No File Found for the Record",
			})
			return
		}

		//Construct a filepath to Serve
		verifiedFileName := filepath.Base(fileName)
		filePath := filepath.Join(consts.FileDirectory, verifiedFileName)

		// Serve the File
		ctx.File(filePath)
	}
}

// Handler to statically Host Preview Files for Records
func HostDataPreview() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Get Filename from Query
		fileName := ctx.Param("filename")

		if fileName == "" {
			//400, If Filename isnt there
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Preview Found for the Record",
			})
			return
		}

		//Construct Filename to Server
		verifiedFileName := filepath.Base(fileName)
		filePath := filepath.Join(consts.PreviewImgDirectory, verifiedFileName)

		// Serve the file
		ctx.File(filePath)

	}
}

// Handler for when the User wants to fetch a Record
func HandleDataGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := &dbUtils.DATA{}

		//Get Data UUID
		uuid := ctx.Param("uuid")
		if uuid == "" {
			//400, if No UUID
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UUID found for the record",
			})
			return
		}

		//Get the Record
		data, err := dbUtils.GetRecord(ctx, uuid, db)
		if err != nil {
			//400, If No record found
			if strings.Contains(err.Error(), "No Records Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No Records Found",
				})
				return
			}
			//500, For any other DB error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			log.Println("ERROR:", err.Error())
			return
		}

		//200, If All goes Well
		ctx.JSON(http.StatusOK, gin.H{
			"message": data,
		})

	}
}

// Handler to Create A Record
func HandleDataCREATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dataRecord := &dbUtils.DATA{}

		//See if file is there IN the request or Not
		document, err := ctx.FormFile("FILE")
		if err != nil {
			//413, IF MaxSizeMiddleware returns an error
			if strings.Contains(err.Error(), "request body too large") {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "File too Large",
				})
				return
			}
			//400, For any other error
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "File Not Found!",
			})
			return
		}

		//Generate Time-UUID
		uuid := utilities.GenerateUUID()
		dataRecord.UUID = uuid

		//Obtain Creator UID From JWT
		dataRecord.CREATOR_ID, err = authUtils.ParseToken(ctx)
		if err != nil {
			//400, If UID not
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JWT Token",
			})
			return
		}

		//Name and Description
		dataRecord.NAME = ctx.PostForm("NAME")
		dataRecord.DESCRIPTION = ctx.PostForm("DESCRIPTION")
		if dataRecord.NAME == "" {
			//400, IF name not present
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Name Provided",
			})
			return
		}

		//Save Preview Img
		previewIMGPath := ""
		previewIMG, err := ctx.FormFile("PREVIEW")

		if err == nil {
			//413, If Image is bigger than Maximum Picture Size
			if previewIMG.Size > consts.MaxPictureSize {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "Preview Picture should be less than " + strconv.Itoa(consts.MaxPictureSize>>20) + "mb!",
				})
				return
			}

			//415, If the Preview Image has an unsupported Image Type
			fileTypeValid, err := filetypeDetector(previewIMG, consts.AllowedImageExtensions)
			if !fileTypeValid || err != nil {
				ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Preview should be a PNG, JPEG or WEBP file",
				})
				return
			}

			//Construct Filepath
			previewIMGPath = filepath.Join(consts.PreviewImgDirectory, uuid+filepath.Ext(previewIMG.Filename))

			//Save the File
			err = ctx.SaveUploadedFile(previewIMG, previewIMGPath)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error Saving Profile Picture",
				})
				return
			}
			dataRecord.PREVIEW_IMG_PATH = previewIMGPath

		}

		//Check if the Document Uploaded matches allowed filetypes
		fileTypeValid, err := filetypeDetector(document, consts.AllowedDocumentExtensions)

		if !fileTypeValid || err != nil {

			//Remove the saved Preview Image if something goes wrong
			if previewIMGPath != "" {
				FSerr := os.Remove(previewIMGPath)
				if FSerr != nil {
					//500, If that file deletion too fails
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Many Errors Together",
					})
					log.Println("ERROR:", FSerr.Error())
					return
				}
			}

			//415, If the document is not of Supported Filetypes
			ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": "Only PDF files are allowed",
			})
			return
		}

		if document.Size > consts.MaxDocumentSize {
			//Delete the previously saved Preview Image file upon error
			if previewIMGPath != "" {
				FSerr := os.Remove(previewIMGPath)
				if FSerr != nil {
					//500, IF file deletion too fails
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Many Errors Together",
					})
					log.Println("ERROR:", FSerr.Error())
					return
				}
			}
			//413, If Document size exceeds Maximum Document Size
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "File must be lesser than " + strconv.Itoa(consts.MaxDocumentSize>>20) + "mb!",
			})
			return
		}

		//Construct Filepath
		documentPath := filepath.Join(consts.FileDirectory, uuid+filepath.Ext(document.Filename))

		err = ctx.SaveUploadedFile(document, documentPath)
		if err != nil {
			//Delete the previously saved Preview Image if Document saving fails
			if previewIMGPath != "" {
				FSerr := os.Remove(previewIMGPath)
				if FSerr != nil {
					//500, IF file deletion fails too
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Many Errors Together",
					})
					log.Println("ERROR:", FSerr.Error())
					return
				}
			}
			//500, If Any errors
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			log.Println("ERROR:", err.Error())
			return
		}
		dataRecord.FILEPATH = documentPath

		//Create The Record
		err = dbUtils.CreateRecord(ctx, dataRecord, db)
		if err != nil {
			//If Any error during Creating the Record, Delete the Files to preserve consistency
			FSerr := os.Remove(documentPath)
			FSerr = os.Remove(previewIMGPath)
			if FSerr != nil {
				//500, If any FS error
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Many Errors Together!",
				})
				log.Println("ERROR:", FSerr.Error())
				return
			}

			//404, If The JWT given gives UID of a User that doesnt exist i.e. No foreign Key possible
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "The Supposed Creator Doesnt exist",
				})
				return
			}

			//500, If any Other errors
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			log.Println("ERROR:", err.Error())
			return
		}

		//200, If All goes well
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Created Record",
		})

	}
}

// DEV NOTE: I wrote this function initially without the deletion of the old file.
// I wrote the old-file deletion code in the commit after 4 May, and wrote the whole thing myself (see the changes between commits of 4 May and after that.).
// This is one of the more complex logics I have designed myself, and am proud that I wrote the entire thing without using any AI.

// Handler to Update the Record
func HandleDataUPDATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Update Info store
		updatedRecordInfo := &dbUtils.DataInfoUpdate{}

		//Get Data UUID
		uuid := ctx.Param("uuid")
		if uuid == "" {
			//400, If no UUID for data
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UUID found for the record",
			})
			return
		}

		//Get User's UID
		CreatorUID, err := authUtils.ParseToken(ctx)
		if err != nil {
			//400, If no UID for User
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JWT Token",
			})
			return
		}

		//Get Details
		updatedRecordInfo.NAME = ctx.PostForm("NAME")
		updatedRecordInfo.DESCRIPTION = ctx.PostForm("DESCRIPTION")

		//Get emptyPreview Parameter
		emptyPreview, err := strconv.ParseBool(ctx.DefaultQuery("emptyPreview", "false"))
		if err != nil {
			emptyPreview = false
		}

		if updatedRecordInfo.NAME == "" {
			//400, If No Name given
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Name Provided",
			})
			return
		}

		//Get Old Record Info
		oldRecordInfo, err := dbUtils.GetRecord(ctx, uuid, db)
		if err != nil {
			if strings.Contains(err.Error(), "No Records Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "The record doesn't exist",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			log.Println("ERROR:", err.Error())
			return
		}

		//Get Preview-Image
		deletedOldPreviewIMGPath := ""
		newPreviewIMGPath := ""
		newPreviewIMG, err := ctx.FormFile("PREVIEW")

		//If the User want an Empty Preview Field
		if emptyPreview {

			//Make the new Preview empty
			updatedRecordInfo.PREVIEW_IMG_PATH = ""

			//Mark the old Preview File for deletion
			if oldRecordInfo.PREVIEW_IMG_PATH != "" {

				//Construct the filepath
				deletedOldPreviewIMGPath = deletedFilePathMaker(oldRecordInfo.PREVIEW_IMG_PATH)

				//Delete the File
				err = os.Rename(oldRecordInfo.PREVIEW_IMG_PATH, deletedOldPreviewIMGPath)
				if err != nil {
					//500, For any FS error
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Something Went Wrong!",
					})
					log.Println("ERROR:", err.Error())
					return
				}
			}
		}

		//If the User wants a Preview
		if !emptyPreview {
			if err != nil {
				//In case of error, make the new Preview the same as old one
				updatedRecordInfo.PREVIEW_IMG_PATH = oldRecordInfo.PREVIEW_IMG_PATH

				//413, If the file is too large
				if strings.Contains(err.Error(), "request body too large") {
					ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "File too Large",
					})
					return
				}
			} else {

				//415, If the Preview Image has an unsupported Image Type
				filetypeValid, err := filetypeDetector(newPreviewIMG, consts.AllowedImageExtensions)
				if !filetypeValid || err != nil {
					ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
						"error": "Preview should be a PNG, JPEG or WEBP file",
					})
					return
				}

				//413, Of the file exceeds document size Limits
				if newPreviewIMG.Size > consts.MaxPictureSize {
					ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "Preview Picture should be less than " + strconv.Itoa(consts.MaxPictureSize>>20) + "mb!",
					})
					return
				}

				//Mark the old image as deleted
				if oldRecordInfo.PREVIEW_IMG_PATH != "" {

					//Construct the File Path
					deletedOldPreviewIMGPath = deletedFilePathMaker(oldRecordInfo.PREVIEW_IMG_PATH)

					//Delete the File
					err = os.Rename(oldRecordInfo.PREVIEW_IMG_PATH, deletedOldPreviewIMGPath)
					if err != nil {
						//500, In case of FS error
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"error": "Something Went Wrong!",
						})
						log.Println("ERROR:", err.Error())
						return
					}
				}

				//Construct the new Profile Picture
				newPreviewIMGPath = filepath.Join(consts.PreviewImgDirectory, uuid+filepath.Ext(newPreviewIMG.Filename))
				//Save the new File
				err = ctx.SaveUploadedFile(newPreviewIMG, newPreviewIMGPath)
				if err != nil {
					//In case of error, Un-mark the Old file for deletion
					if deletedOldPreviewIMGPath != "" {
						FSerr := os.Rename(deletedOldPreviewIMGPath, oldRecordInfo.PREVIEW_IMG_PATH)
						if FSerr != nil {
							//500, For any FS error
							ctx.JSON(http.StatusInternalServerError, gin.H{
								"error": "Many Errors Together",
							})
							log.Println("ERROR:", FSerr.Error())
							return
						}
					}

					//500, For any Errors
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Something Went Wrong!",
					})
					log.Println("ERROR:", err.Error())
					return
				}
				updatedRecordInfo.PREVIEW_IMG_PATH = newPreviewIMGPath
			}
		}

		//Update the actual record
		err = dbUtils.UpdateRecord(ctx, uuid, updatedRecordInfo, CreatorUID, db)
		if err != nil {

			//If DB operation fails, we have to delete the New file and unmark the old file for deletion
			var FSerr error

			//Delete the new File saved
			if newPreviewIMGPath != "" {
				FSerr = os.Remove(newPreviewIMGPath)
				if FSerr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Many Errors Together",
					})
					log.Println("ERROR:", FSerr.Error())
					return
				}
			}

			//Un-mark the Old File for deletion by renaming it to its original name
			if deletedOldPreviewIMGPath != "" {
				FSerr = os.Rename(deletedOldPreviewIMGPath, oldRecordInfo.PREVIEW_IMG_PATH)
				if FSerr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Many Errors Together",
					})
					log.Println("ERROR:", err.Error())
					return
				}
			}

			//404, IF No record found
			if strings.Contains(err.Error(), "No Records Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No Record Found for the UUID",
				})
				return
			}
			//500, For ANy other DB error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			log.Println("ERROR:", err.Error())
			return
		}

		//Delete the file finally
		if deletedOldPreviewIMGPath != "" {
			err = os.Remove(deletedOldPreviewIMGPath)
			if err != nil {
				//200, if file deletion fails, we just tell the user to request a manual deletion of the file since its references have been aleady deleted
				ctx.JSON(http.StatusOK, gin.H{
					"message": "Updated Record Successfully",
					"warning": "Record Updated, but old Preview file still exists. Request manual deletion",
				})
				return
			}

		}

		//200, If All goes well
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
			//400, If No UUID found
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UUID found for the record",
			})
			return
		}

		//Obtain Creator UID From JWT
		CreatorUID, err := authUtils.ParseToken(ctx)
		if err != nil {
			//400, If CreatorUID not found
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JWT Token",
			})
			return
		}

		//Delete the Record from the DB
		err = dbUtils.DeleteRecord(ctx, uuid, CreatorUID, db)

		if err != nil {
			//404, If Record not Found
			if strings.Contains(err.Error(), "No Records Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No Record Found to delete!",
				})
				return
			}
			//500, for any other errors
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			log.Println("ERROR:", err.Error())
			return
		}

		//200, If All goes well
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Record Successfully Deleted",
		})
	}
}

// Handler for When the User wants to search a Record
func HandleDataSEARCH(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		searchResults := []*dbUtils.DATA{}

		//Get Query
		query := ctx.Param("query")
		if query == "" {
			//400, If query not Found
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
		searchResults, err = dbUtils.SearchRecord(ctx, strings.Split(query, " "), db, useDescription)

		//500 If any DB error
		if err != nil {
			//400, IF No record is found
			if strings.Contains(err.Error(), "No Records Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No Records Found",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong",
			})
			log.Println("ERROR:", err.Error())
			return
		}

		//200, Send the results If All goes well
		ctx.JSON(http.StatusOK, gin.H{
			"message": searchResults,
		})

	}
}
