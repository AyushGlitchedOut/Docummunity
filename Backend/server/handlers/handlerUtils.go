package handlers

import (
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// Helper function used in update Routes where old files need to be marked for deletion, parses and constructs routes for renaming for marking
func deletedFilePathMaker(filePath string) string {
	name := filepath.Base(filePath)
	location := filepath.Dir(filePath)

	//append __DELETED__ to the file
	deletedFilePath := filepath.Join(location, "__DELETED__"+name)
	return deletedFilePath
}

// File Type detection using file headers (Magic NUmbers) or MIME type
func filetypeDetector(reqFile *multipart.FileHeader, allowedTypes []string) (bool, error) {

	//Open the file from the gievn Filehader reference
	file, err := reqFile.Open()
	if err != nil {
		return false, err
	}
	defer file.Close()

	//Read the first 512 bytes of the file
	header := make([]byte, 512)
	n, err := file.Read(header)

	//Check for End of file since very small files may be lesser than 512 bytes
	if err != nil && err != io.EOF {
		return false, nil
	}

	//DEV NOTE: Remember, array[a:b] returns a slice from index a (inclusive) to b (exclusive). So, below, header[:n] returns an array from 0 (starting) to n, the number of bytes actually read by the reader. Now, since b is exclusive, you may think that n+1 should be used, but since arrays are 0-indexed, even with nth element excluded, we get the required no. of bytes

	//detect MIME type from the readed bytes
	contentType := http.DetectContentType(header[:n])

	//scan the given list of allowed types and check if the found MIME type matches
	for _, val := range allowedTypes {
		if val == contentType {
			return true, nil
		}
	}

	return false, nil
}
