package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var InternalServerError = "internal server error, please try again later"

func (c *Controller) HandleErrors(ctx *gin.Context) {
	ctx.Next()

	e := ctx.Errors.Last()
	if e != nil {
		switch e.Type {
		case gin.ErrorTypePublic:
			ctx.JSON(-1, errMsg(e.Error()))
		case gin.ErrorTypePrivate:
			ctx.JSON(-1, errMsg(InternalServerError))
			log.Printf("%s %s", e.Meta, e.Err)
		case gin.ErrorTypeBind:
			if err, ok := e.Err.(*validator.ValidationErrors); ok {
				ctx.JSON(-1, errMsg(ValidationErrorToText((*err)[0])))
			} else if err, ok := e.Err.(*json.SyntaxError); ok {
				ctx.JSON(-1, fmt.Sprintf("syntax error (%s) at byte offset %d", err.Error(), err.Offset))
			} else {
				ctx.JSON(-1, errMsg(fmt.Sprintf("malformed request: %s", e.Err)))
			}
		case gin.ErrorTypeAny:
			if pqerr, ok := e.Err.(*pq.Error); ok {
				switch pqerr.Code.Name() {
				case "unique_violation":
					ctx.JSON(http.StatusConflict, errMsg("element was already present"))
				case "name_too_long":
					ctx.JSON(http.StatusBadRequest, errMsg("name was too long"))
				default:
					ctx.JSON(http.StatusInternalServerError, errMsg(InternalServerError))
					log.Printf("%s %s", e.Meta, e.Err)
				}
			} else {
				switch e.Err {
				case sql.ErrNoRows:
					ctx.JSON(http.StatusNotFound, errMsg("element was not found"))
				case bcrypt.ErrMismatchedHashAndPassword:
					ctx.JSON(http.StatusUnauthorized, errMsg("wrong password"))
				default:
					ctx.JSON(http.StatusInternalServerError, errMsg(InternalServerError))
					log.Printf("%s %s", e.Meta, e.Err)
				}
			}
		}
	}
}

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	return strings.ToLower(str)
}

func Split(src string) string {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	var entries []string
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}

	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	for index, word := range entries {
		if index == 0 {
			entries[index] = UcFirst(word)
		} else {
			entries[index] = LcFirst(word)
		}
	}
	justString := strings.Join(entries, " ")
	return justString
}

func ValidationErrorToText(e validator.FieldError) string {
	word := Split(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", word)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", word, e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", word, e.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", word, e.Param())
	}
	return fmt.Sprintf("%s is not valid", word)
}

// only for swaggo
type ErrMsg struct {
	Message string `example:"malformed request"`
}

func errMsg(msg string) gin.H {
	return gin.H{
		"message": msg,
	}
}
