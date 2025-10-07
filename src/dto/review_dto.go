package dto

type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required,min=10,max=1000"`
}

type CreateReviewResponse struct {
	Id       uint   `json:"id"`
	CourseId uint   `json:"course_id"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
	Message  string `json:"message"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment *string `json:"comment" binding:"omitempty,min=10,max=1000"`
}

type UpdateReviewResponse struct {
	Id      uint   `json:"id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Message string `json:"message"`
}

type DeleteReviewResponse struct {
	Message string `json:"message"`
}
