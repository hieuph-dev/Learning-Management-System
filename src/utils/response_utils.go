package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST" // 400
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"             // 404
	ErrCodeConflict     ErrorCode = "CONFLICT"              // 409
	ErrCodeInternal     ErrorCode = "INTERNAL_SERVER_ERROR" // 500
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
)

// Lỗi của bạn muốn tạo
type AppError struct {
	Message string    // Thông báo lỗi cho user
	Code    ErrorCode // Mã lỗi chuẩn
	Err     error     // lỗi gốc
}

func (ae *AppError) Error() string {
	if ae.Message != "" {
		return ae.Message
	}
	if ae.Err != nil {
		return ae.Err.Error()
	}
	return "unknown error"
}

// Tạo mới lỗi của bạn
func NewError(message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

// Thêm lỗi gốc vào
func WrapError(err error, message string, code ErrorCode) error {
	return &AppError{
		Err:     err,     // Lỗi gốc
		Message: message, // Lỗi của mình
		Code:    code,
	}
}

// Trả về lỗi JSON
func ResponseError(ctx *gin.Context, err error) {
	// Kiểm tra nếu là AppError
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)

		response := gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		}

		// Thêm chi tiết lỗi nếu có
		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}
	// Lỗi không xác định
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}

// Trả về JSON thành công
func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
}

// Chỉ trả status code.
func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}

// Trả lỗi validate.
func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
}

// Chuyển mã lỗi sang HTTP status.
func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest, ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
