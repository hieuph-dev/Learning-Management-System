package repository

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"time"

	"gorm.io/gorm"
)

type DBAdminAnalyticsRepository struct {
	db *gorm.DB
}

func NewDBAdminAnalyticsRepository(db *gorm.DB) AdminAnalyticsRepository {
	return &DBAdminAnalyticsRepository{db: db}
}

func (r *DBAdminAnalyticsRepository) GetAdminDashboard() (*dto.AdminDashboardResponse, error) {
	var dashboard dto.AdminDashboardResponse

	// Total users
	var totalUsers int64
	if err := r.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	dashboard.TotalUsers = int(totalUsers)

	// Active users
	var activeUsers int64
	if err := r.db.Model(&models.User{}).Where("status = ?", "active").Count(&activeUsers).Error; err != nil {
		return nil, err
	}
	dashboard.ActiveUsers = int(activeUsers)

	// New users this month
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Truncate(24 * time.Hour)
	var newUsersThisMonth int64
	if err := r.db.Model(&models.User{}).Where("created_at >= ?", startOfMonth).Count(&newUsersThisMonth).Error; err != nil {
		return nil, err
	}
	dashboard.NewUsersThisMonth = int(newUsersThisMonth)

	// Total courses
	var totalCourses int64
	if err := r.db.Model(&models.Course{}).Count(&totalCourses).Error; err != nil {
		return nil, err
	}
	dashboard.TotalCourses = int(totalCourses)

	// Published courses
	var publishedCourses int64
	if err := r.db.Model(&models.Course{}).Where("status = ?", "published").Count(&publishedCourses).Error; err != nil {
		return nil, err
	}
	dashboard.PublishedCourses = int(publishedCourses)

	// Draft courses
	var draftCourses int64
	if err := r.db.Model(&models.Course{}).Where("status = ?", "draft").Count(&draftCourses).Error; err != nil {
		return nil, err
	}
	dashboard.DraftCourses = int(draftCourses)

	// Total instructors
	var totalInstructors int64
	if err := r.db.Model(&models.User{}).Where("role = ?", "instructor").Count(&totalInstructors).Error; err != nil {
		return nil, err
	}
	dashboard.TotalInstructors = int(totalInstructors)

	// Total students
	var totalStudents int64
	if err := r.db.Model(&models.User{}).Where("role = ?", "student").Count(&totalStudents).Error; err != nil {
		return nil, err
	}
	dashboard.TotalStudents = int(totalStudents)

	// Total revenue
	var totalRevenue struct {
		Total float64
	}
	if err := r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(final_price), 0) as total").
		Where("payment_status = ?", "paid").
		Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}
	dashboard.TotalRevenue = totalRevenue.Total

	// Month revenue
	var monthRevenue struct {
		Total float64
	}
	if err := r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(final_price), 0) as total").
		Where("payment_status = ? AND paid_at >= ?", "paid", startOfMonth).
		Scan(&monthRevenue).Error; err != nil {
		return nil, err
	}
	dashboard.MonthRevenue = monthRevenue.Total

	// Total orders
	var totalOrders int64
	if err := r.db.Model(&models.Order{}).Count(&totalOrders).Error; err != nil {
		return nil, err
	}
	dashboard.TotalOrders = int(totalOrders)

	// Pending orders
	var pendingOrders int64
	if err := r.db.Model(&models.Order{}).Where("payment_status = ?", "pending").Count(&pendingOrders).Error; err != nil {
		return nil, err
	}
	dashboard.PendingOrders = int(pendingOrders)

	// Completed orders
	var completedOrders int64
	if err := r.db.Model(&models.Order{}).Where("payment_status = ?", "paid").Count(&completedOrders).Error; err != nil {
		return nil, err
	}
	dashboard.CompletedOrders = int(completedOrders)

	// Total enrollments
	var totalEnrollments int64
	if err := r.db.Model(&models.Enrollment{}).Count(&totalEnrollments).Error; err != nil {
		return nil, err
	}
	dashboard.TotalEnrollments = int(totalEnrollments)

	// Month enrollments
	var monthEnrollments int64
	if err := r.db.Model(&models.Enrollment{}).Where("enrolled_at >= ?", startOfMonth).Count(&monthEnrollments).Error; err != nil {
		return nil, err
	}
	dashboard.MonthEnrollments = int(monthEnrollments)

	// Average rating
	var avgRating struct {
		Avg float32
	}
	if err := r.db.Model(&models.Course{}).
		Select("COALESCE(AVG(rating_avg), 0) as avg").
		Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	dashboard.AverageRating = avgRating.Avg

	// Total reviews
	var totalReviews int64
	if err := r.db.Model(&models.Review{}).Count(&totalReviews).Error; err != nil {
		return nil, err
	}
	dashboard.TotalReviews = int(totalReviews)

	return &dashboard, nil
}

