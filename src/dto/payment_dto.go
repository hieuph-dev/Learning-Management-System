package dto

import "time"

// Payment Method Constants
const (
	PaymentMethodMomo    = "momo"
	PaymentMethodZaloPay = "zalopay"
	PaymentMethodFree    = "free"
)

// Payment Status Constants
const (
	PaymentStatusPending   = "pending"
	PaymentStatusPaid      = "paid"
	PaymentStatusFailed    = "failed"
	PaymentStatusCancelled = "cancelled"
)

// ============ MoMo Payment DTOs ============

type MomoPaymentRequest struct {
	OrderId     string  `json:"orderId" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	OrderInfo   string  `json:"orderInfo" binding:"required"`
	RedirectUrl string  `json:"redirectUrl"`
	IpnUrl      string  `json:"ipnUrl"`
}

type MomoPaymentResponse struct {
	PartnerCode  string `json:"partnerCode"`
	RequestId    string `json:"requestId"`
	OrderId      string `json:"orderId"`
	Amount       int64  `json:"amount"`
	ResponseTime int64  `json:"responseTime"`
	Message      string `json:"message"`
	ResultCode   int    `json:"resultCode"`
	PayUrl       string `json:"payUrl"`
	DeepLink     string `json:"deeplink"`
	QrCodeUrl    string `json:"qrCodeUrl"`
}

type MomoIPNRequest struct {
	PartnerCode  string `json:"partnerCode"`
	OrderId      string `json:"orderId"`
	RequestId    string `json:"requestId"`
	Amount       int64  `json:"amount"`
	OrderInfo    string `json:"orderInfo"`
	OrderType    string `json:"orderType"`
	TransId      int64  `json:"transId"`
	ResultCode   int    `json:"resultCode"`
	Message      string `json:"message"`
	PayType      string `json:"payType"`
	ResponseTime int64  `json:"responseTime"`
	ExtraData    string `json:"extraData"`
	Signature    string `json:"signature"`
}

// ============ ZaloPay Payment DTOs ============

type ZaloPayPaymentRequest struct {
	OrderId     string  `json:"orderId" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
	RedirectUrl string  `json:"redirectUrl"`
	CallbackUrl string  `json:"callbackUrl"`
}

type ZaloPayPaymentResponse struct {
	ReturnCode       int    `json:"return_code"`
	ReturnMessage    string `json:"return_message"`
	SubReturnCode    int    `json:"sub_return_code"`
	SubReturnMessage string `json:"sub_return_message"`
	OrderUrl         string `json:"order_url"`
	ZpTransToken     string `json:"zp_trans_token"`
	OrderToken       string `json:"order_token"`
}

type ZaloPayCallbackRequest struct {
	Data string `json:"data" binding:"required"`
	Mac  string `json:"mac" binding:"required"`
	Type int    `json:"type"`
}

type ZaloPayCallbackData struct {
	AppId          int    `json:"app_id"`
	AppTransId     string `json:"app_trans_id"`
	AppTime        int64  `json:"app_time"`
	AppUser        string `json:"app_user"`
	Amount         int64  `json:"amount"`
	EmbedData      string `json:"embed_data"`
	Item           string `json:"item"`
	ZpTransId      int64  `json:"zp_trans_id"`
	ServerTime     int64  `json:"server_time"`
	Channel        int    `json:"channel"`
	MerchantUserId string `json:"merchant_user_id"`
}

// ============ Unified Payment DTOs ============

type CreatePaymentRequest struct {
	OrderId       uint   `json:"order_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=momo zalopay"`
	ReturnUrl     string `json:"return_url"`
}

type CreatePaymentResponse struct {
	PaymentUrl    string  `json:"payment_url"`
	PaymentMethod string  `json:"payment_method"`
	OrderId       uint    `json:"order_id"`
	OrderCode     string  `json:"order_code"`
	Amount        float64 `json:"amount"`
	QrCodeUrl     string  `json:"qr_code_url,omitempty"`
	DeepLink      string  `json:"deep_link,omitempty"`
	Message       string  `json:"message"`
}

type PaymentCallbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	OrderId uint   `json:"order_id,omitempty"`
}

type CheckPaymentStatusRequest struct {
	OrderId uint `form:"order_id" binding:"required"`
}

type CheckPaymentStatusResponse struct {
	OrderId       uint       `json:"order_id"`
	OrderCode     string     `json:"order_code"`
	PaymentStatus string     `json:"payment_status"`
	PaymentMethod string     `json:"payment_method"`
	Amount        float64    `json:"amount"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	Message       string     `json:"message"`
}
