package httputil

import (
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewError example
func NewError(c *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	c.JSON(status, er)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func CheckFileContentType(fileHeader *multipart.FileHeader, checkType string) bool {
	if fileHeader == nil {
		return false
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	detectedType, err := getFileContentType(file)
	if err != nil {
		log.Panicln(err)
	}

	return detectedType == checkType
}

func getFileContentType(out multipart.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
