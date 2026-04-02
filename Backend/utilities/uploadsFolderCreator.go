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
}
