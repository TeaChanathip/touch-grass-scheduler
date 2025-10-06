package authfx

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/models"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
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
	Role       types.UserRole   `json:"role" binding:"required,oneof='student' 'teacher' 'guardian'"` // Not allow Admin to be registered
	FirstName  string           `json:"first_name" binding:"required,max=128,alpha"`
	MiddleName string           `json:"middle_name" binding:"omitempty,max=128,alpha"`
	LastName   string           `json:"last_name" binding:"omitempty,max=128,alpha"`
	Phone      string           `json:"phone" binding:"required,e164"`
	Gender     types.UserGender `json:"gender" binding:"required,oneof=''male' 'female' 'other prefer_not_to_say'"`
	Email      string           `json:"email" binding:"required,email"`
	Password   string           `json:"password" binding:"required,min=8,max=64"`
	SchoolNum  string           `json:"school_num" binding:"omitempty,number,max=16"` // Be either student_num or teacher_num
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

type LoginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

// ======================== METHODS ========================

func (controller *AuthController) Register(ctx *gin.Context) {
	// Validate request body
	var registerBody RegisterBody
	if err := ctx.ShouldBindBodyWithJSON(&registerBody); err != nil {
		controller.Logger.Debug("Validation error on register request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate school number requirements
	if err := validateSchoolNum(registerBody); err != nil {
		controller.Logger.Debug("Validation error on register request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Business logic
	user, token, err := controller.AuthService.Register(registerBody)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

func (controller *AuthController) Login(ctx *gin.Context) {
	// Validate request body
	var loginBody LoginBody
	if err := ctx.ShouldBindBodyWithJSON(&loginBody); err != nil {
		controller.Logger.Debug("Validation error on login request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Business logic
	user, token, err := controller.AuthService.Login(loginBody)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// ======================== HELPER FUNCTIONS ========================

func validateSchoolNum(registerBody RegisterBody) error {
	schoolPersonnelRoles := []types.UserRole{
		types.UserRoleStudent,
		types.UserRoleTeacher,
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
