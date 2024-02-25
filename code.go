package fserr

import (
	"net/http"
	"sync/atomic"
)

type codeType interface {
	int8 | int16 | int32 | int64 | int |
		uint8 | uint16 | uint32 | uint64 | uint
}

type atomicServiceCode struct {
	atomic.Int64
}

func (a *atomicServiceCode) Load() int {
	return int(a.Int64.Load())
}

func (a *atomicServiceCode) Store(v int64) {
	a.Int64.Store(v)
}

var serviceCode atomicServiceCode

const (
	ErrBasic = iota + 1
	ErrDb
	ErrParam
	ErrRetry
	ErrServiceInvoke
)

// ErrCode 产生error所包含的错误码信息
// 支持http错误码、业务错误码
// 可通过下方快捷函数快速创建指定http码的错误码
type ErrCode struct {
	// HttpCode 该错误码建议的HTTP响应码
	HttpCode int `json:"httpCode,omitempty"`
	// BusinessCode 该错误码对应的业务码
	// +optional
	BusinessCode int `json:"businessCode,omitempty"`
	// Message 该错误码对应的错误信息
	// +optional
	Message string `json:"message,omitempty"`
}

// SetDefault 设置默认错误码
// 当错误码匹配失败时，提供的备选方案，已内置默认错误码，
// 它的HTTP码为200，业务码和信息均为零值
func SetDefault(httpCode, businessCode int, message string) {
	defaultErrCode = ErrCode{httpCode, serviceCode.Load() + businessCode, message}
}

// NewCode 创建指定信息的错误码
func NewCode(httpCode, businessCode int, message string) ErrCode {
	code := ErrCode{httpCode, serviceCode.Load() + businessCode, message}
	register(code)
	return code
}

// === 以下均为见名知意的业务码构建方式 ===

func NewOK(businessCode int, message string) ErrCode {
	return NewCode(http.StatusOK, businessCode, message)
}

func NewNotFound(businessCode int, message string) ErrCode {
	return NewCode(http.StatusNotFound, businessCode, message)
}

func NewAlreadyExists(businessCode int, message string) ErrCode {
	return NewCode(http.StatusConflict, businessCode, message)
}

func NewGenerateNameConflict(businessCode int, message string) ErrCode {
	return NewCode(http.StatusConflict, businessCode, message)
}

func NewUnauthorized(businessCode int, message string) ErrCode {
	return NewCode(http.StatusUnauthorized, businessCode, message)
}

func NewForbidden(businessCode int, message string) ErrCode {
	return NewCode(http.StatusForbidden, businessCode, message)
}

func NewConflict(businessCode int, message string) ErrCode {
	return NewCode(http.StatusConflict, businessCode, message)
}

func NewGone(businessCode int, message string) ErrCode {
	return NewCode(http.StatusGone, businessCode, message)
}

func NewBadRequest(businessCode int, message string) ErrCode {
	return NewCode(http.StatusBadRequest, businessCode, message)
}

func NewTooManyRequests(businessCode int, message string) ErrCode {
	return NewCode(http.StatusTooManyRequests, businessCode, message)
}

func NewServiceUnavailable(businessCode int, message string) ErrCode {
	return NewCode(http.StatusServiceUnavailable, businessCode, message)
}

func NewMethodNotSupported(businessCode int, message string) ErrCode {
	return NewCode(http.StatusMethodNotAllowed, businessCode, message)
}

func NewInternalError(businessCode int, message string) ErrCode {
	return NewCode(http.StatusInternalServerError, businessCode, message)
}

func NewTimeoutError(businessCode int, message string) ErrCode {
	return NewCode(http.StatusGatewayTimeout, businessCode, message)
}

func NewTooManyRequestsError(businessCode int, message string) ErrCode {
	return NewCode(http.StatusTooManyRequests, businessCode, message)
}

func NewRequestEntityTooLargeError(businessCode int, message string) ErrCode {
	return NewCode(http.StatusRequestEntityTooLarge, businessCode, message)
}
