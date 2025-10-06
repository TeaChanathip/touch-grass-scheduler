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
	ErrTokenGeneration = CustomError{StatusCode: http.StatusInternalServerError, Message: "token generation failed"}
	ErrPasswordHashing = CustomError{StatusCode: http.StatusInternalServerError, Message: "password hashing failed"}
	ErrDatabase        = CustomError{StatusCode: http.StatusInternalServerError, Message: "database error"}

	// 400 Bad Request
	ErrDuplicatedEmail = CustomError{StatusCode: http.StatusBadRequest, Message: "email already exists"}

	// 401 Authentication Errors
	ErrInvalidCredentials = CustomError{StatusCode: http.StatusUnauthorized, Message: "invalid credentials"}

	// 404 Not Found
	ErrUserNotFound = CustomError{StatusCode: http.StatusNotFound, Message: "user not found"}
)

// ======================== HELPER FUNCTIONS ========================

func HandleBusinessLogicErr(ctx *gin.Context, err error) {
	if customErr, ok := err.(CustomError); ok {
		ctx.JSON(customErr.StatusCode, gin.H{"error": customErr.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
	}
}
