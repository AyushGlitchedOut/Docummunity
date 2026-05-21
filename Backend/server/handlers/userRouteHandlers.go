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
//For protected Routes, we get UID from JWT
//For free routes, we get UID from parameters

// A Handler to host the routes for statically hosting User's Profile Pictures
func HostUserPROFILE_PIC() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//extract filename from parameters
		fileName := ctx.Param("filename")
		if fileName == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Profile Picture Found for the User",
			})
			return
		}

		//construct relative path to access in the profile picture directory
		verifiedFileName := filepath.Base(fileName)
		filePath := filepath.Join(consts.ProfilePicDirectory, verifiedFileName)

		//TODO: Build a caching system and put a cache-control header
		//Return the file
		ctx.File(filePath)
	}
}

// Handler for getting User's public info
func HandleUserGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		publicUserInfo := &dbUtils.USER_PUBLIC{}

		//Get user's UID
		uid := ctx.Param("uid")
		if uid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UID found for User",
			})
			return
		}
		//Make sure access to deleted user is prevented
		if uid == dbUtils.DeletedUserInfo.UID {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid UID found",
			})
			return
		}

		//Get Public Information about User
		publicUserInfo, err := dbUtils.GetUserInfo(ctx, uid, db)
		if err != nil {
			//404 For No user found
			if strings.Contains(err.Error(), "No Rows Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No User Found",
				})
				return
			}
			//500 For any other error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			return
		}

		//200, return the Info
		ctx.JSON(http.StatusOK, gin.H{
			"message": publicUserInfo,
		})
	}
}

// Handler For when the User requires their own information (including sensitive)
func HandleUserACCOUNT(db *sql.DB) gin.HandlerFunc {
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

		//Get User's Account Info
		userAccountInfo, err = dbUtils.GetUserAccount(ctx, UID, db)
		if err != nil {
			//404 If No user found
			if strings.Contains(err.Error(), "No Rows Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No User Found",
				})
				return
			}
			//500 If any other error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			return
		}

		//Parse settings from the obtained from the User Info from string to map[string]any (object)
		rawSettings := userAccountInfo.SETTINGS
		var settings map[string]any
		err = json.Unmarshal([]byte(rawSettings), &settings)
		if err != nil {
			//200: If any error occurs in parsing the settings i.e. Invalid JSON, we just send it without parsing
			ctx.JSON(http.StatusOK, gin.H{
				"error":   "Error parsing Settings, sending Unparsed",
				"message": userAccountInfo,
			})
			return
		}

		//200, If settings are converted into map[string]any (object), we send a JSON with the same structure as our *dbUtils.User struct, but replace settings from string to object
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

// Handler To get all Records related to a User
func HandleUserRecordsGET(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userRecords []*dbUtils.DATA

		//Get UID of the User
		uid := ctx.Param("uid")
		if uid == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No UID found for User",
			})
			return
		}
		//Prevent the User from accessing the reserved System User
		if uid == dbUtils.DeletedUserInfo.UID {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid UID found",
			})
			return
		}

		//Get the Records for a user
		userRecords, err := dbUtils.GetUserRecords(ctx, uid, db)
		if err != nil {
			//500 for any error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if userRecords == nil {
			//200, if No records found
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No records Found for that User!",
			})
			return
		}

		//200, Send the User's Records
		ctx.JSON(http.StatusOK, gin.H{
			"message": userRecords,
		})

	}
}

