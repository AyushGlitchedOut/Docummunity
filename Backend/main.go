package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AyushGlitchedOut/Docummunity/auth"
	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/server"
	"github.com/AyushGlitchedOut/Docummunity/utilities"
)

func main() {
	port := ":8080"

	//Configure Firebase Admin SDK
	firebase := auth.FirebaseAppCreator()

	//CREATE UPLOADS DIRECTORY
	utilities.CreateUploadsFolder()

	DB, err := dbUtils.InitializeDB(context.Background())
	if err != nil {
		log.Println(err)
	}
	defer DB.Close()

	fmt.Println("-----------------DOCUMMUNITY BACKEND----------------------------")

	server := server.InitServer(port, DB, firebase)

	if err := server.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
