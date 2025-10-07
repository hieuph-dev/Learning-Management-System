package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
)

type authService struct {
	userRepo          repository.UserRepository
	passwordResetRepo repository.PasswordResetRepository
	emailService      EmailService
}

func NewAuthService(userRepo repository.UserRepository, passwordResetRepo repository.PasswordResetRepository, emailService EmailService) AuthService {
	return &authService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		emailService:      emailService,
	}
}

func (as *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {

	// 1. Check email & username co ton tai chua
	req.Email = utils.NormalizeString(req.Email)

	if _, exist := as.userRepo.FindByEmail(req.Email); exist {
		return nil, utils.NewError("email already exists", utils.ErrCodeConflict)
	}

	req.Username = utils.NormalizeString(req.Username)

	if _, exist := as.userRepo.FindByUsername(req.Username); exist {
		return nil, utils.NewError("username already exists", utils.ErrCodeConflict)
	}

	// 2. Hash passowrd
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
	}

	// 3. Tao user moi
	user := models.User{
		Username:      req.Username,
		Email:         req.Email,
		Password:      hashedPassword,
		FullName:      req.FullName,
		Phone:         req.Phone,
		Role:          "student", // Role mac dinh la student
		Status:        "active",
		EmailVerified: false,
	}

	// 4. Luu vao db
	if err := as.userRepo.Create(&user); err != nil {
		return nil, utils.WrapError(err, "failed to create user", utils.ErrCodeInternal)
	}

	// 5. Tao jwt tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.Id, user.Username, user.Role)
	if err != nil {
		return nil, utils.NewError("failed to create tokens", utils.ErrCodeInternal)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserProfile{
			Id:            user.Id,
			Username:      user.Username,
			Email:         user.Email,
			FullName:      user.FullName,
			Phone:         user.Phone,
			Role:          user.Role,
			Status:        user.Status,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
		},
	}, nil
}

func (as *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	req.Email = utils.NormalizeString(req.Email)

	user, exist := as.userRepo.FindByEmail(req.Email)

	if !exist {
		return nil, utils.NewError("invalid credentials mail", utils.ErrCodeUnauthorized)
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, utils.NewError("account is inactive", utils.ErrCodeForbidden)
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, utils.NewError("invalid credentials", utils.ErrCodeUnauthorized)
	}

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.Id, user.Email, user.Role)
	if err != nil {
		return nil, utils.NewError("failed to create tokens", utils.ErrCodeInternal)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserProfile{
			Id:            user.Id,
			Username:      user.Username,
			Email:         user.Email,
			FullName:      user.FullName,
			Phone:         user.Phone,
			Role:          user.Role,
			Status:        user.Status,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
		},
	}, nil
}

func (as *authService) GetProfile(userId uint) (*dto.UserProfile, error) {
	user, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("user not found", utils.ErrCodeNotFound)
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

func (as *authService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, utils.NewError("invalid refresh token", utils.ErrCodeUnauthorized)
	}

	// Kiểm tra xem token có phải là refresh token không
	if claims.Subject != "refresh" {
		return nil, utils.NewError("invalid token type", utils.ErrCodeUnauthorized)
	}

	// Kiểm tra user có tồn tại và active không
	user, err := as.userRepo.FindById(claims.UserId)
	if err != nil {
		return nil, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	if user.Status != "active" {
		return nil, utils.NewError("account is inactive", utils.ErrCodeForbidden)
	}

	// Tạo tokens mới
	newAccessToken, newRefreshToken, err := utils.GenerateTokens(user.Id, user.Username, user.Role)
	if err != nil {
		return nil, utils.NewError("failed to create tokens", utils.ErrCodeInternal)
	}

	return &dto.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (as *authService) ForgotPassword(req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	// 1. Normalize email
	req.Email = utils.NormalizeString(req.Email)

	// 2. Kiểm tra email có tồn tại không
	user, exist := as.userRepo.FindByEmail(req.Email)
	if !exist {
		// Không thông báo email không tồn tại để tránh email enumeration attack
		return &dto.ForgotPasswordResponse{
			Message: "If this email exists, you will receive a password reset link shortly",
			Email:   req.Email,
		}, nil
	}

	// 3. Kiểm tra trạng thái tài khoản
	if user.Status != "active" {
		return nil, utils.NewError("account is inactive", utils.ErrCodeForbidden)
	}

	// 4. Xóa các token reset cũ của email này
	if err := as.passwordResetRepo.DeleteByEmail(req.Email); err != nil {
		fmt.Printf("Error deleting old reset tokens: %v\n", err)
	}

	// 5. Tạo token mới
	secureToken, readableCode, hashToken, err := utils.GeneratePasswordResetToken()
	if err != nil {
		return nil, utils.WrapError(err, "failed to generate reset token", utils.ErrCodeInternal)
	}

	// 6. Tạo record PasswordReset
	resetRecord := &models.PasswordReset{
		Email:     req.Email,
		Token:     hashToken, // Luu hash token, khong luu raw token
		ExpiresAt: utils.GetResetTokenExpiry(),
		Used:      false,
	}

	// 7. Lưu vào database
	if err := as.passwordResetRepo.Create(resetRecord); err != nil {
		return nil, utils.WrapError(err, "failed to create reset record", utils.ErrCodeInternal)
	}

	// 8. Gửi email (sử dụng secureToken raw, không phải hash)
	if err := as.emailService.SendPasswordResetEmail(req.Email, secureToken, readableCode); err != nil {
		fmt.Printf("Failed to send reset email to %s: %v\n", req.Email, err)
		// Không trả lỗi cho user để tránh leak thông tin
	}

	return &dto.ForgotPasswordResponse{
		Message: "If this email exists, you will receive a password reset link shortly",
		Email:   req.Email,
	}, nil
}

func (as *authService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// 1. Hash token để so sánh với DB
	hashedToken := utils.HashToken(req.Token)

	// 2. Tìm token trong DB
	resetRecord, err := as.passwordResetRepo.FindByToken(hashedToken)
	if err != nil {
		return utils.NewError("invalid or expired reset token", utils.ErrCodeUnauthorized)
	}

	// 3. Kiểm tra token đã được sử dụng chưa
	if resetRecord.Used {
		return utils.NewError("reset token has already been used", utils.ErrCodeUnauthorized)
	}

	// 4. Kiểm tra token có hết hạn chưa
	if utils.IsTokenExpired(resetRecord.ExpiresAt) {
		return utils.NewError("reset token has expired", utils.ErrCodeUnauthorized)
	}

	// 5. Tìm user theo email
	user, exist := as.userRepo.FindByEmail(resetRecord.Email)
	if !exist {
		return utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	// 6. Kiểm tra trạng thái tài khoản
	if user.Status != "active" {
		return utils.NewError("account is inactive", utils.ErrCodeForbidden)
	}

	// 7. Hash mật khẩu mới
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
	}

	// 8. Cập nhật mật khẩu user
	if err := as.userRepo.UpdatePassword(user.Id, hashedPassword); err != nil {
		return utils.WrapError(err, "failed  to update password", utils.ErrCodeInternal)
	}

	// 9. Đánh dấu token đã được sử dụng
	if err := as.passwordResetRepo.MarkAsUsed(resetRecord.Id); err != nil {
		fmt.Printf("Failed to mark token as userd: %v\n", err)
	}

	return nil
}
