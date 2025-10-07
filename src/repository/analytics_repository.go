package repository

import (
	"lms/src/dto"
	"lms/src/models"
	"time"

	"gorm.io/gorm"
)

type DBAnalyticsRepository struct {
	db *gorm.DB
}

func NewDBAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &DBAnalyticsRepository{
		db: db,
	}
}

func (ar *DBAnalyticsRepository) GetInstructorOverview(instructorId uint) (*dto.InstructorOverviewResponse, error) {
	var overview dto.InstructorOverviewResponse

	// Total courses
	var totalCourses int64
	if err := ar.db.Model(&models.Course{}).
		Where("instructor_id = ?", instructorId).
		Count(&totalCourses).Error; err != nil {
		return nil, err
	}
	overview.TotalCourses = int(totalCourses)

	// Published course
	var publishedCourses int64
	if err := ar.db.Model(&models.Course{}).
		Where("instructor_id = ? AND status = ?", instructorId, "published").
		Count(&publishedCourses).Error; err != nil {
		return nil, err
	}
	overview.PublishedCourses = int(publishedCourses)

	// Draft courses
	var draftCourses int64
	if err := ar.db.Model(&models.Course{}).
		Where("instructor_id = ? AND status = ?", instructorId, "draft").
		Count(&draftCourses).Error; err != nil {
		return nil, err
	}
	overview.DraftCourses = int(draftCourses)

	// Total students (distinct users enrolled in instructor's courses)
	var totalStudents int64
	if err := ar.db.Model(&models.Enrollment{}).
		Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("courses.instructor_id = ?", instructorId).
		Distinct("enrollments.user_id").
		Count(&totalStudents).Error; err != nil {
		return nil, err
	}
	overview.TotalStudents = int(totalStudents)

	// Active students
	var activeStudents int64
	if err := ar.db.Model(&models.Enrollment{}).
		Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("courses.instructor_id = ? AND enrollments.status = ?", instructorId, "active").
		Distinct("enrollments.user_id").
		Count(&activeStudents).Error; err != nil {
		return nil, err
	}
	overview.ActiveStudents = int(activeStudents)

	// Total revenue
	var totalRevenue struct {
		Total float64
	}
	if err := ar.db.Model(&models.Order{}).
		Select("COALESCE(SUM(final_price), 0) as total").
		Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ? AND orders.payment_status = ?", instructorId, "paid").
		Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}
	overview.TotalRevenue = totalRevenue.Total

	// Month revenue
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Truncate(24 * time.Hour)
	var monthRevenue struct {
		Total float64
	}
	if err := ar.db.Model(&models.Order{}).
		Select("COALESCE(SUM(final_price), 0) as total").
		Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ? AND orders.payment_status = ? AND orders.paid_at >= ?",
			instructorId, "paid", startOfMonth).
		Scan(&monthRevenue).Error; err != nil {
		return nil, err
	}
	overview.MonthRevenue = monthRevenue.Total

	// Total enrollments
	var totalEnrollments int64
	if err := ar.db.Model(&models.Enrollment{}).
		Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("courses.instructor_id = ?", instructorId).
		Count(&totalEnrollments).Error; err != nil {
		return nil, err
	}
	overview.TotalEnrollments = int(totalEnrollments)

	// Month enrollments
	var monthEnrollments int64
	if err := ar.db.Model(&models.Enrollment{}).
		Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("courses.instructor_id = ? AND enrollments.enrolled_at >= ?", instructorId, startOfMonth).
		Count(&monthEnrollments).Error; err != nil {
		return nil, err
	}
	overview.MonthEnrollments = int(monthEnrollments)

	// Average rating
	var avgRating struct {
		Avg float32
	}
	if err := ar.db.Model(&models.Course{}).
		Select("COALESCE(AVG(rating_avg), 0) as avg").
		Where("instructor_id = ?", instructorId).
		Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	overview.AverageRating = avgRating.Avg

	// Total reviews
	var totalReviews int64
	if err := ar.db.Model(&models.Review{}).
		Joins("JOIN courses ON courses.id = reviews.course_id").
		Where("courses.instructor_id = ?", instructorId).
		Count(&totalReviews).Error; err != nil {
		return nil, err
	}
	overview.TotalReviews = int(totalReviews)

	// Completion rate
	var completedEnrollments int64
	if err := ar.db.Model(&models.Enrollment{}).
		Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("courses.instructor_id = ? AND enrollments.status = ?", instructorId, "completed").
		Count(&completedEnrollments).Error; err != nil {
		return nil, err
	}

	if totalEnrollments > 0 {
		overview.CompletionRate = float64(completedEnrollments) / float64(totalEnrollments) * 100
	}

	return &overview, nil
}

