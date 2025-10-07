package service

import (
	"lms/src/dto"
	"lms/src/models"
	"mime/multipart"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	GetProfile(userId uint) (*dto.UserProfile, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(req *dto.ResetPasswordRequest) error
}

// Interface cho EmailService
type EmailService interface {
	SendPasswordResetEmail(email, resetToken, resetCode string) error
	SendWelcomeEmail(email, fullName string) error
}

type UserService interface {
	GetProfile(userId uint) (*dto.UserProfile, error)
	UpdateProfile(userId uint, req *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	ChangePassword(userId uint, req *dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error)
	UploadAvatar(userId uint, file *multipart.FileHeader) (*dto.UploadAvatarResponse, error)
}

type AdminService interface {
	GetUsers(req *dto.GetUsersQueryRequest) (*dto.GetUsersResponse, error)
	GetUserById(userId uint) (*dto.AdminUserDetail, error)
	UpdateUser(userId uint, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	DeleteUser(userId uint) (*dto.DeleteUserResponse, error)
	ChangeUserStatus(userId uint, req *dto.ChangeUserStatusRequest) (*dto.ChangeUserStatusResponse, error)
	GetCourses(req *dto.GetAdminCoursesQueryRequest) (*dto.GetAdminCoursesResponse, error)
	ChangeCourseStatus(courseId uint, req *dto.ChangeCourseStatusRequest) (*dto.ChangeCourseStatusResponse, error)
}

type CategoryService interface {
	GetCategories(req *dto.GetCategoriesQueryRequest) (*dto.GetCategoriesResponse, error)
	GetCategoryById(categoryId uint) (*dto.CategoryDetail, error)
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CreateCategoryResponse, error)
	UpdateCategory(categoryId uint, req *dto.UpdateCategoryRequest) (*dto.UpdateCategoryResponse, error)
	DeleteCategory(categoryId uint) (*dto.DeleteCategoryResponse, error)
}

type CourseService interface {
	GetCourses(req *dto.GetCoursesQueryRequest) (*dto.GetCoursesResponse, error)
	SearchCourses(req *dto.SearchCoursesQueryRequest) (*dto.SearchCoursesResponse, error)
	GetFeaturedCourses(req *dto.GetFeaturedCoursesQueryRequest) (*dto.GetFeaturedCoursesResponse, error)
	GetCourseBySlug(slug string) (*dto.CourseDetail, error)
}

type ReviewService interface {
	GetCourseReviews(courseId uint, req *dto.GetCourseReviewsQueryRequest) (*dto.GetCourseReviewsResponse, error)
	CreateReview(userId, courseId uint, req *dto.CreateReviewRequest) (*dto.CreateReviewResponse, error)
	UpdateReview(userId, reviewId uint, req *dto.UpdateReviewRequest) (*dto.UpdateReviewResponse, error)
	DeleteReview(userId, reviewId uint) (*dto.DeleteReviewResponse, error)
}

type LessonService interface {
	GetCourseLessons(userId, courseId uint) (*dto.GetCourseLessonsResponse, error)
	GetLessonDetail(userId, courseId uint, slug string) (*dto.LessonDetail, error)
}

type EnrollmentService interface {
	EnrollCourse(userId, courseId uint, req *dto.EnrollCourseRequest) (*dto.EnrollCourseResponse, error)
	CheckEnrollment(userId, courseId uint) (*dto.CheckEnrollmentResponse, error)
	GetMyEnrollments(userId uint, req *dto.GetMyEnrollmentsQueryRequest) (*dto.GetMyEnrollmentsResponse, error)
}

type InstructorService interface {
	CreateCourse(instructorId uint, req *dto.CreateCourseRequest) (*dto.CreateCourseResponse, error)
	GetInstructorCourses(instructorId uint, req *dto.GetInstructorCoursesQueryRequest) (*dto.GetInstructorCoursesResponse, error)
	UpdateCourse(instructorId, courseId uint, req *dto.UpdateCourseRequest) (*dto.UpdateCourseResponse, error)
	DeleteCourse(instructorId, courseId uint) (*dto.DeleteCourseResponse, error)
	GetCourseStudents(instructorId, courseId uint, req *dto.GetCourseStudentsQueryRequest) (*dto.GetCourseStudentsResponse, error)
	CreateLesson(instructorId, courseId uint, req *dto.CreateLessonRequest) (*dto.CreateLessonResponse, error)
	UpdateLesson(instructorId, courseId, lessonId uint, req *dto.UpdateLessonRequest) (*dto.UpdateLessonResponse, error)
	DeleteLesson(instructorId, courseId, lessonId uint) (*dto.DeleteLessonResponse, error)
	ReorderLessons(instructorId, lessonId uint, req *dto.ReorderLessonsRequest) (*dto.ReorderLessonsResponse, error)
}

type ProgressService interface {
	GetCourseProgress(userId, courseId uint) (*dto.GetCourseProgressResponse, error)
	CompleteLesson(userId, lessonId uint, req *dto.CompleteLessonRequest) (*dto.CompleteLessonResponse, error)
	UpdateLessonPosition(userId, lessonId uint, req *dto.UpdateLessonPositionRequest) (*dto.UpdateLessonPositionResponse, error)
	updateEnrollmentProgress(userId, courseId uint) error
}

type OrderService interface {
	CreateOrder(userId uint, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	GetOrderHistory(userId uint, req *dto.GetOrderHistoryQueryRequest) (*dto.GetOrderHistoryResponse, error)
	completeOrder(order *models.Order, paymentMethod string) error
	GetOrderDetail(userId uint, orderId uint) (*dto.OrderDetailResponse, error)
	PayOrder(userId uint, orderId uint, req *dto.PayOrderRequest) (*dto.PayOrderResponse, error)
	UpdateOrderStatus(orderId uint, req *dto.UpdateOrderStatusRequest) (*dto.UpdateOrderStatusResponse, error)
	GetAllOrders(req *dto.GetAdminOrdersQueryRequest) (*dto.GetAdminOrdersResponse, error)
}

type CouponService interface {
	ValidateCoupon(req *dto.ValidateCouponRequest) (*dto.ValidateCouponResponse, error)
	GetAdminCoupons(req *dto.GetAdminCouponsQueryRequest) (*dto.GetAdminCouponsResponse, error)
	CheckCoupon(req *dto.CheckCouponRequest) (*dto.CheckCouponResponse, error)
	CreateCoupon(req *dto.CreateCouponRequest) (*dto.CreateCouponResponse, error)
	DeleteCoupon(couponId uint) (*dto.DeleteCouponResponse, error)
	UpdateCoupon(couponId uint, req *dto.UpdateCouponRequest) (*dto.UpdateCouponResponse, error)
}

type AnalyticsService interface {
	GetInstructorOverview(instructorId uint) (*dto.InstructorOverviewResponse, error)
	GetRevenueAnalytics(instructorId uint, req *dto.RevenueAnalyticsRequest) (*dto.RevenueAnalyticsResponse, error)
	GetStudentAnalytics(instructorId uint, req *dto.StudentAnalyticsRequest) (*dto.StudentAnalyticsResponse, error)
}

type AdminAnalyticsService interface {
	GetAdminDashboard() (*dto.AdminDashboardResponse, error)
	GetAdminRevenueAnalytics(req *dto.AdminRevenueAnalyticsRequest) (*dto.AdminRevenueAnalyticsResponse, error)
	GetAdminUsersAnalytics(req *dto.AdminUsersAnalyticsRequest) (*dto.AdminUsersAnalyticsResponse, error)
	GetAdminCoursesAnalytics(req *dto.AdminCoursesAnalyticsRequest) (*dto.AdminCoursesAnalyticsResponse, error)
}

type PaymentService interface {
	CreatePayment(userId uint, req *dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error)
	HandleMomoCallback(data map[string]interface{}) (*dto.PaymentCallbackResponse, error)
	HandleZaloPayCallback(data map[string]interface{}) (*dto.PaymentCallbackResponse, error)
	CheckPaymentStatus(userId uint, req *dto.CheckPaymentStatusRequest) (*dto.CheckPaymentStatusResponse, error)
}
