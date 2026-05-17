package utilities

import (
	"log"
	"os"

	"github.com/AyushGlitchedOut/Docummunity/consts"
)

//Create all the folders required upon App Initialization if not already existing

func CreateUploadsFolder() {

	//Make Uploads directory,
	err := os.MkdirAll("./uploads", 0o755)
	if err != nil {
		log.Println(err)
	}

	//check if file exists or not
	_, err = os.Stat(consts.DBLocation)
	if err == nil {
		log.Println("DB already exists")
	} else {
		//Create DDATABASE.db File
		file, err := os.Create(consts.DBLocation)
		if err != nil {
			log.Fatal(err)
		}
		//close the connection to the file since we dont have anything to do with the file itself, and another method is gonna use the file
		file.Close()
	}

	//Create directories for all 3 types of file to store
	err = os.MkdirAll(consts.ProfilePicDirectory, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(consts.FileDirectory, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(consts.PreviewImgDirectory, 0o755)
	if err != nil {
		log.Fatal(err)
	}

}
