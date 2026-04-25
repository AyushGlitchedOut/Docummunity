package utilities

import (
	"log"
	"os"

	"github.com/AyushGlitchedOut/Docummunity/consts"
)

const ()

func CreateUploadsFolder() {
	err := os.MkdirAll("./uploads", 0o755)
	if err != nil {
		log.Println(err)
	}

	//check if file exists or not
	_, err = os.Stat(consts.DBLocation)
	if err == nil {
		log.Println("DB already exists")
	} else {
		file, err := os.Create(consts.DBLocation)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

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
