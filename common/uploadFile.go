package common

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) ([]string, error) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	destinationDir := "./uploads"
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		err := os.Mkdir(destinationDir, os.ModePerm)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating directory: %s", err.Error()))
			return nil, err
		}
	}
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
