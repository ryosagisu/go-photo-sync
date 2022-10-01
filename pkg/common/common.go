package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func ListLocalImages(imagePath string) map[string]bool {
	imageFiles := make(map[string]bool)
	files, err := ioutil.ReadDir(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() || file.Size() == 0 {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", imagePath, file.Name())
		if !IsValidImage(filePath) {
			continue
		}
		imageFiles[strings.TrimSuffix(file.Name(), ".jpg")] = true
	}
	return imageFiles
}

func IsValidImage(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open image: %v\n", err)
		return false
	}
	defer f.Close()

	buff := make([]byte, 512)
	if _, err := f.Read(buff); err != nil {
		log.Printf("failed to read image: %v\n", err)
		return false
	}

	return http.DetectContentType(buff) == "image/jpeg"
}

func DeleteLocalFile(filePath string, localImages map[string]bool) {
	for k, v := range localImages {
		if v {
			imageFile := fmt.Sprintf(k, filePath)
			err := os.Remove(imageFile)
			if err != nil {
				log.Printf("failed to delete local files: %v\n", err)
				return
			}
			log.Printf("%s deleted\n", imageFile)
		}
	}
}