func (r *DBAdminAnalyticsRepository) GetAdminRevenueAnalytics(req *dto.AdminRevenueAnalyticsRequest) (*dto.AdminRevenueAnalyticsResponse, error) {
	var response dto.AdminRevenueAnalyticsResponse

	// Parse dates
	var startDate, endDate time.Time
	var err error

	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			startDate = time.Now().AddDate(0, -1, 0)
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			endDate = time.Now()
		}
	} else {
		endDate = time.Now()
	}

	// Base query
	query := r.db.Model(&models.Order{}).
		Where("payment_status = ? AND paid_at BETWEEN ? AND ?", "paid", startDate, endDate)

	// Total revenue and orders
	var stats struct {
		Total  float64
		Orders int64
	}
	if err := query.Select("COALESCE(SUM(final_price), 0) as total, COUNT(*) as orders").
		Scan(&stats).Error; err != nil {
		return nil, err
	}
	response.TotalRevenue = stats.Total
	response.TotalOrders = int(stats.Orders)

	// Average order value
	if response.TotalOrders > 0 {
		response.AverageOrderValue = response.TotalRevenue / float64(response.TotalOrders)
	}

	// Revenue by period
	period := req.Period
	if period == "" {
		period = "day"
	}

	var periodFormat string
	switch period {
	case "day":
		periodFormat = "YYYY-MM-DD"
	case "week":
		periodFormat = "IYYY-IW"
	case "month":
		periodFormat = "YYYY-MM"
	case "year":
		periodFormat = "YYYY"
	default:
		periodFormat = "YYYY-MM-DD"
	}

	var revenueByPeriod []struct {
		Period   string
		Revenue  float64
		Orders   int
		Students int
	}

	if err := r.db.Raw(`
		SELECT 
			TO_CHAR(paid_at, ?) as period,
			COALESCE(SUM(final_price), 0) as revenue,
			COUNT(id) as orders,
			COUNT(DISTINCT user_id) as students
		FROM orders
		WHERE payment_status = ? AND paid_at BETWEEN ? AND ?
		GROUP BY period
		ORDER BY period
	`, periodFormat, "paid", startDate, endDate).Scan(&revenueByPeriod).Error; err != nil {
		return nil, err
	}

	for _, item := range revenueByPeriod {
		response.RevenueByPeriod = append(response.RevenueByPeriod, dto.RevenuePeriodItem{
			Period:   item.Period,
			Revenue:  item.Revenue,
			Orders:   item.Orders,
			Students: item.Students,
		})
	}

	// Revenue by instructor
	var revenueByInstructor []struct {
		InstructorId   uint
		InstructorName string
		Revenue        float64
		Orders         int
		Courses        int
	}

	if err := r.db.Raw(`
		SELECT 
			users.id as instructor_id,
			users.full_name as instructor_name,
			COALESCE(SUM(orders.final_price), 0) as revenue,
			COUNT(DISTINCT orders.id) as orders,
			COUNT(DISTINCT courses.id) as courses
		FROM users
		LEFT JOIN courses ON courses.instructor_id = users.id
		LEFT JOIN orders ON orders.course_id = courses.id 
			AND orders.payment_status = ?
			AND orders.paid_at BETWEEN ? AND ?
		WHERE users.role = ?
		GROUP BY users.id, users.full_name
		ORDER BY revenue DESC
		LIMIT 10
	`, "paid", startDate, endDate, "instructor").Scan(&revenueByInstructor).Error; err != nil {
		return nil, err
	}

	for _, item := range revenueByInstructor {
		response.RevenueByInstructor = append(response.RevenueByInstructor, dto.InstructorRevenueItem{
			InstructorId:   item.InstructorId,
			InstructorName: item.InstructorName,
			Revenue:        item.Revenue,
			Orders:         item.Orders,
			Courses:        item.Courses,
		})
	}

	// Revenue by category
	var revenueByCategory []struct {
		CategoryId   uint
		CategoryName string
		Revenue      float64
		Orders       int
		Courses      int
	}

	if err := r.db.Raw(`
		SELECT 
			categories.id as category_id,
			categories.name as category_name,
			COALESCE(SUM(orders.final_price), 0) as revenue,
			COUNT(DISTINCT orders.id) as orders,
			COUNT(DISTINCT courses.id) as courses
		FROM categories
		LEFT JOIN courses ON courses.category_id = categories.id
		LEFT JOIN orders ON orders.course_id = courses.id 
			AND orders.payment_status = ?
			AND orders.paid_at BETWEEN ? AND ?
		GROUP BY categories.id, categories.name
		ORDER BY revenue DESC
	`, "paid", startDate, endDate).Scan(&revenueByCategory).Error; err != nil {
		return nil, err
	}

	for _, item := range revenueByCategory {
		response.RevenueByCategory = append(response.RevenueByCategory, dto.CategoryRevenueItem{
			CategoryId:   item.CategoryId,
			CategoryName: item.CategoryName,
			Revenue:      item.Revenue,
			Orders:       item.Orders,
			Courses:      item.Courses,
		})
	}

	// Top courses
	var topCourses []struct {
		CourseId    uint
		CourseTitle string
		Revenue     float64
		Orders      int
		Students    int
	}

	if err := r.db.Raw(`
		SELECT 
			courses.id as course_id,
			courses.title as course_title,
			COALESCE(SUM(orders.final_price), 0) as revenue,
			COUNT(orders.id) as orders,
			COUNT(DISTINCT orders.user_id) as students
		FROM courses
		LEFT JOIN orders ON orders.course_id = courses.id 
			AND orders.payment_status = ?
			AND orders.paid_at BETWEEN ? AND ?
		GROUP BY courses.id, courses.title
		ORDER BY revenue DESC
		LIMIT 10
	`, "paid", startDate, endDate).Scan(&topCourses).Error; err != nil {
		return nil, err
	}

	for _, item := range topCourses {
		response.TopCourses = append(response.TopCourses, dto.TopCourseRevenue{
			CourseId:    item.CourseId,
			CourseTitle: item.CourseTitle,
			Revenue:     item.Revenue,
			Orders:      item.Orders,
			Students:    item.Students,
		})
	}

	// Payment method stats
	var paymentStats []struct {
		Method string
		Count  int
		Amount float64
	}

	if err := r.db.Raw(`
		SELECT 
			payment_method as method,
			COUNT(*) as count,
			COALESCE(SUM(final_price), 0) as amount
		FROM orders
		WHERE payment_status = ? AND paid_at BETWEEN ? AND ?
		GROUP BY payment_method
		ORDER BY amount DESC
	`, "paid", startDate, endDate).Scan(&paymentStats).Error; err != nil {
		return nil, err
	}

	for _, item := range paymentStats {
		response.PaymentMethodStats = append(response.PaymentMethodStats, dto.PaymentMethodStat{
			Method: item.Method,
			Count:  item.Count,
			Amount: item.Amount,
		})
	}

	// Revenue growth
	periodDiff := endDate.Sub(startDate)
	previousStartDate := startDate.Add(-periodDiff)
	previousEndDate := startDate

	var previousRevenue struct {
		Total float64
	}
	if err := r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(final_price), 0) as total").
		Where("payment_status = ? AND paid_at BETWEEN ? AND ?", "paid", previousStartDate, previousEndDate).
		Scan(&previousRevenue).Error; err != nil {
		return nil, err
	}

	if previousRevenue.Total > 0 {
		response.RevenueGrowth = ((response.TotalRevenue - previousRevenue.Total) / previousRevenue.Total) * 100
	}

	return &response, nil
}

