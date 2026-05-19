package main

import (
	"context"
	"log"

	"github.com/AyushGlitchedOut/Docummunity/authUtils"
	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/server"
	"github.com/AyushGlitchedOut/Docummunity/server/handlers"
	"github.com/AyushGlitchedOut/Docummunity/utilities"
)

func main() {
	port := ":8080"

	//Configure Firebase Admin SDK
	firebase := authUtils.FirebaseAppCreator()

	//CREATE UPLOADS DIRECTORY
	utilities.CreateUploadsFolder()

	//Initialize Database
	DB, err := dbUtils.InitializeDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//make sure DB connection closes once the App closes
	defer DB.Close()

	log.Println("-----------------------------------------DOCUMMUNITY BACKEND-------------------------------------")

	//Initialize the App instance
	app := server.InitServer(port, DB, firebase)

	//Start the ClientList Cleanup service
	handlers.StartClientListCleanupService()

	//Run the Server, which will pause the main program at this point and not execute any code written below this
	if err := app.ListenAndServe(); err != nil {
		log.Fatal("Failed to run server: ", err)
	}

	//-------------------NO CODE BELOW HERE------------------//
}
