package common

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func CheckDestinationDir() (string, error) {
	destinationDir := "./uploads"
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		err := os.Mkdir(destinationDir, os.ModePerm)
		if err != nil {
			return "", nil
		}
	}
	return destinationDir, nil
}

func UploadMultiFile(c *gin.Context) ([]string, error) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	destinationDir, _ := CheckDestinationDir()
	var filePaths []string
	for _, file := range files {
		savePath := filepath.Join(destinationDir, file.Filename)
		err := c.SaveUploadedFile(file, savePath)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error saving file '%s': %s", file.Filename, err.Error()))
			return nil, err
		}
		filePaths = append(filePaths, savePath)
	}
	return filePaths, nil
}

func UploadFile(c *gin.Context) (string, error) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	destinationDir, _ := CheckDestinationDir()
	savePath := filepath.Join(destinationDir, file.Filename)
	err := c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error saving file '%s': %s", file.Filename, err.Error()))
		return "", err
	}
	return savePath, nil
}
