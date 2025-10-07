package repository

import (
	"lms/src/dto"
	"lms/src/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, bool)
	FindByUsername(username string) (*models.User, bool)
	FindById(id uint) (*models.User, error)
	UpdatePassword(userId uint, hashedPassword string) error
	UpdateProfile(userId uint, updates map[string]interface{}) error
	ChangePassword(userId uint, hashedPassword string) error
	UpdateAvatar(userId uint, avatarURL string) error
	GetUsersWithPagination(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.User, int, error)
	DeleteUser(userId uint) error
}

type PasswordResetRepository interface {
	Create(reset *models.PasswordReset) error
	FindByToken(token string) (*models.PasswordReset, error)
	MarkAsUsed(id uint) error
	DeleteExpired() error
	DeleteByEmail(email string) error
}

type CategoryRepository interface {
	GetCategories(filters map[string]interface{}) ([]models.Category, int, error)
	FindById(id uint) (*models.Category, error)
	Create(category *models.Category) error
	FindBySlug(slug string) (*models.Category, bool)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	HasChildren(id uint) (bool, error)
	FindBySlugExcept(slug string, excludeId uint) (*models.Category, bool)
}

type CourseRepository interface {
	GetCoursesWithPagination(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Course, int, error)
	SearchCourses(query string, offset, limit int, filters map[string]interface{}, sortBy, order string) ([]models.Course, int, error)
	GetSearchFilters(query string) (*dto.SearchFilters, error)
	GetFeaturedCourses(limit int, filters map[string]interface{}) ([]models.Course, int, error)
	FindBySlug(slug string) (*models.Course, error)
	FindById(courseId uint) (*models.Course, error)
	UpdateCourseStatus(courseId uint, status string) error
}

type ReviewRepository interface {
	GetCourseReviews(courseId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Review, int, error)
	GetCourseReviewStats(courseId uint) (*dto.ReviewStats, error)
	FindByUserAndCourse(userId, courseId uint) (*models.Review, error)
	Create(review *models.Review) error
	FindById(reviewId uint) (*models.Review, error)
	UpdateCourseRatingStats(courseId uint) error
	Delete(reviewId uint) error
	Update(reviewId uint, updates map[string]interface{}) error
}

type LessonRepository interface {
	GetCourseLessons(courseId uint) ([]models.Lesson, error)
	CheckUserEnrollment(userId, courseId uint) (bool, error)
	GetLessonProgress(userId uint, lessonIds []uint) (map[uint]bool, error)
	GetLessonProgressDetail(userId, lessonId uint) (*models.Progress, error)
	GetPreviousLesson(courseId uint, currentOrder int) (*models.Lesson, error)
	GetNextLesson(courseId uint, currentOrder int) (*models.Lesson, error)
	FindLessonBySlugAndCourse(slug string, courseId uint) (*models.Lesson, error)
	FindLessonByIds(lessonIds []uint) ([]models.Lesson, error)
}

type CouponRepository interface {
	FindByCode(code string) (*models.Coupon, error)
	FindById(id uint) (*models.Coupon, error)
	IncrementUsedCount(couponId uint) error
	IsValidCoupon(coupon *models.Coupon) bool
	FindByCodeExcept(code string, excludeId uint) (*models.Coupon, bool)
	Delete(couponId uint) error
	Update(couponId uint, updates map[string]interface{}) error
	Create(coupon *models.Coupon) error
	GetCouponsWithPagination(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Coupon, int, error)
}

type OrderRepository interface {
	Create(order *models.Order) error
	Update(order *models.Order) error
	FindById(orderId uint) (*models.Order, error)
	FindByOrderCode(orderCode string) (*models.Order, error)
	UpdatePaymentStatus(orderId uint, status string) error
	GetUsersOrders(userId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Order, int, error)
	FindPendingOrderByUserAndCourse(userId, courseId uint) (*models.Order, error)
	GetAllOrders(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Order, int, error)
	UpdateOrderStatus(orderId uint, status string) error
	GetOrderStatistics(filters map[string]interface{}) (*dto.OrderStatistics, error)
}

type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	CheckEnrollment(userId, courseId uint) (*models.Enrollment, bool)
	CheckUserEnrollment(userId, courseId uint) (bool, error)
	GetUserEnrollments(userId uint, offset, limit int, filters map[string]interface{}) ([]models.Enrollment, int, error)
	CompleteEnrollment(enrollmentId uint) error
	UpdateEnrollmentProgress(enrollmentId uint, updates map[string]interface{}) error
}

type InstructorRepository interface {
	GetInstructorCourses(instructorId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Course, int, error)
	CreateCourse(course *models.Course) error
	FindCourseBySlug(slug string) (*models.Course, bool)
	FindCourseById(courseId uint) (*models.Course, error)
	FindCourseByIdAndInstructor(courseId, instructorId uint) (*models.Course, error)
	UpdateCourse(courseId uint, updates map[string]interface{}) error
	DeleteCourse(courseId uint) error
	CountEnrollmentsByCourse(courseId uint) (int64, error)
	GetCourseStudents(courseId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Enrollment, int, error)
	GetStudentStatistics(courseId uint) (*dto.StudentStatistics, error)
	CreateLesson(lesson *models.Lesson) error
	FindLessonBySlug(slug string, courseId uint) (*models.Lesson, bool)
	CheckLessonOrderExists(courseId uint, lessonOrder int) (bool, error)
	FindLessonByIdAndCourse(lessonId, courseId uint) (*models.Lesson, error)
	UpdateLesson(lessonId uint, updates map[string]interface{}) error
	DeleteLesson(lessonId uint) error
	CheckLessonOrderExistsExcept(courseId uint, lessonOrder int, excludeId uint) (bool, error)
	FindLessonsByIds(lessonIds []uint) ([]models.Lesson, error)
	UpdateLessonOrder(lessonId uint, newOrder int) error
	BeginTransaction() *gorm.DB
}

type ProgressRepository interface {
	CountCompletedLessons(userId, courseId uint) (int, error)
	GetCourseProgress(userId, courseId uint) ([]models.Progress, error)
	UpdateProgress(progress *models.Progress) error
	GetLessonProgress(userId, lessonId uint) (*models.Progress, error)
}

type AnalyticsRepository interface {
	GetInstructorOverview(instructorId uint) (*dto.InstructorOverviewResponse, error)
	GetRevenueAnalytics(instructorId uint, req *dto.RevenueAnalyticsRequest) (*dto.RevenueAnalyticsResponse, error)
	GetStudentAnalytics(instructorId uint, req *dto.StudentAnalyticsRequest) (*dto.StudentAnalyticsResponse, error)
}

// Thêm interface này vào file interfaces.go
type AdminAnalyticsRepository interface {
	GetAdminDashboard() (*dto.AdminDashboardResponse, error)
	GetAdminRevenueAnalytics(req *dto.AdminRevenueAnalyticsRequest) (*dto.AdminRevenueAnalyticsResponse, error)
	GetAdminUsersAnalytics(req *dto.AdminUsersAnalyticsRequest) (*dto.AdminUsersAnalyticsResponse, error)
	GetAdminCoursesAnalytics(req *dto.AdminCoursesAnalyticsRequest) (*dto.AdminCoursesAnalyticsResponse, error)
}
