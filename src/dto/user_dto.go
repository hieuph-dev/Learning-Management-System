package dto

import "time"

type UpdateProfileRequest struct {
	FullName  string `json:"full_name" binding:"omitempty,min=2,max=100"`
	Phone     string `json:"phone" binding:"omitempty,max=20"`
	Bio       string `json:"bio" binding:"omitempty,max=500"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

type UpdateProfileResponse struct {
	Id            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Bio           string    `json:"bio"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,password_strong,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

type UploadAvatarResponse struct {
	Message   string `json:"message"`
	AvatarURL string `json:"avatar_url"`
}
