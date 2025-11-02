package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomError struct {
	StatusCode int
	Message    string
}

func (e CustomError) Error() string {
	return e.Message
}

var (
	// 500 Internal Server Errors
	ErrTokenGeneration = CustomError{StatusCode: http.StatusInternalServerError, Message: "accessToken generation failed"}
	ErrPasswordHashing = CustomError{StatusCode: http.StatusInternalServerError, Message: "password hashing failed"}
	ErrDatabase        = CustomError{StatusCode: http.StatusInternalServerError, Message: "database error"}
	ErrMailHTMLSetting = CustomError{StatusCode: http.StatusInternalServerError, Message: "set HTML to mail message failed"}
	ErrMailSending     = CustomError{StatusCode: http.StatusInternalServerError, Message: "mail sending failed"}
	ErrVariableParsing = CustomError{StatusCode: http.StatusInternalServerError, Message: "variable parsing failed"}
	ErrUUIDGenerating  = CustomError{StatusCode: http.StatusInternalServerError, Message: "UUID generation failed"}
	ErrStorage         = CustomError{StatusCode: http.StatusInternalServerError, Message: "storage error"}
	ErrURLSigning      = CustomError{StatusCode: http.StatusInternalServerError, Message: "signing url failed"}

	// 400 Bad Request
	ErrDuplicatedEmail          = CustomError{StatusCode: http.StatusBadRequest, Message: "email already exists"}
	ErrActionTokenParsing       = CustomError{StatusCode: http.StatusBadRequest, Message: "actionToken parsing failed"}
	ErrActionTokenClaimsGetting = CustomError{StatusCode: http.StatusBadRequest, Message: "actionToken getting claims failed"}

	// 401 Authentication Errors
	ErrInvalidCredentials = CustomError{StatusCode: http.StatusUnauthorized, Message: "invalid credentials"}

	// 404 Not Found
	ErrUserNotFound          = CustomError{StatusCode: http.StatusNotFound, Message: "user not found"}
	ErrStorageObjectNotFound = CustomError{StatusCode: http.StatusNotFound, Message: "object not found"}
	ErrPendingUploadNotFound = CustomError{StatusCode: http.StatusNotFound, Message: "pending upload not found"}
)

// ======================== HELPER FUNCTIONS ========================

func HandleBusinessLogicErr(ctx *gin.Context, err error) {
	if customErr, ok := err.(CustomError); ok {
		ctx.JSON(customErr.StatusCode, gin.H{"error": customErr.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error"})
	}
}
