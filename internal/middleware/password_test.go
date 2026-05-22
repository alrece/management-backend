package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		pwd     string
		minLen  int
		wantErr bool
	}{
		{"有效密码", "Abc123!@", 8, false},
		{"长度不足", "Ab1!", 8, true},
		{"无大写", "abc123!@xyz", 8, true},
		{"无小写", "ABC123!@XYZ", 8, true},
		{"无数字", "Abcdefg!@", 8, true},
		{"无特殊字符", "Abcdefg12", 8, true},
		{"刚好满足", "Abc123!@", 8, false},
		{"空密码", "", 8, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.pwd, tt.minLen)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHashPassword_And_CheckPassword(t *testing.T) {
	pwd := "Admin123!@"
	hash, err := HashPassword(pwd)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	assert.True(t, CheckPassword(pwd, hash))
	assert.False(t, CheckPassword("wrongpwd", hash))
}
