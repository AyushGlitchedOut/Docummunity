package main

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/server"

	"google.golang.org/api/option"
)

func main() {

	//Configure Firebase Admin SDK
	opt := option.WithCredentialsFile("./firebase_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)

	//CREATE UPLOADS DIRECTORY
	err = os.MkdirAll("./uploads", 0o755)
	if err != nil {
		log.Println(err)
	}

	//Init DB
	db, err := dbUtils.InitDatabase()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("-----------------DOCUMMUNITY BACKEND----------------------------")
	port := ":8080"

	server := server.InitServer(port, db, app)

	if err := server.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// client, err := app.Auth(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// userInfo := (&auth.UserToCreate{}).Email("test@gmail.com").Password("ayush2011").DisplayName("TESTER")

	// user, err := client.CreateUser(context.Background(), userInfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(user.UID)

}