func (ar *DBAnalyticsRepository) GetRevenueAnalytics(instructorId uint, req *dto.RevenueAnalyticsRequest) (*dto.RevenueAnalyticsResponse, error) {
	var response dto.RevenueAnalyticsResponse

	// Parse dates
	var startDate, endDate time.Time
	var err error

	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			startDate = time.Now().AddDate(0, -1, 0) // Default: 1 month ago
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
	query := ar.db.Model(&models.Order{}).
		Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ? AND orders.payment_status = ?", instructorId, "paid")

	if req.CourseId != 0 {
		query = query.Where("orders.course_id = ?", req.CourseId)
	}

	query = query.Where("orders.paid_at BETWEEN ? AND ?", startDate, endDate)

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

	if err := ar.db.Raw(`
		SELECT 
			TO_CHAR(orders.paid_at, ?) as period,
			COALESCE(SUM(orders.final_price), 0) as revenue,
			COUNT(orders.id) as orders,
			COUNT(DISTINCT orders.user_id) as students
		FROM orders
		JOIN courses ON courses.id = orders.course_id
		WHERE courses.instructor_id = ? 
			AND orders.payment_status = ?
			AND orders.paid_at BETWEEN ? AND ?
			AND (? = 0 OR orders.course_id = ?)
		GROUP BY period
		ORDER BY period
	`, periodFormat, instructorId, "paid", startDate, endDate, req.CourseId, req.CourseId).
		Scan(&revenueByPeriod).Error; err != nil {
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

	// Top courses by revenue
	var topCourses []struct {
		CourseId    uint
		CourseTitle string
		Revenue     float64
		Orders      int
		Students    int
	}

	if err := ar.db.Raw(`
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
		WHERE courses.instructor_id = ?
			AND (? = 0 OR courses.id = ?)
		GROUP BY courses.id, courses.title
		ORDER BY revenue DESC
		LIMIT 5
	`, "paid", startDate, endDate, instructorId, req.CourseId, req.CourseId).
		Scan(&topCourses).Error; err != nil {
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

	// Calculate revenue growth (compare with previous period)
	periodDiff := endDate.Sub(startDate)
	previousStartDate := startDate.Add(-periodDiff)
	previousEndDate := startDate

	var previousRevenue struct {
		Total float64
	}
	if err := ar.db.Model(&models.Order{}).
		Select("COALESCE(SUM(final_price), 0) as total").
		Joins("JOIN courses ON courses.id = orders.course_id").
		Where("courses.instructor_id = ? AND orders.payment_status = ? AND orders.paid_at BETWEEN ? AND ?",
			instructorId, "paid", previousStartDate, previousEndDate).
		Scan(&previousRevenue).Error; err != nil {
		return nil, err
	}

	if previousRevenue.Total > 0 {
		response.RevenueGrowth = ((response.TotalRevenue - previousRevenue.Total) / previousRevenue.Total) * 100
	}

	return &response, nil
}

func (ar *DBAnalyticsRepository) GetStudentAnalytics(instructorId uint, req *dto.StudentAnalyticsRequest) (*dto.StudentAnalyticsResponse, error) {
	var response dto.StudentAnalyticsResponse

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
	baseQuery := ar.db.Model(&models.Enrollment{}).
		Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("courses.instructor_id = ?", instructorId)

	if req.CourseId != 0 {
		baseQuery = baseQuery.Where("enrollments.course_id = ?", req.CourseId)
	}

	// Total students
	var totalStudents int64
	if err := baseQuery.Count(&totalStudents).Error; err != nil {
		return nil, err
	}
	response.TotalStudents = int(totalStudents)

	// Active students
	var activeStudents int64
	if err := baseQuery.Where("enrollments.status = ?", "active").Count(&activeStudents).Error; err != nil {
		return nil, err
	}
	response.ActiveStudents = int(activeStudents)

	// Completed students
	var completedStudents int64
	if err := baseQuery.Where("enrollments.status = ?", "completed").Count(&completedStudents).Error; err != nil {
		return nil, err
	}
	response.CompletedStudents = int(completedStudents)

	// Dropped students
	var droppedStudents int64
	if err := baseQuery.Where("enrollments.status = ?", "dropped").Count(&droppedStudents).Error; err != nil {
		return nil, err
	}
	response.DroppedStudents = int(droppedStudents)

	// New students in period
	var newStudents int64
	if err := baseQuery.Where("enrollments.enrolled_at BETWEEN ? AND ?", startDate, endDate).
		Count(&newStudents).Error; err != nil {
		return nil, err
	}
	response.NewStudents = int(newStudents)

	// Students by period
	var studentsByPeriod []struct {
		Period         string
		NewStudents    int
		ActiveStudents int
		Completed      int
	}

	if err := ar.db.Raw(`
		SELECT 
			TO_CHAR(enrollments.enrolled_at, 'YYYY-MM-DD') as period,
			COUNT(*) as new_students,
			COUNT(CASE WHEN enrollments.status = 'active' THEN 1 END) as active_students,
			COUNT(CASE WHEN enrollments.status = 'completed' THEN 1 END) as completed
		FROM enrollments
		JOIN courses ON courses.id = enrollments.course_id
		WHERE courses.instructor_id = ?
			AND enrollments.enrolled_at BETWEEN ? AND ?
			AND (? = 0 OR enrollments.course_id = ?)
		GROUP BY period
		ORDER BY period
	`, instructorId, startDate, endDate, req.CourseId, req.CourseId).
		Scan(&studentsByPeriod).Error; err != nil {
		return nil, err
	}

	for _, item := range studentsByPeriod {
		response.StudentsByPeriod = append(response.StudentsByPeriod, dto.StudentPeriodItem{
			Period:         item.Period,
			NewStudents:    item.NewStudents,
			ActiveStudents: item.ActiveStudents,
			Completed:      item.Completed,
		})
	}

	// Students by course
	var studentsByCourse []struct {
		CourseId          uint
		CourseTitle       string
		TotalStudents     int
		ActiveStudents    int
		CompletedStudents int
		AverageProgress   float64
	}

	if err := ar.db.Raw(`
		SELECT 
			courses.id as course_id,
			courses.title as course_title,
			COUNT(enrollments.id) as total_students,
			COUNT(CASE WHEN enrollments.status = 'active' THEN 1 END) as active_students,
			COUNT(CASE WHEN enrollments.status = 'completed' THEN 1 END) as completed_students,
			COALESCE(AVG(enrollments.progress_percentage), 0) as average_progress
		FROM courses
		LEFT JOIN enrollments ON enrollments.course_id = courses.id
		WHERE courses.instructor_id = ?
			AND (? = 0 OR courses.id = ?)
		GROUP BY courses.id, courses.title
		ORDER BY total_students DESC
	`, instructorId, req.CourseId, req.CourseId).
		Scan(&studentsByCourse).Error; err != nil {
		return nil, err
	}

	for _, item := range studentsByCourse {
		completionRate := 0.0
		if item.TotalStudents > 0 {
			completionRate = float64(item.CompletedStudents) / float64(item.TotalStudents) * 100
		}

		response.StudentsByCourse = append(response.StudentsByCourse, dto.StudentCourseItem{
			CourseId:          item.CourseId,
			CourseTitle:       item.CourseTitle,
			TotalStudents:     item.TotalStudents,
			ActiveStudents:    item.ActiveStudents,
			CompletedStudents: item.CompletedStudents,
			AverageProgress:   item.AverageProgress,
			CompletionRate:    completionRate,
		})
	}

	// Average progress
	var avgProgress struct {
		Avg float64
	}
	if err := baseQuery.Select("COALESCE(AVG(enrollments.progress_percentage), 0) as avg").
		Scan(&avgProgress).Error; err != nil {
		return nil, err
	}
	response.AverageProgress = avgProgress.Avg

	// Completion rate
	if response.TotalStudents > 0 {
		response.CompletionRate = float64(response.CompletedStudents) / float64(response.TotalStudents) * 100
	}

	// Retention rate (active + completed / total)
	if response.TotalStudents > 0 {
		response.RetentionRate = float64(response.ActiveStudents+response.CompletedStudents) / float64(response.TotalStudents) * 100
	}

	return &response, nil
}
