package utilities

import (
	"log"
	"os"
	"strings"
)

func FindFileFromDirectory(ID string) string {

	files, err := os.ReadDir("../uploads")

	if err != nil {
		log.Println("Something Went Wrong")
		return ""
	}
	for _, file := range files {
		fileParts := strings.Split(file.Name(), ".")
		if len(fileParts) <= 1 {
			filename := strings.Join(fileParts, "")
			if filename == ID {
				return file.Name()
			}
		}
		var fileNameArray []string
		fileNameArray = append(fileNameArray, fileParts[:len(fileParts)-1]...)
		filename := strings.Join(fileNameArray, ".")
		if filename == ID {
			return file.Name()
		}

	}
	return ""
}
