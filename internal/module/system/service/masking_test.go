package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"13812345678", "138****5678"},
		{"123", "123"},
		{"", ""},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expect, maskPhone(tt.input))
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"test@example.com", "t***@example.com"},
		{"a@b.cn", "a***@b.cn"},
		{"@", "@"},
		{"", ""},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expect, maskEmail(tt.input))
	}
}

func TestMaskIDCard(t *testing.T) {
	assert.Equal(t, "3201**********1234", maskIDCard("320123199001011234"))
	assert.Equal(t, "1234", maskIDCard("1234"))
}

func TestMaskBankCard(t *testing.T) {
	assert.Equal(t, "6222****0123", maskBankCard("6222021234567890123"))
	assert.Equal(t, "1234567", maskBankCard("1234567"))
}

func TestMaskName(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"张三", "张*"},
		{"李", "*"},
		{"欧阳修", "欧**"},
		{"", ""},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expect, maskName(tt.input))
	}
}

func TestMaskAddress(t *testing.T) {
	assert.Equal(t, "北京市朝阳区***", maskAddress("北京市朝阳区建国门外大街1号"))
	assert.Equal(t, "北京市", maskAddress("北京市"))
}

func TestMaskStruct(t *testing.T) {
	type User struct {
		Name   string `masking:"name"`
		Phone  string `masking:"phone"`
		Email  string `masking:"email"`
		Normal string
	}

	u := User{
		Name:   "张三",
		Phone:  "13812345678",
		Email:  "test@example.com",
		Normal: "不变",
	}
	MaskStruct(&u)

	assert.Equal(t, "张*", u.Name)
	assert.Equal(t, "138****5678", u.Phone)
	assert.Equal(t, "t***@example.com", u.Email)
	assert.Equal(t, "不变", u.Normal)
}

func TestMaskStruct_NonStruct(t *testing.T) {
	// 非 struct 应不 panic
	MaskStruct("hello")
	MaskStruct(123)

	s := "hello"
	MaskStruct(&s) // *string 不是 struct，应跳过
}
