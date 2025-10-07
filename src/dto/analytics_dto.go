package dto

// Overview Analytics Response
type InstructorOverviewResponse struct {
	TotalCourses     int     `json:"total_courses"`
	PublishedCourses int     `json:"published_courses"`
	DraftCourses     int     `json:"draft_courses"`
	TotalStudents    int     `json:"total_students"`
	ActiveStudents   int     `json:"active_students"`
	TotalRevenue     float64 `json:"total_revenue"`
	MonthRevenue     float64 `json:"month_revenue"`
	TotalEnrollments int     `json:"total_enrollments"`
	MonthEnrollments int     `json:"month_enrollments"`
	AverageRating    float32 `json:"average_rating"`
	TotalReviews     int     `json:"total_reviews"`
	CompletionRate   float64 `json:"completion_rate"`
}

// Revenue Analytics Request
type RevenueAnalyticsRequest struct {
	StartDate string `form:"start_date" binding:"omitempty"`
	EndDate   string `form:"end_date" binding:"omitempty"`
	Period    string `form:"period" binding:"omitempty,oneof=day week month year"`
	CourseId  uint   `form:"course_id" binding:"omitempty"`
}

// Revenue Analytics Response
type RevenueAnalyticsResponse struct {
	TotalRevenue      float64             `json:"total_revenue"`
	TotalOrders       int                 `json:"total_orders"`
	AverageOrderValue float64             `json:"average_order_value"`
	RevenueByPeriod   []RevenuePeriodItem `json:"revenue_by_period"`
	TopCourses        []TopCourseRevenue  `json:"top_courses"`
	RevenueGrowth     float64             `json:"revenue_growth"`
}

type RevenuePeriodItem struct {
	Period   string  `json:"period"`
	Revenue  float64 `json:"revenue"`
	Orders   int     `json:"orders"`
	Students int     `json:"students"`
}

type TopCourseRevenue struct {
	CourseId    uint    `json:"course_id"`
	CourseTitle string  `json:"course_title"`
	Revenue     float64 `json:"revenue"`
	Orders      int     `json:"orders"`
	Students    int     `json:"students"`
}

// Student Analytics Request
type StudentAnalyticsRequest struct {
	StartDate string `form:"start_date" binding:"omitempty"`
	EndDate   string `form:"end_date" binding:"omitempty"`
	CourseId  uint   `form:"course_id" binding:"omitempty"`
	Status    string `form:"status" binding:"omitempty,oneof=active completed dropped"`
}

// Student Analytics Response
type StudentAnalyticsResponse struct {
	TotalStudents     int                 `json:"total_students"`
	ActiveStudents    int                 `json:"active_students"`
	CompletedStudents int                 `json:"completed_students"`
	DroppedStudents   int                 `json:"dropped_students"`
	NewStudents       int                 `json:"new_students"`
	StudentsByPeriod  []StudentPeriodItem `json:"students_by_period"`
	StudentsByCourse  []StudentCourseItem `json:"students_by_course"`
	AverageProgress   float64             `json:"average_progress"`
	CompletionRate    float64             `json:"completion_rate"`
	RetentionRate     float64             `json:"retention_rate"`
}

type StudentPeriodItem struct {
	Period         string `json:"period"`
	NewStudents    int    `json:"new_students"`
	ActiveStudents int    `json:"active_students"`
	Completed      int    `json:"completed"`
}

type StudentCourseItem struct {
	CourseId          uint    `json:"course_id"`
	CourseTitle       string  `json:"course_title"`
	TotalStudents     int     `json:"total_students"`
	ActiveStudents    int     `json:"active_students"`
	CompletedStudents int     `json:"completed_students"`
	AverageProgress   float64 `json:"average_progress"`
	CompletionRate    float64 `json:"completion_rate"`
}