func (r *DBAdminAnalyticsRepository) GetAdminUsersAnalytics(req *dto.AdminUsersAnalyticsRequest) (*dto.AdminUsersAnalyticsResponse, error) {
	var response dto.AdminUsersAnalyticsResponse

	// Parse dates
	var startDate, endDate time.Time
	var err error

	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			startDate = time.Now().AddDate(0, -1, 0)
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			endDate = time.Now()
		}
	} else {
		endDate = time.Now()
	}

	// Base query
	baseQuery := r.db.Model(&models.User{})
	if req.Role != "" {
		baseQuery = baseQuery.Where("role = ?", req.Role)
	}

	// Total users
	var totalUsers int64
	if err := baseQuery.Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	response.TotalUsers = int(totalUsers)

	// Active users
	var activeUsers int64
	if err := baseQuery.Where("status = ?", "active").Count(&activeUsers).Error; err != nil {
		return nil, err
	}
	response.ActiveUsers = int(activeUsers)

	// Inactive users
	var inactiveUsers int64
	if err := baseQuery.Where("status = ?", "inactive").Count(&inactiveUsers).Error; err != nil {
		return nil, err
	}
	response.InactiveUsers = int(inactiveUsers)

	// Banned users
	var bannedUsers int64
	if err := baseQuery.Where("status = ?", "banned").Count(&bannedUsers).Error; err != nil {
		return nil, err
	}
	response.BannedUsers = int(bannedUsers)

	// New users in period
	newQuery := r.db.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)
	if req.Role != "" {
		newQuery = newQuery.Where("role = ?", req.Role)
	}
	var newUsers int64
	if err := newQuery.Count(&newUsers).Error; err != nil {
		return nil, err
	}
	response.NewUsers = int(newUsers)

	// Users by period
	var usersByPeriod []struct {
		Period      string
		NewUsers    int
		ActiveUsers int
	}

	roleFilter := ""
	if req.Role != "" {
		roleFilter = " AND role = '" + req.Role + "'"
	}

	if err := r.db.Raw(`
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM-DD') as period,
			COUNT(*) as new_users,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_users
		FROM users
		WHERE created_at BETWEEN ? AND ?`+roleFilter+`
		GROUP BY period
		ORDER BY period
	`, startDate, endDate).Scan(&usersByPeriod).Error; err != nil {
		return nil, err
	}

	for _, item := range usersByPeriod {
		response.UsersByPeriod = append(response.UsersByPeriod, dto.UserPeriodItem{
			Period:      item.Period,
			NewUsers:    item.NewUsers,
			ActiveUsers: item.ActiveUsers,
		})
	}

	// Users by role
	var usersByRole []struct {
		Role  string
		Count int
	}

	if err := r.db.Raw(`
		SELECT role, COUNT(*) as count
		FROM users
		GROUP BY role
		ORDER BY count DESC
	`).Scan(&usersByRole).Error; err != nil {
		return nil, err
	}

	for _, item := range usersByRole {
		response.UsersByRole = append(response.UsersByRole, dto.UserRoleItem{
			Role:  item.Role,
			Count: item.Count,
		})
	}

	// User growth rate
	// periodDiff := endDate.Sub(startDate)
	// previousStartDate := startDate.Add(-periodDiff)

	var previousUsers int64
	previousQuery := r.db.Model(&models.User{}).Where("created_at < ?", startDate)
	if req.Role != "" {
		previousQuery = previousQuery.Where("role = ?", req.Role)
	}
	if err := previousQuery.Count(&previousUsers).Error; err != nil {
		return nil, err
	}

	if previousUsers > 0 {
		response.UserGrowthRate = (float64(newUsers) / float64(previousUsers)) * 100
	}

	// Email verified rate
	var verifiedUsers int64
	verifiedQuery := r.db.Model(&models.User{}).Where("email_verified = ?", true)
	if req.Role != "" {
		verifiedQuery = verifiedQuery.Where("role = ?", req.Role)
	}
	if err := verifiedQuery.Count(&verifiedUsers).Error; err != nil {
		return nil, err
	}

	if totalUsers > 0 {
		response.EmailVerifiedRate = (float64(verifiedUsers) / float64(totalUsers)) * 100
	}

	return &response, nil
}

