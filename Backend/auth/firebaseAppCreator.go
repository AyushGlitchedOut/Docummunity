package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func FirebaseAppCreator() *firebase.App {
	opt := option.WithCredentialsFile("./firebase_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("Error Initializing the ADMIN SDK")
	}
	return app
}
