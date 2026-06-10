package errors

import (
	stderrors "errors"
	"net/http"
)

// 用户端 HTTP 错误响应契约。
type ResultCode struct {
	// 全局唯一业务 ID。
	Code int
	// 默认英文文案。
	Message string
	// HTTP 状态。
	HTTPStatus int
}

var (
	// 请求参数校验失败。
	ParameterValidationError = ResultCode{Code: 10005, Message: "Parameter validation failed", HTTPStatus: http.StatusBadRequest}
	// 认证失败或登录凭证无效。
	InvalidCredentials = ResultCode{Code: 20002, Message: "Authentication failed", HTTPStatus: http.StatusUnauthorized}
	// 当前用户缺少访问权限。
	AccessDenied = ResultCode{Code: 403, Message: "Access denied", HTTPStatus: http.StatusForbidden}
	// 请求的数据不存在。
	DataNotFound = ResultCode{Code: 40004, Message: "Data not found", HTTPStatus: http.StatusNotFound}
	// 服务端内部错误。
	InternalError = ResultCode{Code: 50000, Message: "Service error", HTTPStatus: http.StatusInternalServerError}
	// 接口尚未实现。
	NotImplementedError = ResultCode{Code: 50001, Message: "API is not implemented yet", HTTPStatus: http.StatusNotImplemented}
)

// 用户端可直接写入响应的业务错误。
type Error struct {
	// 全局唯一业务 ID。
	Code int
	// 对外稳定文案。
	Message string
	// HTTP 状态。
	HTTPStatus int
}

// 实现 error 接口，返回对外稳定文案。
func (e *Error) Error() string {
	return e.Message
}

// 按用户端错误码注册项创建业务错误。
func New(resultCode ResultCode) *Error {
	return &Error{Code: resultCode.Code, Message: resultCode.Message, HTTPStatus: resultCode.HTTPStatus}
}

// 保留业务码和 HTTP 状态，只替换对外文案。
func WithMessage(resultCode ResultCode, message string) *Error {
	return &Error{Code: resultCode.Code, Message: message, HTTPStatus: resultCode.HTTPStatus}
}

// 请求参数错误。
func BadRequest(message string) *Error {
	if message == "" {
		return New(ParameterValidationError)
	}
	return WithMessage(ParameterValidationError, message)
}

// 认证失败。
func Unauthorized() *Error {
	return New(InvalidCredentials)
}

// 权限不足。
func Forbidden() *Error {
	return New(AccessDenied)
}

// 数据不存在。
func NotFound() *Error {
	return New(DataNotFound)
}

// 服务端内部错误。
func Internal() *Error {
	return New(InternalError)
}

// 接口尚未实现。
func NotImplemented() *Error {
	return New(NotImplementedError)
}

// 将普通错误转换为用户端响应错误。
func AsAppError(err error) *Error {
	if err == nil {
		return nil
	}
	var appErr *Error
	if stderrors.As(err, &appErr) {
		return appErr
	}
	return New(InternalError)
}

// 将业务错误转换为指定语言的用户端文案。
func Localize(err error, language string) *Error {
	appErr := AsAppError(err)
	if appErr == nil {
		return nil
	}
	if appErr.Message != "" && appErr.Message != MessageForCode(appErr.Code, "en") {
		return appErr
	}
	message := MessageForCode(appErr.Code, language)
	return &Error{Code: appErr.Code, Message: message, HTTPStatus: appErr.HTTPStatus}
}
