package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/utilities"
	"github.com/gin-gonic/gin"
)

// User Functions

func HandleUserGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		publicUserInfo := &dbUtils.USER_PUBLIC{}

		uid := ctx.Param("uid")

		if uid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UID found for User",
			})
		}

		result, err := dbUtils.GetUserInfo(ctx, uid, db)
		if err != nil {
			if strings.Contains(err.Error(), "No Rows Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No User Found",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			return
		}
		publicUserInfo = result

		ctx.JSON(http.StatusOK, gin.H{
			"message": publicUserInfo,
		})
	}
}

func HandleUserACCOUNT(db *sql.DB) gin.HandlerFunc {
	//For when the User requires their own information (including sensitive)
	return func(ctx *gin.Context) {

		userAccountInfo := &dbUtils.USER{}

		token, err := utilities.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		result, err := dbUtils.GetUserAccount(ctx, token, db)
		if err != nil {
			if strings.Contains(err.Error(), "No Rows Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No User Found",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			return
		}
		userAccountInfo = result

		//Parse settings
		rawSettings := result.SETTINGS
		var settings map[string]any
		err = json.Unmarshal([]byte(rawSettings), &settings)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"error":   "Error parsing Settings, sending Unparsed",
				"message": userAccountInfo,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": map[string]any{
				"UID":           userAccountInfo.UID,
				"DISPLAY_NAME":  userAccountInfo.DISPLAY_NAME,
				"BIO":           userAccountInfo.BIO,
				"PROFILE_PIC":   userAccountInfo.PROFILE_PIC,
				"CREATION_DATE": userAccountInfo.CREATION_DATE,
				"SETTINGS":      settings,
			},
		})

	}
}

func HandleUserRecordsGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userRecords []*dbUtils.DATA

		uid := ctx.Param("uid")
		if uid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UID found for User",
			})
			return
		}

		results, err := dbUtils.GetUserRecords(ctx, uid, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		userRecords = results

		if userRecords == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No records Found for that User!",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": userRecords,
		})

	}
}

func HandleUserCREATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := &dbUtils.USER{}

		//Obtain UID From JWT
		UID, err := utilities.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		user.UID = UID

		//Profile Pic Saving logic
		pictureFile, err := ctx.FormFile("PROFILE_PIC")
		if err == nil && pictureFile != nil {
			if pictureFile.Size > maxPictureSize {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "Profile Picture should be less than " + strconv.Itoa(maxPictureSize>>20) + "mb!",
				})
				return
			}

			path := utilities.ProfilePicDirectory + user.UID + filepath.Ext(pictureFile.Filename)
			err = ctx.SaveUploadedFile(pictureFile, path)
			if err == nil {
				user.PROFILE_PIC = path
			}
		}

		//Save fields from FormData
		user.DISPLAY_NAME = ctx.PostForm("DISPLAY_NAME")
		user.BIO = ctx.PostForm("BIO")
		user.CREATION_DATE = time.Now().UTC().Format("02/01/2006 03:04:05 PM MST")
		user.SETTINGS = ctx.PostForm("SETTINGS")

		if user.DISPLAY_NAME == "" {
			user.DISPLAY_NAME = "[Unknown]"
		}

		//Insert User
		err = dbUtils.CreateUser(ctx, user, db)
		if err != nil {
			//If User is repeating Sign-In / Sign-Up, it would try to create the User's record again
			if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "User Not Created Since it already exists",
				})
				return
			}
			//DB errors
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		//All Goes well
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "User Created",
		})
	}
}

func HandleUserUPDATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func HandleUserDELETE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func HandleUserDELETEWithRecords(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func HandleUserSEARCH(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
