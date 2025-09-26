package main

import (
	"github.com/TeaChanathip/touch-grass-scheduler/internal/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	flagConfig := configs.LoadFlags()
	configs.LoadConfig(flagConfig.Environment)
	configs.GetDatabase()

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()
	r.Run()
}
