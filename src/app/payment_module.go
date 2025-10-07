package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/payment"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
	"lms/src/utils"
	"strconv"
)

type PaymentModule struct {
	routes routes.Route
}

func NewPaymentModule() *PaymentModule {
	// Repositories
	orderRepo := repository.NewDBOrderRepository(db.DB)
	enrollmentRepo := repository.NewDBEnrollmentRepository(db.DB)
	// courseRepo := repository.NewDBCourseRepository(db.DB)
	couponRepo := repository.NewDBCouponRepository(db.DB)

	// Payment gateways configuration
	momoConfig := payment.MomoConfig{
		PartnerCode: utils.GetEnv("MOMO_PARTNER_CODE", ""),
		AccessKey:   utils.GetEnv("MOMO_ACCESS_KEY", ""),
		SecretKey:   utils.GetEnv("MOMO_SECRET_KEY", ""),
		Endpoint:    utils.GetEnv("MOMO_ENDPOINT", "https://test-payment.momo.vn/v2/gateway/api/create"),
		ReturnURL:   utils.GetEnv("MOMO_RETURN_URL", "http://localhost:3000/payment/result"),
		IPNUrl:      utils.GetEnv("MOMO_IPN_URL", "http://localhost:8080/api/v1/payments/momo/callback"),
	}

	appId, _ := strconv.Atoi(utils.GetEnv("ZALOPAY_APP_ID", "2553"))
	zaloPayConfig := payment.ZaloPayConfig{
		AppId:       appId,
		Key1:        utils.GetEnv("ZALOPAY_KEY1", ""),
		Key2:        utils.GetEnv("ZALOPAY_KEY2", ""),
		Endpoint:    utils.GetEnv("ZALOPAY_ENDPOINT", "https://sb-openapi.zalopay.vn/v2/create"),
		CallbackURL: utils.GetEnv("ZALOPAY_CALLBACK_URL", "http://localhost:8080/api/v1/payments/zalopay/callback"),
		ReturnURL:   utils.GetEnv("ZALOPAY_RETURN_URL", "http://localhost:3000/payment/result"),
	}

	// Initialize payment gateways
	momoGateway := payment.NewMomoPayment(momoConfig)
	zaloPayGateway := payment.NewZaloPayPayment(zaloPayConfig)

	// Payment service
	paymentService := service.NewPaymentService(
		orderRepo,
		enrollmentRepo,
		couponRepo,
		momoGateway,
		zaloPayGateway,
	)

	// Payment handler
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Payment routes
	paymentRoutes := routes.NewPaymentRoutes(paymentHandler)

	return &PaymentModule{routes: paymentRoutes}
}

func (pm *PaymentModule) Routes() routes.Route {
	return pm.routes
}
