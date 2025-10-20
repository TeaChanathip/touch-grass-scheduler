package authfx

import (
	"fmt"
	"net/http"
	"slices"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthControllerParams struct {
	fx.In
	FlagConfig  *configfx.FlagConfig
	AppConfig   *configfx.AppConfig
	Logger      *zap.Logger
	AuthService AuthServiceInterface
}

type AuthController struct {
	FlagConfig  *configfx.FlagConfig
	AppConfig   *configfx.AppConfig
	Logger      *zap.Logger
	AuthService AuthServiceInterface
}

func NewAuthController(params AuthControllerParams) *AuthController {
	return &AuthController{
		FlagConfig:  params.FlagConfig,
		AppConfig:   params.AppConfig,
		Logger:      params.Logger,
		AuthService: params.AuthService,
	}
}

// ======================== REQUEST BODY ========================

type GetRegistrationMailBody struct {
	Email string `json:"email" binding:"required,email"`
}

type RegisterBody struct {
	Role       types.UserRole   `json:"role" binding:"required,oneof='student' 'teacher' 'guardian'"` // Not allow Admin to be registered
	FirstName  string           `json:"first_name" binding:"required,max=128,alpha"`
	MiddleName string           `json:"middle_name" binding:"omitempty,max=128,alpha"`
	LastName   string           `json:"last_name" binding:"omitempty,max=128,alpha"`
	Phone      string           `json:"phone" binding:"required,e164"`
	Gender     types.UserGender `json:"gender" binding:"required,oneof=''male' 'female' 'other' 'prefer_not_to_say'"`
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
		Password:   rb.Password,
		SchoolNum:  rb.SchoolNum,
	}
}

type LoginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

// ======================== METHODS ========================

func (controller *AuthController) GetRegistrationMail(ctx *gin.Context) {
	// Get validated body from context that set by RequestBodyValidator
	validatedBody, _ := ctx.Get("validatedBody")
	getRegistrationMailBody, _ := validatedBody.(*GetRegistrationMailBody)

	// Business logic
	err := controller.AuthService.GetRegistrationMail(getRegistrationMailBody.Email)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *AuthController) Register(ctx *gin.Context) {
	// Get registrationToken from params
	registrationTokenString, exists := ctx.Params.Get("registrationToken")
	if !exists {
		controller.Logger.Debug("The registrationToken is not exists.")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "registrationToken is required"})
		return
	}

	// Get validated body from context that set by RequestBodyValidator
	validatedBody, _ := ctx.Get("validatedBody")
	registerBody, _ := validatedBody.(*RegisterBody)

	// Validate school number requirements
	if err := validateSchoolNum(registerBody); err != nil {
		controller.Logger.Debug("Validation error on register request:", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Business logic
	user, accessToken, err := controller.AuthService.Register(registrationTokenString, registerBody)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	// Convert user struct to map with snake_case key
	userMap, err := common.StructToSnakeMap(user)
	if err != nil {
		controller.Logger.Error("Internal error while converting user struct to map:", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	controller.setAccessTokenCookie(ctx, accessToken)
	ctx.JSON(http.StatusCreated, gin.H{
		"user": userMap,
	})
}

func (controller *AuthController) Login(ctx *gin.Context) {
	// Get validated body from context that set by RequestBodyValidator
	validatedBody, _ := ctx.Get("validatedBody")
	loginBody, _ := validatedBody.(*LoginBody)

	// Business logic
	user, accessToken, err := controller.AuthService.Login(loginBody)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	// Convert user struct to map with snake_case key
	userMap, err := common.StructToSnakeMap(user)
	if err != nil {
		controller.Logger.Error("Internal error while converting user struct to map:", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	controller.setAccessTokenCookie(ctx, accessToken)
	ctx.JSON(http.StatusOK, gin.H{
		"user": userMap,
	})
}

// ======================== HELPER METHODS ========================

func (controller *AuthController) setAccessTokenCookie(ctx *gin.Context, accessToken string) {
	isProduction := controller.FlagConfig.Environment == "production"
	maxAge := controller.AppConfig.JWTExpiresIn * 3600

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"accessToken",
		accessToken,
		maxAge,
		"/",
		controller.AppConfig.AppDomain,
		isProduction,
		true,
	)
}

// ======================== HELPER FUNCTIONS ========================

func validateSchoolNum(registerBody *RegisterBody) error {
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
