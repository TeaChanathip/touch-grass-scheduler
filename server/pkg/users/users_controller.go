package usersfx

import (
	"fmt"
	"net/http"

	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UsersControllerParams struct {
	fx.In
	Logger      *zap.Logger
	UserService UserServiceInterface
}

type UsersController struct {
	Logger      *zap.Logger
	UserService UserServiceInterface
}

func NewUsersController(params UsersControllerParams) *UsersController {
	return &UsersController{
		Logger:      params.Logger,
		UserService: params.UserService,
	}
}

// ======================== REQUEST BODY ========================

type UpdateUserBody struct {
	FirstName  *string           `json:"first_name" binding:"omitempty,max=128,alpha"`
	MiddleName *string           `json:"middle_name" binding:"omitempty,max=128,len=0|alpha"`
	LastName   *string           `json:"last_name" binding:"omitempty,max=128,len=0|alpha"`
	Phone      *string           `json:"phone" binding:"omitempty,e164"`
	Gender     *types.UserGender `json:"gender" binding:"omitempty,oneof=''male' 'female' 'other' 'prefer_not_to_say'"`
}

// ======================== METHODS ========================

func (controller *UsersController) GetMe(ctx *gin.Context) {
	// Get userID Context that set by AuthMiddleware
	_userID, _ := ctx.Get("user_id")
	userID, err := uuid.Parse(_userID.(string))
	if err != nil {
		controller.Logger.Debug("Failed to parse userID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	user, err := controller.UserService.GetPublicUserByID(userID)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (controller *UsersController) GetUserByID(ctx *gin.Context) {
	// Get id from params
	id := ctx.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		controller.Logger.Debug("Failed to parse userID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "request ID is not UUID"})
		return
	}

	user, err := controller.UserService.GetPublicUserByID(userID)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (controller *UsersController) UpdateUserByID(ctx *gin.Context) {
	// Get userID Context that set by AuthMiddleware
	_userID, _ := ctx.Get("user_id")
	userID, err := uuid.Parse(_userID.(string))
	if err != nil {
		controller.Logger.Debug("Failed to parse userID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	validatedBody, _ := ctx.Get("validatedBody")
	updateUserBody, _ := validatedBody.(*UpdateUserBody)
	controller.Logger.Debug(fmt.Sprintf("%+v\n", updateUserBody))

	user, err := controller.UserService.UpdateUserByID(userID, updateUserBody)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (controller *UsersController) GetUploadAvatarSignedURL(ctx *gin.Context) {
	// Get userID Context that set by AuthMiddleware
	_userID, _ := ctx.Get("user_id")
	userID, err := uuid.Parse(_userID.(string))
	if err != nil {
		controller.Logger.Debug("Failed to parse userID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	response, err := controller.UserService.GetUploadAvatarSignedURL(userID)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (controller *UsersController) HandleAvatarUpload(ctx *gin.Context) {
	// Get userID Context that set by AuthMiddleware
	_userID, _ := ctx.Get("user_id")
	userID, err := uuid.Parse(_userID.(string))
	if err != nil {
		controller.Logger.Debug("Failed to parse userID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	url, err := controller.UserService.HandleAvatarUpload(ctx.Request.Context(), userID)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"avatar_url": url.String()})
}
