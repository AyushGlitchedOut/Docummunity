package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AyushGlitchedOut/Docummunity/authUtils"
	"github.com/AyushGlitchedOut/Docummunity/consts"
	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/server"
	"github.com/AyushGlitchedOut/Docummunity/server/handlers"
	"github.com/AyushGlitchedOut/Docummunity/utilities"
)

func main() {

	//Configure Firebase Admin SDK
	firebase := authUtils.FirebaseAppCreator()

	//CREATE UPLOADS DIRECTORY
	utilities.CreateUploadsFolder()

	//Initialize Database
	DB, err := dbUtils.InitializeDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	log.Println("-----------------------------------------DOCUMMUNITY BACKEND-------------------------------------")

	app := server.InitServer(consts.Port, DB, firebase)

	//Start the ClientList Cleanup service
	handlers.StartClientListCleanupService()

	go func() {
		//Run the Server
		if err := app.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal("Failed to run server: ", err)
			}
		}
	}()

	//Setup For Graceful Shutdown:
	exitDetector := make(chan os.Signal, 1)
	signal.Notify(exitDetector, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(exitDetector)

	//Detect termination message and receive it
	message := <-exitDetector

	//Set Timelimit for the shutdown process
	cancelCtx, cancel := context.WithTimeout(context.Background(), consts.ShutdownTimeLimit)
	defer cancel()

	//Shutdown the HTTP server
	err = app.Shutdown(cancelCtx)
	if err != nil {
		log.Println("Error Shutting Down Server: ", err)
	}

	log.Println("SHUTTING DOWN SERVER: ", message)

}
