package app

import (
	"lms/src/config"
	"lms/src/db"
	"lms/src/routes"
	"lms/src/validation"
	"log"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *config.ServerConfig
	router  *gin.Engine
	modules []Module // Ds các module
}

func NewApplication(cfg *config.ServerConfig) *Application {
	// Khởi tạo validation
	if err := validation.InitValidation(); err != nil {
		log.Fatalf("Validator init failed %v", err)
	}

	// Load biến môi trường
	config.LoadEnv()

	// Kết nối DB
	err := db.InitDB()
	if err != nil {
		log.Fatal("unable to connect to db")
	}

	// Tạo Gin router
	r := gin.Default()

	// Định nghĩa các module
	modules := []Module{
		NewAuthModule(),
		NewUserModule(),
		NewAdminModule(),
		NewCategoryModule(),
		NewCourseModule(),
		NewLessonModule(),
		NewEnrollmentModule(),
		NewInstructorModule(),
		NewProgressModule(),
		NewOrderModule(),
		NewCouponModule(),
	}

	// Đăng ký routes cho tất cả modules
	routes.RegisterRoutes(r, getModuleRoutes(modules)...)

	// Trả về Application instance
	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}
}

func (a *Application) Run() error { // a chính là &Application{config: cfg, router: r,}
	return a.router.Run(a.config.ServerAddress) // Hàm Run này là của Gin
}

func getModuleRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
