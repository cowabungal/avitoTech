package handler

import (
	"avitoTech/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	// gin.Default()

	router.Use(gin.Logger())
	//router.NoRoute(h.notFound)
	router.POST("/top-up", h.TopUp)
	router.GET("/balance", h.Balance)
	router.POST("/debit", h.Debit)
	router.POST("/transfer", h.Transfer)
	router.GET("/transaction", h.Transaction)
	return router
}
