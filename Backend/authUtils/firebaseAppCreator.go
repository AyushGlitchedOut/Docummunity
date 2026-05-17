package authUtils

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Simply takes the "./firebase_key.json" from the current directory and returns a firebaseApp instance. Crashes the server if error occurs
func FirebaseAppCreator() *firebase.App {

	//Initialize App with Settings File
	opt := option.WithCredentialsFile("./firebase_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatal("Error Initializing The ADMIN SDK")
	}

	return app
}
