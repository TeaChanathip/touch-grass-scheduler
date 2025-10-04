package authfx

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/models"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/mytypes"
	usersfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthControllerParams struct {
	fx.In
	Logger      *zap.Logger
	AuthService *AuthService
	UserService *usersfx.UserService
}

type AuthController struct {
	Logger      *zap.Logger
	AuthService *AuthService
	UserService *usersfx.UserService
}

func NewAuthController(params AuthControllerParams) *AuthController {
	return &AuthController{
		Logger:      params.Logger,
		AuthService: params.AuthService,
		UserService: params.UserService,
	}
}

// ======================== REQUEST BODY ========================

type RegisterBody struct {
	Role       mytypes.UserRole   `json:"role" binding:"required,oneof='student' 'teacher' 'guardian'"` // Not allow Admin to be registered
	FirstName  string             `json:"first_name" binding:"required,max=128,alpha"`
	MiddleName string             `json:"middle_name" binding:"omitempty,max=128,alpha"`
	LastName   string             `json:"last_name" binding:"omitempty,max=128,alpha"`
	Phone      string             `json:"phone" binding:"required,e164"`
	Gender     mytypes.UserGender `json:"gender" binding:"required,oneof=''male' 'female' 'other prefer_not_to_say'"`
	Email      string             `json:"email" binding:"required,email"`
	Password   string             `json:"password" binding:"required,min=8,max=64"`
	SchoolNum  string             `json:"school_num" binding:"omitempty,number,max=16"` // Be either student_num or teacher_num
}

func (rb RegisterBody) ToUserModel() *models.User {
	return &models.User{
		Role:       rb.Role,
		FirstName:  rb.FirstName,
		MiddleName: rb.MiddleName,
		LastName:   rb.LastName,
		Phone:      rb.Phone,
		Gender:     rb.Gender,
		Email:      rb.Email,
		Password:   rb.Password,
		SchoolNum:  rb.SchoolNum,
	}
}

// ======================== METHODS ========================

func (controller *AuthController) Register(ctx *gin.Context) {
	// Validate request body
	var registerBody RegisterBody
	if err := ctx.ShouldBindBodyWithJSON(&registerBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate school number requirements
	if err := validateSchoolNum(registerBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Add logic to check if SchoolNumber is valid
	// TODO: Send the verification link to the user's email

	// Create new user
	user := registerBody.ToUserModel()
	if err := controller.UserService.CreateUser(user); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Remove sensitive attributes
	publicUser := user.ToPublic()

	// Generate JWT token
	token, err := controller.AuthService.GenerateToken(user)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":  publicUser,
		"token": token,
	})
}

// ======================== Helper Functions ========================

func validateSchoolNum(registerBody RegisterBody) error {
	schoolPersonnelRoles := []mytypes.UserRole{
		mytypes.UserRoleStudent,
		mytypes.UserRoleTeacher,
	}

	isSchoolPersonnel := slices.Contains(schoolPersonnelRoles, registerBody.Role)

	if isSchoolPersonnel && registerBody.SchoolNum == "" {
		return fmt.Errorf("%s must provide school_num", registerBody.Role)
	}

	if !isSchoolPersonnel && registerBody.SchoolNum != "" {
		return fmt.Errorf("%s should not provide school_num", registerBody.Role)
	}

	return nil
}
