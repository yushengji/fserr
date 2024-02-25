package fserr

import (
	"errors"
	"fmt"
	"strings"
)

// New 创建新的错误，支持格式化占位符
func New(format string, args ...any) error {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return &fundamental{
		msg:   msg,
		stack: callers(),
	}
}

// Wrap 包装已有错误，支持格式化占位符
func Wrap(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	if len(strings.TrimSpace(format)) == 0 {
		return err
	}
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return &withMessage{
		cause: wrapStack(err),
		msg:   msg,
	}
}

// WithCode 创建带有错误码的error，支持格式化占位符
// 使用option可以替换其中信息
func WithCode[T codeType](err error, businessCode T, options ...Option) error {
	code := getCode(serviceCode.Load() + int(businessCode))
	ret := &withCode{
		cause:        wrapStack(err),
		Msg:          code.Message,
		HttpCode:     code.HttpCode,
		BusinessCode: code.BusinessCode,
	}
	for _, option := range options {
		option(ret)
	}

	if err == nil {
		return &withStack{
			error: ret,
			stack: callers(),
		}
	}

	return ret
}

func WithStack(err error) error {
	return &withStack{
		error: err,
		stack: callers(),
	}
}

// UnWrap 获取包装过的error
// 若error1使用 Wrap 包装后产出错误error2
// 当使用 UmWrap 后，返回的是error1
func UnWrap(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Is 判断err是否是target类型
// 相比==判断错误，该方式会进行不断地类似于 UmWrap 的操作，
// 将 UmWrap 后的错误进行比较
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 匹配最外层的与target类型相同的error，将其赋值给target
func As(err error, target any) bool { return errors.As(err, target) }

// ParseCode 将错误解析为错误码错误
// 若err不是错误码错误，则包裹传递错误，其他信息为默认错误码信息
// 若err为错误码错误，将其转换，将最外层错误信息作为最终错误信息返回
// 若想得到原始的错误码错误，可以使用As方法
func ParseCode(err error) *withCode {
	var target *withCode
	if As(err, &target) {
		return &withCode{
			cause:        target.cause,
			Msg:          outerMsg(err),
			HttpCode:     target.HttpCode,
			BusinessCode: target.BusinessCode,
		}
	}

	return &withCode{
		cause:        nil,
		Msg:          err.Error(),
		HttpCode:     defaultErrCode.HttpCode,
		BusinessCode: serviceCode.Load(),
	}
}

// IsCode 判断某个错误是否为某个错误码
func IsCode[T codeType](err error, code T) bool {
	var target *withCode
	if !As(err, &target) {
		return false
	}
	return target.BusinessCode == serviceCode.Load()+int(code)
}

// SetAppCode 设置服务错误码
// 模块码、模块错误码一共四位，指定应用码将拼接在前
// 例如应用码位101，模块码为1，模块错误码为21，那么最终业务错误码为:1010121
func SetAppCode[T codeType](code T) {
	serviceCode.Store(int64(code) * 10000)
}

// outerMsg 获取最外层的错误信息
func outerMsg(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(*withCode); ok {
		return e.Msg
	}
	if e, ok := err.(*withMessage); ok {
		return e.msg
	}
	if e, ok := err.(*withStack); ok {
		return outerMsg(e.error)
	}
	return err.Error()
}

func wrapStack(err error) error {
	switch err.(type) {
	case *fundamental, *withCode, *withMessage, *withStack:
		return err
	default:
		return &withStack{
			error: err,
			stack: callers(),
		}
	}
}
