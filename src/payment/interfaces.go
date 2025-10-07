package payment

import "lms/src/dto"

// PaymentGateway định nghĩa interface chung cho các cổng thanh toán
type PaymentGateway interface {
	// CreatePayment tạo payment và trả về URL thanh toán
	CreatePayment(orderCode string, amount float64, description string) (*dto.CreatePaymentResponse, error)

	// VerifyCallback xác thực callback từ payment gateway
	VerifyCallback(data map[string]interface{}) (bool, error)

	// GetPaymentStatus kiểm tra trạng thái thanh toán
	GetPaymentStatus(orderCode string) (string, error)

	// GetName trả về tên của payment gateway
	GetName() string
}

// PaymentConfig chứa cấu hình chung cho payment
type PaymentConfig struct {
	BaseURL     string
	ReturnURL   string
	CallbackURL string
}

// MomoConfig cấu hình cho MoMo
type MomoConfig struct {
	PartnerCode string
	AccessKey   string
	SecretKey   string
	Endpoint    string
	ReturnURL   string
	IPNUrl      string
}

// ZaloPayConfig cấu hình cho ZaloPay
type ZaloPayConfig struct {
	AppId       int
	Key1        string
	Key2        string
	Endpoint    string
	CallbackURL string
	ReturnURL   string
}
