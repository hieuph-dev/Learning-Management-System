package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type PaymentRoutes struct {
	handler *handler.PaymentHandler
}

func NewPaymentRoutes(handler *handler.PaymentHandler) *PaymentRoutes {
	return &PaymentRoutes{
		handler: handler,
	}
}

func (pr *PaymentRoutes) Register(r *gin.RouterGroup) {
	payments := r.Group("/payments")
	{
		// Protected routes - require authentication
		payments.Use(middleware.AuthMiddleware())
		{
			// Create payment
			payments.POST("/create", pr.handler.CreatePayment)

			// Check payment status
			payments.GET("/status", pr.handler.CheckPaymentStatus)
		}

		// Public callback endpoints - no authentication needed
		// These will be called by payment gateways
		callbacks := payments.Group("/")
		callbacks.Use() // Remove auth middleware for callbacks
		{
			// MoMo IPN callback
			payments.POST("/momo/callback", pr.handler.MomoCallback)

			// ZaloPay callback
			payments.POST("/zalopay/callback", pr.handler.ZaloPayCallback)
		}
	}
}
