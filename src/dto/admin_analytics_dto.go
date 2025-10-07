package dto

// Admin Dashboard Response
type AdminDashboardResponse struct {
	TotalUsers        int     `json:"total_users"`
	ActiveUsers       int     `json:"active_users"`
	NewUsersThisMonth int     `json:"new_users_this_month"`
	TotalCourses      int     `json:"total_courses"`
	PublishedCourses  int     `json:"published_courses"`
	DraftCourses      int     `json:"draft_courses"`
	TotalInstructors  int     `json:"total_instructors"`
	TotalStudents     int     `json:"total_students"`
	TotalRevenue      float64 `json:"total_revenue"`
	MonthRevenue      float64 `json:"month_revenue"`
	TotalOrders       int     `json:"total_orders"`
	PendingOrders     int     `json:"pending_orders"`
	CompletedOrders   int     `json:"completed_orders"`
	TotalEnrollments  int     `json:"total_enrollments"`
	MonthEnrollments  int     `json:"month_enrollments"`
	AverageRating     float32 `json:"average_rating"`
	TotalReviews      int     `json:"total_reviews"`
}

// Admin Revenue Analytics Request
type AdminRevenueAnalyticsRequest struct {
	StartDate string `form:"start_date" binding:"omitempty"`
	EndDate   string `form:"end_date" binding:"omitempty"`
	Period    string `form:"period" binding:"omitempty,oneof=day week month year"`
}

// Admin Revenue Analytics Response
type AdminRevenueAnalyticsResponse struct {
	TotalRevenue        float64                 `json:"total_revenue"`
	TotalOrders         int                     `json:"total_orders"`
	AverageOrderValue   float64                 `json:"average_order_value"`
	RevenueByPeriod     []RevenuePeriodItem     `json:"revenue_by_period"`
	RevenueByInstructor []InstructorRevenueItem `json:"revenue_by_instructor"`
	RevenueByCategory   []CategoryRevenueItem   `json:"revenue_by_category"`
	TopCourses          []TopCourseRevenue      `json:"top_courses"`
	RevenueGrowth       float64                 `json:"revenue_growth"`
	PaymentMethodStats  []PaymentMethodStat     `json:"payment_method_stats"`
}

type InstructorRevenueItem struct {
	InstructorId   uint    `json:"instructor_id"`
	InstructorName string  `json:"instructor_name"`
	Revenue        float64 `json:"revenue"`
	Orders         int     `json:"orders"`
	Courses        int     `json:"courses"`
}

type CategoryRevenueItem struct {
	CategoryId   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Revenue      float64 `json:"revenue"`
	Orders       int     `json:"orders"`
	Courses      int     `json:"courses"`
}

type PaymentMethodStat struct {
	Method string  `json:"method"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
}

// Admin Users Analytics Request
type AdminUsersAnalyticsRequest struct {
	StartDate string `form:"start_date" binding:"omitempty"`
	EndDate   string `form:"end_date" binding:"omitempty"`
	Role      string `form:"role" binding:"omitempty,oneof=student instructor admin"`
}

// Admin Users Analytics Response
type AdminUsersAnalyticsResponse struct {
	TotalUsers        int              `json:"total_users"`
	ActiveUsers       int              `json:"active_users"`
	InactiveUsers     int              `json:"inactive_users"`
	BannedUsers       int              `json:"banned_users"`
	NewUsers          int              `json:"new_users"`
	UsersByPeriod     []UserPeriodItem `json:"users_by_period"`
	UsersByRole       []UserRoleItem   `json:"users_by_role"`
	UserGrowthRate    float64          `json:"user_growth_rate"`
	EmailVerifiedRate float64          `json:"email_verified_rate"`
}

type UserPeriodItem struct {
	Period      string `json:"period"`
	NewUsers    int    `json:"new_users"`
	ActiveUsers int    `json:"active_users"`
}

type UserRoleItem struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}

// Admin Courses Analytics Request
type AdminCoursesAnalyticsRequest struct {
	StartDate  string `form:"start_date" binding:"omitempty"`
	EndDate    string `form:"end_date" binding:"omitempty"`
	CategoryId uint   `form:"category_id" binding:"omitempty"`
	Status     string `form:"status" binding:"omitempty,oneof=draft published archived"`
}

// Admin Courses Analytics Response
type AdminCoursesAnalyticsResponse struct {
	TotalCourses        int                    `json:"total_courses"`
	PublishedCourses    int                    `json:"published_courses"`
	DraftCourses        int                    `json:"draft_courses"`
	ArchivedCourses     int                    `json:"archived_courses"`
	NewCourses          int                    `json:"new_courses"`
	CoursesByPeriod     []CoursePeriodItem     `json:"courses_by_period"`
	CoursesByCategory   []CourseCategoryItem   `json:"courses_by_category"`
	CoursesByInstructor []CourseInstructorItem `json:"courses_by_instructor"`
	AverageRating       float32                `json:"average_rating"`
	TotalEnrollments    int                    `json:"total_enrollments"`
	AverageEnrollments  float64                `json:"average_enrollments"`
	PublishRate         float64                `json:"publish_rate"`
}

type CoursePeriodItem struct {
	Period      string `json:"period"`
	NewCourses  int    `json:"new_courses"`
	Published   int    `json:"published"`
	Enrollments int    `json:"enrollments"`
}

type CourseCategoryItem struct {
	CategoryId   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Courses      int     `json:"courses"`
	Enrollments  int     `json:"enrollments"`
	AvgRating    float32 `json:"avg_rating"`
}

type CourseInstructorItem struct {
	InstructorId   uint    `json:"instructor_id"`
	InstructorName string  `json:"instructor_name"`
	Courses        int     `json:"courses"`
	Enrollments    int     `json:"enrollments"`
	AvgRating      float32 `json:"avg_rating"`
}
