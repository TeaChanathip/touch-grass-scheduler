package usersfx

import (
	"net/http"

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

// ======================== METHODS ========================

func (controller *UsersController) GetUser(ctx *gin.Context) {
	// Get user's infomation from Context
	_userID, _ := ctx.Get("user_id")
	userID, err := uuid.Parse(_userID.(string))
	if err != nil {
		controller.Logger.Debug("Failed to parse userID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	user, err := controller.UserService.GetUserByID(userID)
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

	user, err := controller.UserService.GetUserByID(userID)
	if err != nil {
		common.HandleBusinessLogicErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// func (controller *UsersController) UpdateUser(ctx *gin.Context) {
// }

// func (controller *UsersController) DeleteUser(ctx *gin.Context) {
// }
