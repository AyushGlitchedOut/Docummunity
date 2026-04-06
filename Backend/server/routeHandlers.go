package server

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePING(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Server Active!",
	})
}

// Test functions
func VerifyTest(ctx *gin.Context) {

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