// Handler to Create a User
func HandleUserCREATE(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := &dbUtils.USER{}

		//Check if profile picture exists or not in the request
		pictureFile, profilePicErr := ctx.FormFile("PROFILE_PIC")
		if profilePicErr != nil {
			//413, if we get error from MaxSizeMiddleware while parsing the request
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

			//415, If the Profile Picture is an unsupported type
			if filepath.Ext(pictureFile.Filename) != ".png" && filepath.Ext(pictureFile.Filename) != ".jpg" {
				ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Profile Picturw should be of type .jpg or .png",
				})
				return
			}

			if pictureFile.Size > consts.MaxPictureSize {
				//413, if the image received is larger than Maximum Picture Size
				ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "Profile Picture should be less than " + strconv.Itoa(consts.MaxPictureSize>>20) + "mb!",
				})
				return
			}

			//Save the image
			path := filepath.Join(consts.ProfilePicDirectory, user.UID+filepath.Ext(pictureFile.Filename))
			err = ctx.SaveUploadedFile(pictureFile, path)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error Saving Profile Picture",
				})
				return
			}
			user.PROFILE_PIC = path
		}

		//Save fields from FormData
		user.DISPLAY_NAME = ctx.PostForm("DISPLAY_NAME")
		user.BIO = ctx.PostForm("BIO")
		user.CREATION_DATE = time.Now().UTC().Format("02/01/2006 03:04:05 PM MST")
		user.SETTINGS = ctx.PostForm("SETTINGS")

		if user.DISPLAY_NAME == "" {
			user.DISPLAY_NAME = "[Unknown]"
		}

		//Create the User
		err = dbUtils.CreateUser(ctx, user, db)
		if err != nil {
			//Delete the Profile Picture file if DB operation fails
			if user.PROFILE_PIC != "" {
				FSerr := os.Remove(user.PROFILE_PIC)
				if FSerr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Many Errors Together!",
					})
					return
				}
			}

			//409, If User is repeating Sign-In / Sign-Up, it would try to create the User's record again
			if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "User Not Created Since it already exists",
				})
				return
			}
			//500, Any other DB errors
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		//200, All Goes well
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "User Created",
		})
	}
}

