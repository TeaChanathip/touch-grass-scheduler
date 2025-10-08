package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func RequestBodyValidator(logger *zap.Logger, name string, structType any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Create a new instance of the struct type using reflection
		structValue := reflect.New(reflect.TypeOf(structType))
		validateBody := structValue.Interface()

		if err := ctx.ShouldBindBodyWithJSON(&validateBody); err != nil {
			logger.Debug(fmt.Sprintf("Validation error on %s request", name), zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": parseValidationErrors(err)})
			ctx.Abort()
			return
		}

		// Store the validated body in context for the controller to use
		ctx.Set("validatedBody", validateBody)
		ctx.Next()
	}
}

// ======================== HELPER FUNCTIONS ========================

func parseValidationErrors(err error) map[string]string {
	valErrors := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		// If it's not a validation error, return generic message
		return map[string]string{"general": "Invalid request format"}
	}

	for _, fieldError := range validationErrors {
		valErrors[fieldError.Field()] = getValidationMessage(fieldError)
	}

	return valErrors
}

// returns user-friendly error messages
func getValidationMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "Please provide a valid email address"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, param)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "e164":
		return "Please provide a valid phone number (e.g., +1234567890)"
	case "number":
		return fmt.Sprintf("%s must be a valid number", field)
	case "oneof":
		values := strings.ReplaceAll(param, " ", ", ")
		return fmt.Sprintf("%s must be one of: %s", field, values)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
