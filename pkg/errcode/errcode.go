package errcode

import "fmt"

// BusinessError 业务错误（支持 errors.As 提取）
type BusinessError struct {
	Code    int
	Message string
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Err 创建业务错误
func Err(code int) *BusinessError {
	return &BusinessError{Code: code, Message: Msg(code)}
}

// WithMsg 创建带自定义消息的业务错误
func WithMsg(code int, msg string) *BusinessError {
	return &BusinessError{Code: code, Message: msg}
}
