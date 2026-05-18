package service

import (
	"reflect"
	"strings"
	"unicode/utf8"
)

// Masking 支持的数据脱敏类型
// 使用 struct tag `masking:"phone"` 等标记字段
var maskingFuncs = map[string]func(string) string{
	"phone":    maskPhone,
	"email":    maskEmail,
	"idCard":   maskIDCard,
	"bankCard": maskBankCard,
	"name":     maskName,
	"address":  maskAddress,
}

// MaskStruct 对 struct 进行脱敏处理
func MaskStruct(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	maskStructFields(val)
}

func maskStructFields(val reflect.Value) {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			maskStructFields(val.Field(i))
			continue
		}
		maskTag := field.Tag.Get("masking")
		if maskTag == "" {
			continue
		}
		fn, ok := maskingFuncs[maskTag]
		if !ok {
			continue
		}
		f := val.Field(i)
		if f.Kind() == reflect.String && f.String() != "" {
			f.SetString(fn(f.String()))
		}
	}
}

// 手机号脱敏: 13812345678 → 138****5678
func maskPhone(s string) string {
	if utf8.RuneCountInString(s) < 7 {
		return s
	}
	runes := []rune(s)
	return string(runes[:3]) + "****" + string(runes[7:])
}

// 邮箱脱敏: test@example.com → t***@example.com
func maskEmail(s string) string {
	at := strings.Index(s, "@")
	if at <= 0 {
		return s
	}
	return string(s[0]) + "***" + s[at:]
}

// 身份证脱敏: 320123199001011234 → 3201**********1234
func maskIDCard(s string) string {
	if utf8.RuneCountInString(s) < 8 {
		return s
	}
	runes := []rune(s)
	return string(runes[:4]) + "**********" + string(runes[len(runes)-4:])
}

// 银行卡脱敏: 6222021234567890123 → 6222****0123
func maskBankCard(s string) string {
	if utf8.RuneCountInString(s) < 8 {
		return s
	}
	runes := []rune(s)
	return string(runes[:4]) + "****" + string(runes[len(runes)-4:])
}

// 姓名脱敏: 张三 → 张*
func maskName(s string) string {
	runes := []rune(s)
	switch len(runes) {
	case 0:
		return s
	case 1:
		return "*"
	case 2:
		return string(runes[0]) + "*"
	default:
		return string(runes[0]) + strings.Repeat("*", len(runes)-1)
	}
}

// 地址脱敏: 保留前6个字符
func maskAddress(s string) string {
	runes := []rune(s)
	if len(runes) <= 6 {
		return s
	}
	return string(runes[:6]) + "***"
}