// Handler for Updating User
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
		if updatedUser.DISPLAY_NAME == "" {
			//400, If No DisplayName
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Display Name Not Provided",
			})
			return
		}
		if emptyProfilePicture == "" {
			//400, If No emptyProfilePicture argument
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Please Provide the emptyProfilePicture argument",
			})
			return
		}

		//Get User's Old Info
		oldUserInfo, err := dbUtils.GetUserInfo(ctx, UID, db)
		if err != nil {
			//404, If the User doesnt exist
			if strings.Contains(err.Error(), "No Rows Found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No User Found",
				})
				return
			}
			//500, for any other server error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			return
		}

		//Save Profile Picture
		deletedOldUserProfilePicturePath := ""
		updatedUserProfilePicture := ""
		profilePicture, err := ctx.FormFile("PROFILE_PIC")

		//If the User doesnt wanna have a profile picture, we just simply mark the old file for deletion
		if emptyProfilePicture == "true" {
			//Make the new Profile Picture field empty
			updatedUser.PROFILE_PIC = ""

			//Mark the Old Profile Picture file for deletion
			if oldUserInfo.PROFILE_PIC != "" {

				//Construct Filepath for old File
				oldUserProfilePicturePathFilename := filepath.Base(oldUserInfo.PROFILE_PIC)
				oldUserProfilePicturePathLocation := filepath.Dir(oldUserInfo.PROFILE_PIC)
				deletedOldUserProfilePicturePath = filepath.Join(oldUserProfilePicturePathLocation, "__DELETED__"+oldUserProfilePicturePathFilename)

				err = os.Rename(oldUserInfo.PROFILE_PIC, deletedOldUserProfilePicturePath)
				if err != nil {
					//500, Incase of DB error
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

			}
		}

		//If the User actually wants to have a Profile Picture
		if emptyProfilePicture != "true" {

			if err != nil {
				//Incase of Error, We make the updated filename as old filename and proceed
				updatedUser.PROFILE_PIC = oldUserInfo.PROFILE_PIC

				//413, If File too large from MaxSizeMiddleware
				if strings.Contains(err.Error(), "request body too large") {
					ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "File Too Large",
					})
					return
				}

			} else {

				//415, If the picture is not of Supported file types
				if filepath.Ext(profilePicture.Filename) != ".png" && filepath.Ext(profilePicture.Filename) != ".jpg" {
					ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
						"error": "Please upload a supported Profile Picture format",
					})
					return
				}

				//413, If picture exceeds Max Picture Size
				if profilePicture.Size > consts.MaxPictureSize {
					ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
						"error": "Profile Picture should be less than " + strconv.Itoa(consts.MaxPictureSize>>20) + "mb!",
					})
					return
				}

				//mark old file for deletion
				if oldUserInfo.PROFILE_PIC != "" {
					//Construct Filepath for old file
					oldUserProfilePicturePathFilename := filepath.Base(oldUserInfo.PROFILE_PIC)
					oldUserProfilePicturePathLocation := filepath.Dir(oldUserInfo.PROFILE_PIC)
					deletedOldUserProfilePicturePath = filepath.Join(oldUserProfilePicturePathLocation, "__DELETED__"+oldUserProfilePicturePathFilename)

					err = os.Rename(oldUserInfo.PROFILE_PIC, deletedOldUserProfilePicturePath)
					if err != nil {
						//500, In case of any FS error
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"error": err.Error(),
						})
						return
					}
				}

				//Construct new Profile Picture path
				updatedUserProfilePicture = filepath.Join(consts.ProfilePicDirectory, UID+filepath.Ext(profilePicture.Filename))
				//Save new Profile Picture
				err = ctx.SaveUploadedFile(profilePicture, updatedUserProfilePicture)
				if err != nil {
					//In case of error, unmark Old file for deletion and stop
					if oldUserInfo.PROFILE_PIC != "" {
						FSerr := os.Rename(deletedOldUserProfilePicturePath, oldUserInfo.PROFILE_PIC)
						if FSerr != nil {
							//500, In case of FS error
							ctx.JSON(http.StatusInternalServerError, gin.H{
								"error": FSerr.Error(),
							})
							return
						}
					}
					//500, for File saving error
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}
				updatedUser.PROFILE_PIC = updatedUserProfilePicture
			}
		}

		//Finally Update the actual User Info
		err = dbUtils.UpdateUserInfo(ctx, UID, updatedUser, db)
		if err != nil {

			//In case of error, we need to delete the files in order to maintain consistency
			var FSerr error

			//Delete the new, Updated File
			if updatedUserProfilePicture != "" {
				FSerr = os.Remove(updatedUserProfilePicture)
				if FSerr != nil {
					//500 in case of FS error
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": FSerr.Error(),
					})
					return
				}
			}
			//Unmark the old File with its old name without "__DELETED__"
			if deletedOldUserProfilePicturePath != "" {
				FSerr = os.Rename(deletedOldUserProfilePicturePath, oldUserInfo.PROFILE_PIC)
				if FSerr != nil {
					//500, in case of FS error
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": FSerr.Error(),
					})
					return
				}
			}

			//404, If No User is found
			if strings.Contains(err.Error(), "No User found") {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No User Found",
				})
				return
			}
			//500, If anything else Happens
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something Went Wrong!",
			})
			return
		}

		//Delete the marked file
		if deletedOldUserProfilePicturePath != "" {
			err = os.Remove(deletedOldUserProfilePicturePath)
			if err != nil {
				//200, but if file deletion goes wrong, We just tell the user to request manual deletion of file later.
				ctx.JSON(http.StatusOK, gin.H{
					"message": "Updated User",
					"warning": "User Updated, but old Profile Picture file still exists. Request manual deletion",
				})
				return
			}

		}

		//200, If All goes well
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User Updated",
		})

	}
}

// Handler to Delete a User, the boolean is keepRecords, stating whether to delete the records with the User
func HandleUserDELETE(db *sql.DB, keepRecords bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Obtain UID From JWT
		UID, err := authUtils.ParseToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = dbUtils.DeleteUser(ctx, UID, db, keepRecords)
		if err != nil {
			//400, Both of those errors will mean the same thing, no User found
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

			//500, for any other errors
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		//200, if All goes well along with whether the records were preserved or Not
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User Successfully Deleted, Records Preserved: " + strconv.FormatBool(keepRecords),
		})
	}
}

// Handler for when the User searches for another User
func HandleUserSEARCH(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		searchResults := []*dbUtils.USER_PUBLIC{}

		//Get Query
		query := ctx.Param("query")
		if query == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "No Query Found",
			})
			return
		}

		//Obtain Search Results
		searchResults, err := dbUtils.SearchUser(ctx, strings.Split(query, " "), db)
		if err != nil {
			//500, for any DB errors
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if len(searchResults) < 1 {
			//404, if no User is found for the query
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No User Found for the Query",
			})
			return
		}

		//200, if All goes well
		ctx.JSON(http.StatusOK, gin.H{
			"message": searchResults,
		})

	}
}
