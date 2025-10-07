package service

import (
	"fmt"
	"lms/src/utils"
)

type emailService struct {
	// Có thể thêm SMTP config, template engine, etc.
}

func NewEmailService() EmailService {
	return &emailService{}
}

func (es *emailService) SendPasswordResetEmail(email, resetToken, resetCode string) error {
	// Tạo reset URL
	baseURL := utils.GetEnv("FRONTEND_URL", "http://localhost:3000")
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, resetToken)
	// Template email (trong production nên dùng HTML template)
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
	Dear User,

	You have requested to reset your password. Please use one of the following methods:

	Method 1: Click the link below
	%s

	Method 2: Use this code: %s

	This link and code will expire in 1 hour.

	If you did not request this, please ignore this email.

	Best regards,
	LMS Team
`, resetURL, resetCode)

	// Trong development, chỉ log ra console
	fmt.Printf("=== PASSWORD RESET EMAIL ===\n")
	fmt.Printf("To: %s\n", email)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Body:\n%s\n", body)
	fmt.Printf("===============================\n")

	// TODO: Implement thật sự với SMTP
	// return es.sendSMTPEmail(email, subject, body)

	return nil
}

func (es *emailService) SendWelcomeEmail(email, fullName string) error {
	fmt.Printf("=== WELCOME EMAIL ===\n")
	fmt.Printf("To: %s\n", email)
	fmt.Printf("Welcome %s!\n", fullName)
	fmt.Printf("====================\n")
	return nil
}
