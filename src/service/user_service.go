package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
	"mime/multipart"
	"strings"
	"time"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (us *userService) GetProfile(userId uint) (*dto.UserProfile, error) {
	user, err := us.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	return &dto.UserProfile{
		Id:            user.Id,
		Username:      user.Username,
		Email:         user.Email,
		FullName:      user.FullName,
		Phone:         user.Phone,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
	}, nil
}

func (us *userService) UpdateProfile(userId uint, req *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	// 1. Kiểm tra user có tồn tại không
	existingUser, err := us.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}
	if existingUser.Status != "active" {
		return nil, utils.NewError("User account is not active", utils.ErrCodeForbidden)
	}

	// 2. Chuẩn bị dữ liệu cập nhật - chỉ cập nhật các field có giá trị
	updates := make(map[string]interface{})

	if req.FullName != "" {
		updates["full_name"] = strings.TrimSpace(req.FullName)
	}

	if req.Phone != "" {
		updates["phone"] = strings.TrimSpace(req.Phone)
	}

	if req.Bio != "" {
		updates["bio"] = strings.TrimSpace(req.Bio)
	}

	if req.AvatarURL != "" {
		updates["avatar_url"] = strings.TrimSpace(req.AvatarURL)
	}

	// Luôn cập nhật updated_at
	updates["updated_at"] = time.Now()

	// 3. Cập nhật nếu có dữ liệu
	if len(updates) > 1 {
		if err := us.userRepo.UpdateProfile(userId, updates); err != nil {
			return nil, utils.WrapError(err, "Failed to update profile", utils.ErrCodeInternal)
		}
	}

	// 4. Lấy thông tin user đã cập nhật
	updatedUser, err := us.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get updated profile", utils.ErrCodeInternal)
	}

	// 5. Trả về response
	return &dto.UpdateProfileResponse{
		Id:            updatedUser.Id,
		Username:      updatedUser.Username,
		Email:         updatedUser.Email,
		FullName:      updatedUser.FullName,
		Phone:         updatedUser.Phone,
		Bio:           updatedUser.Bio,
		AvatarURL:     updatedUser.AvatarURL,
		Role:          updatedUser.Role,
		Status:        updatedUser.Status,
		EmailVerified: updatedUser.EmailVerified,
		CreatedAt:     updatedUser.CreatedAt,
		UpdatedAt:     updatedUser.UpdatedAt,
	}, nil
}

func (us *userService) ChangePassword(userId uint, req *dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error) {
	// 1. Validate confirm password
	if req.NewPassword != req.ConfirmPassword {
		return nil, utils.NewError("New password and confirm password do not match", utils.ErrCodeBadRequest)
	}

	// 2. Kiểm tra user có tồn tại không
	existingUser, err := us.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// 3. Kiểm tra trạng thái tài khoản
	if existingUser.Status != "active" {
		return nil, utils.NewError("User account is not active", utils.ErrCodeForbidden)
	}

	// 4. Verify current password
	if !utils.CheckPassword(existingUser.Password, req.CurrentPassword) {
		return nil, utils.NewError("Current password is incorrect", utils.ErrCodeUnauthorized)
	}

	// 5. Kiểm tra mật khẩu mới không được giống mật khẩu cũ
	if utils.CheckPassword(existingUser.Password, req.NewPassword) {
		return nil, utils.NewError("New password must be different from current password", utils.ErrCodeBadRequest)
	}

	// 6. Hash mật khẩu mới
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to hash new password", utils.ErrCodeInternal)
	}

	// 7. Cập nhật mật khẩu
	if err := us.userRepo.ChangePassword(userId, hashedPassword); err != nil {
		return nil, utils.WrapError(err, "Failed to change password", utils.ErrCodeInternal)
	}

	return &dto.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}

func (us *userService) UploadAvatar(userId uint, file *multipart.FileHeader) (*dto.UploadAvatarResponse, error) {
	// 1. Kiểm tra user có tồn tại không
	existingUser, err := us.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra trạng thái tài khoản
	if existingUser.Status != "active" {
		return nil, utils.NewError("User account is not active", utils.ErrCodeForbidden)
	}

	// 3. Validate file specifically for avatar (optional - sử dụng hàm riêng)
	if err := utils.ValidateAvatarFile(file); err != nil {
		return nil, utils.WrapError(err, "Invalid avatar file", utils.ErrCodeBadRequest)
	}

	// 4. Validate và lưu file
	uploadDir := "../../src/uploads/avatars"

	fileName, err := utils.ValidateAndSaveFile(file, uploadDir)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to upload avatar", utils.ErrCodeBadRequest)
	}

	// 5. Tạo URL cho avatar
	baseURL := utils.GetEnv("BASE_URL", "http://localhost:8080")
	avatarURL := fmt.Sprintf("%s/uploads/avatars/%s", baseURL, fileName)

	// 6. Cập nhật avatar URL trong database
	if err := us.userRepo.UpdateAvatar(userId, avatarURL); err != nil {
		return nil, utils.WrapError(err, "Failed to update avatar URL", utils.ErrCodeInternal)
	}

	return &dto.UploadAvatarResponse{
		Message:   "Avatar uploaded successfully",
		AvatarURL: avatarURL,
	}, nil
}
