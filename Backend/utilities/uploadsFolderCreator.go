package utilities

import (
	"log"
	"os"
)

func CreateUploadsFolder() {
	err := os.MkdirAll("./uploads", 0o755)
	if err != nil {
		log.Println(err)
	}

	//check if file exists or not
	_, err = os.Stat("./uploads/DATABASE.db")
	if err == nil {
		log.Println("DB already exists")
		return
	}
	if !os.IsNotExist(err) {
		log.Fatal("Error Creating File")
	}

	file, err := os.Create("./uploads/DATABASE.db")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

}
