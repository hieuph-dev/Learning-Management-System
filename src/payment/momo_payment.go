package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"lms/src/dto"
	"net/http"

	"github.com/google/uuid"
)

type MomoPayment struct {
	config MomoConfig
}

func NewMomoPayment(config MomoConfig) PaymentGateway {
	return &MomoPayment{
		config: config,
	}
}

func (m *MomoPayment) GetName() string {
	return "momo"
}

func (m *MomoPayment) CreatePayment(orderCode string, amount float64, description string) (*dto.CreatePaymentResponse, error) {
	requestId := uuid.New().String()
	orderId := orderCode
	amountInt := int64(amount)
	orderInfo := description
	requestType := "captureWallet"
	extraData := ""

	// Tạo raw signature
	rawSignature := fmt.Sprintf(
		"accessKey=%s&amount=%d&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=%s",
		m.config.AccessKey,
		amountInt,
		extraData,
		m.config.IPNUrl,
		orderId,
		orderInfo,
		m.config.PartnerCode,
		m.config.ReturnURL,
		requestId,
		requestType,
	)

	// Tạo signature
	signature := m.generateSignature(rawSignature, m.config.SecretKey)

	// Tạo request body
	requestBody := map[string]interface{}{
		"partnerCode": m.config.PartnerCode,
		"partnerName": "LMS Platform",
		"storeId":     "LMSStore",
		"requestId":   requestId,
		"amount":      amountInt,
		"orderId":     orderId,
		"orderInfo":   orderInfo,
		"redirectUrl": m.config.ReturnURL,
		"ipnUrl":      m.config.IPNUrl,
		"lang":        "vi",
		"extraData":   extraData,
		"requestType": requestType,
		"signature":   signature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Gửi request tới MoMo
	resp, err := http.Post(m.config.Endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to MoMo: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var momoResp dto.MomoPaymentResponse
	if err := json.Unmarshal(body, &momoResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Kiểm tra kết quả
	if momoResp.ResultCode != 0 {
		return nil, fmt.Errorf("MoMo error: %s (code: %d)", momoResp.Message, momoResp.ResultCode)
	}

	return &dto.CreatePaymentResponse{
		PaymentUrl:    momoResp.PayUrl,
		PaymentMethod: "momo",
		OrderCode:     orderCode,
		Amount:        amount,
		QrCodeUrl:     momoResp.QrCodeUrl,
		DeepLink:      momoResp.DeepLink,
		Message:       "Payment URL created successfully",
	}, nil
}

func (m *MomoPayment) VerifyCallback(data map[string]interface{}) (bool, error) {
	// Lấy các thông tin từ callback
	partnerCode, _ := data["partnerCode"].(string)
	orderId, _ := data["orderId"].(string)
	requestId, _ := data["requestId"].(string)
	amount, _ := data["amount"].(float64)
	orderInfo, _ := data["orderInfo"].(string)
	orderType, _ := data["orderType"].(string)
	transId, _ := data["transId"].(float64)
	resultCode, _ := data["resultCode"].(float64)
	message, _ := data["message"].(string)
	payType, _ := data["payType"].(string)
	responseTime, _ := data["responseTime"].(float64)
	extraData, _ := data["extraData"].(string)
	receivedSignature, _ := data["signature"].(string)

	// Tạo raw signature để verify
	rawSignature := fmt.Sprintf(
		"accessKey=%s&amount=%.0f&extraData=%s&message=%s&orderId=%s&orderInfo=%s&orderType=%s&partnerCode=%s&payType=%s&requestId=%s&responseTime=%.0f&resultCode=%.0f&transId=%.0f",
		m.config.AccessKey,
		amount,
		extraData,
		message,
		orderId,
		orderInfo,
		orderType,
		partnerCode,
		payType,
		requestId,
		responseTime,
		resultCode,
		transId,
	)

	expectedSignature := m.generateSignature(rawSignature, m.config.SecretKey)

	return expectedSignature == receivedSignature, nil
}

func (m *MomoPayment) GetPaymentStatus(orderCode string) (string, error) {
	requestId := uuid.New().String()

	// Tạo raw signature cho query
	rawSignature := fmt.Sprintf(
		"accessKey=%s&orderId=%s&partnerCode=%s&requestId=%s",
		m.config.AccessKey,
		orderCode,
		m.config.PartnerCode,
		requestId,
	)

	signature := m.generateSignature(rawSignature, m.config.SecretKey)

	requestBody := map[string]interface{}{
		"partnerCode": m.config.PartnerCode,
		"requestId":   requestId,
		"orderId":     orderCode,
		"lang":        "vi",
		"signature":   signature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Endpoint để query status
	queryEndpoint := m.config.Endpoint + "/query"
	resp, err := http.Post(queryEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to query MoMo: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var queryResp map[string]interface{}
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	resultCode, ok := queryResp["resultCode"].(float64)
	if !ok {
		return dto.PaymentStatusPending, nil
	}

	if resultCode == 0 {
		return dto.PaymentStatusPaid, nil
	} else if resultCode == 1006 {
		return dto.PaymentStatusPending, nil
	}

	return dto.PaymentStatusFailed, nil
}

// generateSignature tạo HMAC SHA256 signature
func (m *MomoPayment) generateSignature(data string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Helper function để parse MoMo callback request
func ParseMomoCallback(data map[string]interface{}) (*dto.MomoIPNRequest, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var ipnReq dto.MomoIPNRequest
	if err := json.Unmarshal(jsonData, &ipnReq); err != nil {
		return nil, err
	}

	return &ipnReq, nil
}

// GetOrderIdFromCallback lấy order ID từ MoMo callback
func (m *MomoPayment) GetOrderIdFromCallback(data map[string]interface{}) (string, error) {
	orderId, ok := data["orderId"].(string)
	if !ok {
		return "", fmt.Errorf("orderId not found in callback")
	}
	return orderId, nil
}

// IsPaymentSuccess kiểm tra thanh toán thành công từ callback
func (m *MomoPayment) IsPaymentSuccess(data map[string]interface{}) bool {
	resultCode, ok := data["resultCode"].(float64)
	if !ok {
		return false
	}
	return resultCode == 0
}
