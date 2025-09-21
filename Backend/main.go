package main

import (
	"fmt"
	"net/http"

	"github.com/AyushGlitchedOut/Docummunity/Backend/consts"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	//Change to gin.ReleaseMode in production
	gin.SetMode(gin.DebugMode)

	router := gin.New()

	//Give the CORS config in the constants file
	router.Use(cors.New(consts.CorsConfig))

	//Use Recovery
	router.Use(gin.Recovery())

	//Test Route
	router.GET("/test", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello Docummunity!!")
	})

	//START THE SERVER
	fmt.Print("-------------------STARTING THE SERVER------------------\n\n\n")
	err := router.Run(":8000")

	if err != nil {
		fmt.Println("Failed to Start Server:", err)
	}
}
