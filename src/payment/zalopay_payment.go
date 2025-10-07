package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"lms/src/dto"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ZaloPayPayment struct {
	config ZaloPayConfig
}

func NewZaloPayPayment(config ZaloPayConfig) PaymentGateway {
	return &ZaloPayPayment{
		config: config,
	}
}

func (z *ZaloPayPayment) GetName() string {
	return "zalopay"
}

func (z *ZaloPayPayment) CreatePayment(orderCode string, amount float64, description string) (*dto.CreatePaymentResponse, error) {
	// Tạo app_trans_id theo format: yymmdd_xxxx
	transID := fmt.Sprintf("%s_%s", time.Now().Format("060102"), orderCode)
	appTime := time.Now().UnixMilli()
	amountInt := int64(amount)

	// Embed data (có thể chứa thông tin bổ sung)
	embedData := map[string]interface{}{
		"redirecturl": z.config.ReturnURL,
	}
	embedDataStr, _ := json.Marshal(embedData)

	// Item data
	item := []map[string]interface{}{
		{
			"itemid":       orderCode,
			"itemname":     description,
			"itemprice":    amountInt,
			"itemquantity": 1,
		},
	}
	itemStr, _ := json.Marshal(item)

	// Tạo data string để tạo MAC
	dataStr := fmt.Sprintf("%d|%s|%s|%d|%s|%s|%s",
		z.config.AppId,
		transID,
		"LMS_USER",
		amountInt,
		appTime,
		embedDataStr,
		itemStr,
	)

	// Tạo MAC
	mac := z.generateMAC(dataStr, z.config.Key1)

	// Tạo request body
	requestBody := url.Values{}
	requestBody.Set("app_id", strconv.Itoa(z.config.AppId))
	requestBody.Set("app_user", "LMS_USER")
	requestBody.Set("app_time", strconv.FormatInt(appTime, 10))
	requestBody.Set("amount", strconv.FormatInt(amountInt, 10))
	requestBody.Set("app_trans_id", transID)
	requestBody.Set("embed_data", string(embedDataStr))
	requestBody.Set("item", string(itemStr))
	requestBody.Set("description", description)
	requestBody.Set("bank_code", "")
	requestBody.Set("callback_url", z.config.CallbackURL)
	requestBody.Set("mac", mac)

	// Gửi request
	resp, err := http.PostForm(z.config.Endpoint+"/create", requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to ZaloPay: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var zaloResp dto.ZaloPayPaymentResponse
	if err := json.Unmarshal(body, &zaloResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Kiểm tra kết quả
	if zaloResp.ReturnCode != 1 {
		return nil, fmt.Errorf("ZaloPay error: %s (code: %d)", zaloResp.ReturnMessage, zaloResp.ReturnCode)
	}

	return &dto.CreatePaymentResponse{
		PaymentUrl:    zaloResp.OrderUrl,
		PaymentMethod: "zalopay",
		OrderCode:     orderCode,
		Amount:        amount,
		Message:       "Payment URL created successfully",
	}, nil
}

func (z *ZaloPayPayment) VerifyCallback(data map[string]interface{}) (bool, error) {
	dataStr, ok := data["data"].(string)
	if !ok {
		return false, fmt.Errorf("data field not found")
	}

	receivedMac, ok := data["mac"].(string)
	if !ok {
		return false, fmt.Errorf("mac field not found")
	}

	// Tạo MAC để verify
	expectedMac := z.generateMAC(dataStr, z.config.Key2)

	return expectedMac == receivedMac, nil
}

func (z *ZaloPayPayment) GetPaymentStatus(orderCode string) (string, error) {
	// Tạo app_trans_id
	transID := fmt.Sprintf("%s_%s", time.Now().Format("060102"), orderCode)
	appTime := time.Now().UnixMilli()

	// Tạo data string
	dataStr := fmt.Sprintf("%d|%s|%d", z.config.AppId, transID, appTime)

	// Tạo MAC
	mac := z.generateMAC(dataStr, z.config.Key1)

	// Tạo request body
	requestBody := url.Values{}
	requestBody.Set("app_id", strconv.Itoa(z.config.AppId))
	requestBody.Set("app_trans_id", transID)
	requestBody.Set("mac", mac)

	// Gửi request
	queryEndpoint := z.config.Endpoint + "/query"
	resp, err := http.PostForm(queryEndpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to query ZaloPay: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var queryResp map[string]interface{}
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	returnCode, ok := queryResp["return_code"].(float64)
	if !ok {
		return dto.PaymentStatusPending, nil
	}

	switch int(returnCode) {
	case 1:
		return dto.PaymentStatusPaid, nil
	case 2:
		return dto.PaymentStatusFailed, nil
	case 3:
		return dto.PaymentStatusPending, nil
	default:
		return dto.PaymentStatusFailed, nil
	}
}

// generateMAC tạo HMAC SHA256 MAC
func (z *ZaloPayPayment) generateMAC(data string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// ParseZaloPayCallback parse callback data từ ZaloPay
func ParseZaloPayCallback(callbackData string) (*dto.ZaloPayCallbackData, error) {
	var data dto.ZaloPayCallbackData
	if err := json.Unmarshal([]byte(callbackData), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// GetOrderIdFromCallback lấy order ID từ ZaloPay callback
func (z *ZaloPayPayment) GetOrderIdFromCallback(data map[string]interface{}) (string, error) {
	dataStr, ok := data["data"].(string)
	if !ok {
		return "", fmt.Errorf("data field not found")
	}

	var callbackData dto.ZaloPayCallbackData
	if err := json.Unmarshal([]byte(dataStr), &callbackData); err != nil {
		return "", err
	}

	// app_trans_id format: yymmdd_orderCode
	// Extract orderCode từ app_trans_id
	appTransId := callbackData.AppTransId
	if len(appTransId) > 7 {
		return appTransId[7:], nil // Bỏ qua phần "yymmdd_"
	}

	return "", fmt.Errorf("invalid app_trans_id format")
}

// IsPaymentSuccess kiểm tra thanh toán thành công từ callback
func (z *ZaloPayPayment) IsPaymentSuccess(data map[string]interface{}) bool {
	dataStr, ok := data["data"].(string)
	if !ok {
		return false
	}

	var callbackData dto.ZaloPayCallbackData
	if err := json.Unmarshal([]byte(dataStr), &callbackData); err != nil {
		return false
	}

	return callbackData.ZpTransId > 0
}
