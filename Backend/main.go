package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AyushGlitchedOut/Docummunity/authUtils"
	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/server"
	"github.com/AyushGlitchedOut/Docummunity/utilities"
)

func main() {
	port := ":8080"

	//Configure Firebase Admin SDK
	firebase := authUtils.FirebaseAppCreator()

	//CREATE UPLOADS DIRECTORY
	utilities.CreateUploadsFolder()

	DB, err := dbUtils.InitializeDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	fmt.Println("-----------------DOCUMMUNITY BACKEND----------------------------")

	app := server.InitServer(port, DB, firebase)

	if err := app.ListenAndServe(); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
