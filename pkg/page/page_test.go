package page

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReq_GetPage(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{-1, 1},
		{1, 1},
		{5, 5},
		{100, 100},
	}
	for _, tt := range tests {
		r := &Req{Page: tt.input}
		assert.Equal(t, tt.expected, r.GetPage(), "Page=%d", tt.input)
	}
}

func TestReq_GetPageSize(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 10},
		{-1, 10},
		{10, 10},
		{50, 50},
		{100, 100},
	}
	for _, tt := range tests {
		r := &Req{PageSize: tt.input}
		assert.Equal(t, tt.expected, r.GetPageSize(), "PageSize=%d", tt.input)
	}
}

func TestReq_Offset(t *testing.T) {
	tests := []struct {
		page     int
		pageSize int
		offset   int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 20, 40},
		{0, 0, 0},   // 默认 page=1, pageSize=10
		{5, 0, 40},  // 默认 pageSize=10 → offset=40
	}
	for _, tt := range tests {
		r := &Req{Page: tt.page, PageSize: tt.pageSize}
		assert.Equal(t, tt.offset, r.Offset(), "page=%d,pageSize=%d", tt.page, tt.pageSize)
	}
}

func TestReq_Limit(t *testing.T) {
	r := &Req{PageSize: 25}
	assert.Equal(t, 25, r.Limit())

	r = &Req{}
	assert.Equal(t, 10, r.Limit())
}
