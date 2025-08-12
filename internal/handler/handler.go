package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/realdanielursul/order-service/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.StaticFile("/", "./web/index.html")

	router.GET("/order/:order_uid", func(c *gin.Context) {
		orderUID := c.Param("order_uid")
		if orderUID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "empty order_uid param"})
			return
		}

		order, err := h.services.GetOrder(c.Request.Context(), orderUID)
		if err != nil {
			if err.Error() == "order not found" {
				c.JSON(http.StatusNotFound, gin.H{"reason": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	return router
}
