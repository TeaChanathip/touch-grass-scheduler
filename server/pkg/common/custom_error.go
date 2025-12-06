package common

import (
	"errors"
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
	ErrDatabase = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "database error",
	}
	ErrStorage = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "storage error",
	}
	ErrTokenGeneration = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "failed generating access token",
	}
	ErrPasswordHashing = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "failed hashing password",
	}
	ErrMailHTMLSetting = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "failed setting mail html",
	}
	ErrMailSending = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "failed sending mail",
	}
	ErrVariableParsing = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "failed parsing variable",
	}
	ErrUUIDGeneration = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "failed generating uuid",
	}
	ErrURLSigning = CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "failed signing url",
	}

	// 400 Bad Request
	ErrDuplicatedEmail = CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    "email already exists",
	}
	ErrActionTokenParsing = CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    "failed parsing action token",
	}
	ErrActionTokenClaimsGetting = CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    "failed getting action token claims",
	}
	ErrActionTokenExpired = CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    "action token already expired",
	}

	// 401 Authentication/Authorization Errors
	ErrInvalidCredentials = CustomError{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid credentials",
	}

	// 404 Not Found
	ErrUserNotFound = CustomError{
		StatusCode: http.StatusNotFound,
		Message:    "user not found",
	}
	ErrStorageObjectNotFound = CustomError{
		StatusCode: http.StatusNotFound,
		Message:    "object not found",
	}
	ErrPendingUploadNotFound = CustomError{
		StatusCode: http.StatusNotFound,
		Message:    "pending upload not found",
	}
)

// ======================== HELPER FUNCTIONS ========================

func HandleBusinessLogicErr(ctx *gin.Context, err error) {
	var customErr CustomError
	if errors.As(err, &customErr) {
		ctx.JSON(customErr.StatusCode, gin.H{"error": customErr.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error"})
	}
}
