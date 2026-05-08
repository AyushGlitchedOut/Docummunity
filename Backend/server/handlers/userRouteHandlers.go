package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AyushGlitchedOut/Docummunity/authUtils"
	"github.com/AyushGlitchedOut/Docummunity/consts"
	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/gin-gonic/gin"
)

// User Functions

func HostUserPROFILE_PIC() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileName := ctx.Param("filename")
		if fileName == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Profile Picture Found for the User",
			})
			return
		}

		verifiedFileName := filepath.Base(fileName)

		filePath := filepath.Join(consts.ProfilePicDirectory, verifiedFileName)

		//TODO: Build a caching system and put a cache-control header
		ctx.File(filePath)
	}
}

func HandleUserGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		publicUserInfo := &dbUtils.USER_PUBLIC{}

		uid := ctx.Param("uid")

		if uid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UID found for User",
			})
			return
		}
		if uid == "000" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid UID found",
			})
			return
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

		//Obtain UID From JWT
		UID, err := authUtils.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		result, err := dbUtils.GetUserAccount(ctx, UID, db)
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
		if uid == "000" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid UID found",
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

		//Check if profile picture exists or not
		pictureFile, profilePicErr := ctx.FormFile("PROFILE_PIC")
		if profilePicErr != nil {
			if strings.Contains(profilePicErr.Error(), "request body too large") {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "File too Large",
				})
				return
			}
		}

		//Obtain UID From JWT
		UID, err := authUtils.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		user.UID = UID

		//Profile Pic Saving logic
		if profilePicErr == nil && pictureFile != nil {
			if pictureFile.Size > consts.MaxPictureSize {
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "Profile Picture should be less than " + strconv.Itoa(consts.MaxPictureSize>>20) + "mb!",
				})
				return
			}

			path := consts.ProfilePicDirectory + user.UID + filepath.Ext(pictureFile.Filename)
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

		updatedUser := &dbUtils.UserInfoUpdate{}

		//Obtain UID From JWT
		UID, err := authUtils.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//Obtain Details
		updatedUser.DISPLAY_NAME = ctx.PostForm("NAME")
		updatedUser.BIO = ctx.PostForm("BIO")
		updatedUser.SETTINGS = ctx.PostForm("SETTINGS")
		emptyProfilePicture := ctx.PostForm("emptyProfilePicture")

		//verify display name
		if updatedUser.DISPLAY_NAME == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Display Name Not Provided",
			})
			return
		}

		//verify emptyProfilePicture argument
		if emptyProfilePicture == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Please Provide the emptyProfilePicture argument",
			})
			return
		}

		//Get User's Old Info
		oldUserInfo, err := dbUtils.GetUserInfo(ctx, UID, db)
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

		//Save Profile Picture

		deletedOldUserProfilePicturePath := ""
		updatedUserProfilePicture := ""
		profilePicture, err := ctx.FormFile("PROFILE_PIC")
		if emptyProfilePicture == "true" {
			updatedUser.PROFILE_PIC = ""
			if oldUserInfo.PROFILE_PIC != "" {
				oldUserProfilePicturePathFilename := filepath.Base(oldUserInfo.PROFILE_PIC)
				oldUserProfilePicturePathLocation := filepath.Dir(oldUserInfo.PROFILE_PIC)

				deletedOldUserProfilePicturePath = filepath.Join(oldUserProfilePicturePathLocation, "__DELETED__"+oldUserProfilePicturePathFilename)

				err = os.Rename(oldUserInfo.PROFILE_PIC, deletedOldUserProfilePicturePath)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

			}
		}
		if emptyProfilePicture != "true" {
			if err != nil {
				updatedUser.PROFILE_PIC = oldUserInfo.PROFILE_PIC
				if strings.Contains(err.Error(), "request body too large") {
					ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "File Too Large",
					})
					return
				}
			} else {
				if profilePicture.Size > consts.MaxPictureSize {
					ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "Profile Picture should be less than " + strconv.Itoa(consts.MaxPictureSize>>20) + "mb!",
					})
					return
				}

				//mark old file for deletion

				if oldUserInfo.PROFILE_PIC != "" {
					oldUserProfilePicturePathFilename := filepath.Base(oldUserInfo.PROFILE_PIC)
					oldUserProfilePicturePathLocation := filepath.Dir(oldUserInfo.PROFILE_PIC)

					deletedOldUserProfilePicturePath = filepath.Join(oldUserProfilePicturePathLocation, "__DELETED__"+oldUserProfilePicturePathFilename)

					err = os.Rename(oldUserInfo.PROFILE_PIC, deletedOldUserProfilePicturePath)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}

				}

				updatedUserProfilePicture = consts.ProfilePicDirectory + UID + filepath.Ext(profilePicture.Filename)
				err = ctx.SaveUploadedFile(profilePicture, updatedUserProfilePicture)
				if err != nil {
					if oldUserInfo.PROFILE_PIC != "" {
						FSerr := os.Rename(deletedOldUserProfilePicturePath, oldUserInfo.PROFILE_PIC)
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
				updatedUser.PROFILE_PIC = updatedUserProfilePicture
			}
		}

		err = dbUtils.UpdateUserInfo(ctx, UID, updatedUser, db)
		if err != nil {

			var FSerr error
			if updatedUserProfilePicture != "" {
				FSerr = os.Remove(updatedUserProfilePicture)
				if FSerr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": FSerr.Error(),
					})
					return
				}
			}
			if deletedOldUserProfilePicturePath != "" {
				FSerr = os.Rename(deletedOldUserProfilePicturePath, oldUserInfo.PROFILE_PIC)
				if FSerr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": FSerr.Error(),
					})
					return
				}
			}

			if strings.Contains(err.Error(), "No User found") {
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

		//delete the marked file
		if deletedOldUserProfilePicturePath != "" {
			err = os.Remove(deletedOldUserProfilePicturePath)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "Updated User",
					"warning": "User Updated, but old Profile Picture file still exists. Request manual deletion",
				})
				return
			}

		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "User Updated",
		})

	}
}

func HandleUserDELETE(db *sql.DB, keeprecords bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Obtain UID From JWT
		UID, err := authUtils.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = dbUtils.DeleteUser(ctx, UID, db, keeprecords)
		if err != nil {
			if strings.Contains(err.Error(), "No User found with UID:") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No User Found to Delete",
				})
				return
			}
			if strings.Contains(err.Error(), "No Rows Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No User Found to Delete",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "User Successfully Deleted, Records Preserved: " + strconv.FormatBool(keeprecords),
		})
	}
}

func HandleUserSEARCH(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get Query
		query := ctx.Param("query")
		if query == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Query Found",
			})
			return
		}
	}
}
