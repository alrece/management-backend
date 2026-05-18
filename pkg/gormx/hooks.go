package gormx

import (
	"reflect"

	"management-backend/pkg/crypto"

	"gorm.io/gorm"
)

var encryptionKey = []byte("32-byte-key-replace-via-config!")

// SetEncryptionKey 设置 AES-256 加密密钥（必须 32 字节）
func SetEncryptionKey(key []byte) {
	if len(key) == 32 {
		encryptionKey = make([]byte, 32)
		copy(encryptionKey, key)
	}
}

// EncryptedFields 注册 GORM 加解密 Hook
func EncryptedFields(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("encrypted:before_create", encryptHook)
	db.Callback().Query().After("gorm:query").Register("encrypted:after_query", decryptHook)
	db.Callback().Update().Before("gorm:update").Register("encrypted:before_update", encryptHook)
}

func encryptHook(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	val := reflect.ValueOf(db.Statement.Dest)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	encryptStructFields(val)
}

func decryptHook(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	val := reflect.ValueOf(db.Statement.Dest)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			if elem.Kind() == reflect.Struct {
				decryptStructFields(elem)
			}
		}
	case reflect.Ptr:
		if val.Elem().Kind() == reflect.Struct {
			decryptStructFields(val.Elem())
		}
	case reflect.Struct:
		decryptStructFields(val)
	}
}

func encryptStructFields(val reflect.Value) {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			encryptStructFields(val.Field(i))
			continue
		}
		if _, ok := field.Tag.Lookup("encrypted"); !ok {
			continue
		}
		f := val.Field(i)
		if f.Kind() != reflect.String || f.String() == "" {
			continue
		}
		encrypted, err := crypto.AESEncrypt([]byte(f.String()), encryptionKey)
		if err == nil {
			f.SetString(encrypted)
		}
	}
}

func decryptStructFields(val reflect.Value) {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			decryptStructFields(val.Field(i))
			continue
		}
		if _, ok := field.Tag.Lookup("encrypted"); !ok {
			continue
		}
		f := val.Field(i)
		if f.Kind() != reflect.String || f.String() == "" {
			continue
		}
		decrypted, err := crypto.AESDecrypt(f.String(), encryptionKey)
		if err == nil {
			f.SetString(string(decrypted))
		}
	}
}
