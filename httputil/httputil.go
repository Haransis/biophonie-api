package httputil

import (
	"log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype/matchers"
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

func CheckFileContentType(fileHeader *multipart.FileHeader, matcher matchers.Matcher) bool {
	if fileHeader == nil {
		return false
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()

	buffer := make([]byte, 261)

	_, err = file.Read(buffer)
	if err != nil {
		return false
	}

	return matcher(buffer)
}
