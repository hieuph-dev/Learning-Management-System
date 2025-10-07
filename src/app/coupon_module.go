package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type CouponModule struct {
	handler *handler.CouponHandler
}

func NewCouponModule() *CouponModule {
	couponRepo := repository.NewDBCouponRepository(db.DB)
	courseRepo := repository.NewDBCourseRepository(db.DB)
	couponService := service.NewCouponService(couponRepo, courseRepo)
	couponHandler := handler.NewCouponHandler(couponService)

	return &CouponModule{
		handler: couponHandler,
	}
}

func (cm *CouponModule) Routes() routes.Route {
	return cm.handler
}

// Thêm method này để có thể inject vào AdminModule
func (cm *CouponModule) GetHandler() *handler.CouponHandler {
	return cm.handler
}
