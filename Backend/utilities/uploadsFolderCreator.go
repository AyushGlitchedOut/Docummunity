package utilities

import (
	"log"
	"os"
)

const (
	UploadsDirectory    = "./uploads"
	DBLocation          = UploadsDirectory + "/DATABASE.db"
	ProfilePicDirectory = UploadsDirectory + "/PROFILE_PIC/"
	FileDirectory       = UploadsDirectory + "/FILES/"
	PreviewImgDirectory = UploadsDirectory + "/PREVIEW/"
)

func CreateUploadsFolder() {
	err := os.MkdirAll("./uploads", 0o755)
	if err != nil {
		log.Println(err)
	}

	//check if file exists or not
	_, err = os.Stat(DBLocation)
	if err == nil {
		log.Println("DB already exists")
	} else {
		file, err := os.Create(DBLocation)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	err = os.MkdirAll(ProfilePicDirectory, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(FileDirectory, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(PreviewImgDirectory, 0o755)
	if err != nil {
		log.Fatal(err)
	}

}