func (r *DBAdminAnalyticsRepository) GetAdminCoursesAnalytics(req *dto.AdminCoursesAnalyticsRequest) (*dto.AdminCoursesAnalyticsResponse, error) {
	var response dto.AdminCoursesAnalyticsResponse

	// Parse dates
	var startDate, endDate time.Time
	var err error

	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			startDate = time.Now().AddDate(0, -1, 0)
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			endDate = time.Now()
		}
	} else {
		endDate = time.Now()
	}

	// Base query
	baseQuery := r.db.Model(&models.Course{})
	if req.CategoryId != 0 {
		baseQuery = baseQuery.Where("category_id = ?", req.CategoryId)
	}
	if req.Status != "" {
		baseQuery = baseQuery.Where("status = ?", req.Status)
	}

	// Total courses
	var totalCourses int64
	if err := baseQuery.Count(&totalCourses).Error; err != nil {
		return nil, err
	}
	response.TotalCourses = int(totalCourses)

	// Published courses
	var publishedCourses int64
	pubQuery := r.db.Model(&models.Course{}).Where("status = ?", "published")
	if req.CategoryId != 0 {
		pubQuery = pubQuery.Where("category_id = ?", req.CategoryId)
	}
	if err := pubQuery.Count(&publishedCourses).Error; err != nil {
		return nil, err
	}
	response.PublishedCourses = int(publishedCourses)

	// Draft courses
	var draftCourses int64
	draftQuery := r.db.Model(&models.Course{}).Where("status = ?", "draft")
	if req.CategoryId != 0 {
		draftQuery = draftQuery.Where("category_id = ?", req.CategoryId)
	}
	if err := draftQuery.Count(&draftCourses).Error; err != nil {
		return nil, err
	}
	response.DraftCourses = int(draftCourses)

	// Archived courses
	var archivedCourses int64
	archQuery := r.db.Model(&models.Course{}).Where("status = ?", "archived")
	if req.CategoryId != 0 {
		archQuery = archQuery.Where("category_id = ?", req.CategoryId)
	}
	if err := archQuery.Count(&archivedCourses).Error; err != nil {
		return nil, err
	}
	response.ArchivedCourses = int(archivedCourses)

	// New courses in period
	newQuery := r.db.Model(&models.Course{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)
	if req.CategoryId != 0 {
		newQuery = newQuery.Where("category_id = ?", req.CategoryId)
	}
	if req.Status != "" {
		newQuery = newQuery.Where("status = ?", req.Status)
	}

	var newCourses int64
	if err := newQuery.Count(&newCourses).Error; err != nil {
		return nil, err
	}
	response.NewCourses = int(newCourses)

	// Courses by period
	var coursesByPeriod []struct {
		Period      string
		NewCourses  int
		Published   int
		Enrollments int
	}

	categoryFilter := ""
	if req.CategoryId != 0 {
		categoryFilter = " AND category_id = " + fmt.Sprintf("%d", req.CategoryId)
	}

	statusFilter := ""
	if req.Status != "" {
		statusFilter = " AND status = '" + req.Status + "'"
	}

	if err := r.db.Raw(`
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM-DD') as period,
			COUNT(*) as new_courses,
			COUNT(CASE WHEN status = 'published' THEN 1 END) as published,
			COALESCE(SUM(enrolled_count), 0) as enrollments
		FROM courses
		WHERE created_at BETWEEN ? AND ?`+categoryFilter+statusFilter+`
		GROUP BY period
		ORDER BY period
	`, startDate, endDate).Scan(&coursesByPeriod).Error; err != nil {
		return nil, err
	}

	for _, item := range coursesByPeriod {
		response.CoursesByPeriod = append(response.CoursesByPeriod, dto.CoursePeriodItem{
			Period:      item.Period,
			NewCourses:  item.NewCourses,
			Published:   item.Published,
			Enrollments: item.Enrollments,
		})
	}

	// Courses by category
	var coursesByCategory []struct {
		CategoryId   uint
		CategoryName string
		Courses      int
		Enrollments  int
		AvgRating    float32
	}

	if err := r.db.Raw(`
		SELECT 
			categories.id as category_id,
			categories.name as category_name,
			COUNT(courses.id) as courses,
			COALESCE(SUM(courses.enrolled_count), 0) as enrollments,
			COALESCE(AVG(courses.rating_avg), 0) as avg_rating
		FROM categories
		LEFT JOIN courses ON courses.category_id = categories.id
		GROUP BY categories.id, categories.name
		ORDER BY courses DESC
	`).Scan(&coursesByCategory).Error; err != nil {
		return nil, err
	}

	for _, item := range coursesByCategory {
		response.CoursesByCategory = append(response.CoursesByCategory, dto.CourseCategoryItem{
			CategoryId:   item.CategoryId,
			CategoryName: item.CategoryName,
			Courses:      item.Courses,
			Enrollments:  item.Enrollments,
			AvgRating:    item.AvgRating,
		})
	}

	// Courses by instructor
	var coursesByInstructor []struct {
		InstructorId   uint
		InstructorName string
		Courses        int
		Enrollments    int
		AvgRating      float32
	}

	if err := r.db.Raw(`
		SELECT 
			users.id as instructor_id,
			users.full_name as instructor_name,
			COUNT(courses.id) as courses,
			COALESCE(SUM(courses.enrolled_count), 0) as enrollments,
			COALESCE(AVG(courses.rating_avg), 0) as avg_rating
		FROM users
		LEFT JOIN courses ON courses.instructor_id = users.id
		WHERE users.role = ?
		GROUP BY users.id, users.full_name
		ORDER BY courses DESC
		LIMIT 10
	`, "instructor").Scan(&coursesByInstructor).Error; err != nil {
		return nil, err
	}

	for _, item := range coursesByInstructor {
		response.CoursesByInstructor = append(response.CoursesByInstructor, dto.CourseInstructorItem{
			InstructorId:   item.InstructorId,
			InstructorName: item.InstructorName,
			Courses:        item.Courses,
			Enrollments:    item.Enrollments,
			AvgRating:      item.AvgRating,
		})
	}

	// Average rating
	var avgRating struct {
		Avg float32
	}
	avgQuery := r.db.Model(&models.Course{})
	if req.CategoryId != 0 {
		avgQuery = avgQuery.Where("category_id = ?", req.CategoryId)
	}
	if req.Status != "" {
		avgQuery = avgQuery.Where("status = ?", req.Status)
	}

	if err := avgQuery.Select("COALESCE(AVG(rating_avg), 0) as avg").Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	response.AverageRating = avgRating.Avg

	// Total enrollments
	var totalEnrollments int64
	enrollQuery := r.db.Model(&models.Enrollment{})
	if req.CategoryId != 0 {
		enrollQuery = enrollQuery.Joins("JOIN courses ON courses.id = enrollments.course_id").
			Where("courses.category_id = ?", req.CategoryId)
	}
	if err := enrollQuery.Count(&totalEnrollments).Error; err != nil {
		return nil, err
	}
	response.TotalEnrollments = int(totalEnrollments)

	// Average enrollments per course
	if response.TotalCourses > 0 {
		response.AverageEnrollments = float64(response.TotalEnrollments) / float64(response.TotalCourses)
	}

	// Publish rate
	if response.TotalCourses > 0 {
		response.PublishRate = float64(response.PublishedCourses) / float64(response.TotalCourses) * 100
	}

	return &response, nil
}
