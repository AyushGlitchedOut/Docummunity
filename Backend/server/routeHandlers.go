package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/utilities"
	"github.com/gin-gonic/gin"
)

const (
	maxProfilePictureSize = 2 << 20 //2mb
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

// User Functions
func HandleUserGET(db *sql.DB) gin.HandlerFunc {
	//When fetching information about another user
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Fetch User",
		})
	}
}
func HandleUserACCOUNT(db *sql.DB) gin.HandlerFunc {
	//For when the User requires their own information (including sensitive)
	return func(ctx *gin.Context) {}
}

func HandleUserCREATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := &dbUtils.USER{}

		//Obtain UID From JWT
		UID, exists := ctx.Get("tokenUID")
		if !exists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Error getting UID",
			})
			return
		}
		UIDstr, ok := UID.(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Error getting UID",
			})
			return
		}
		user.UID = UIDstr
		//Validation
		if user.UID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UID Provided",
			})
			return
		}

		//Profile Pic Saving logic
		pictureFile, err := ctx.FormFile("PROFILE_PIC")
		if err == nil && pictureFile != nil {
			if pictureFile.Size > maxProfilePictureSize {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Profile Picture should be less than 2mb",
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

// Data Functions
func HandleDataGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Fetch Data",
		})
	}
}

func HandleDataCREATE(db *sql.DB) gin.HandlerFunc {
	//NOTE: Get Creator_ID from JWT, not the request itself, since otherwise anyone can create records on anyone's behalf
	return func(ctx *gin.Context) {}
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
