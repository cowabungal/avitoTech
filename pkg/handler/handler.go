package handler

import (
	"avitoTech/pkg/service"
	"github.com/gin-gonic/gin"
)

// Handler хранит данные из сервисов (service)
type Handler struct {
	services *service.Service
}

// NewHandler создает новый объект *Handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

// InitRoutes инициализирует routes
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	// gin.Default()

	router.Use(gin.Logger())
	//router.NoRoute(h.notFound)
	router.POST("/top-up", h.topUp)
	router.GET("/balance", h.balance)
	router.POST("/debit", h.debit)
	router.POST("/transfer", h.transfer)
	return router
}
