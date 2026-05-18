package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const captchaChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// CaptchaService 验证码服务
type CaptchaService struct {
	rdb redis.Cmdable
}

// NewCaptchaService 创建验证码服务
func NewCaptchaService(rdb redis.Cmdable) *CaptchaService {
	return &CaptchaService{rdb: rdb}
}

// CaptchaResult 验证码结果
type CaptchaResult struct {
	CaptchaKey   string `json:"captchaKey"`
	CaptchaImage string `json:"captchaImage"`
}

// Generate 生成图形验证码
func (s *CaptchaService) Generate(ctx context.Context) (*CaptchaResult, error) {
	code := randomCode(4)
	img := generateImage(code, 120, 40)

	key := randomCode(16)
	redisKey := fmt.Sprintf("captcha:%s", key)
	if err := s.rdb.Set(ctx, redisKey, code, 5*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("验证码存储失败: %w", err)
	}

	return &CaptchaResult{
		CaptchaKey:   key,
		CaptchaImage: encodeImageBase64(img),
	}, nil
}

// Verify 校验验证码（一次性）
func (s *CaptchaService) Verify(ctx context.Context, key, code string) bool {
	redisKey := fmt.Sprintf("captcha:%s", key)
	stored, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return false
	}
	s.rdb.Del(ctx, redisKey)
	return stored == code
}

func randomCode(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(captchaChars))))
		result[i] = captchaChars[n.Int64()]
	}
	return string(result)
}

func generateImage(code string, w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{245, 245, 245, 255}}, image.Point{}, draw.Src)

	// 干扰点
	for i := 0; i < 100; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(w)))
		y, _ := rand.Int(rand.Reader, big.NewInt(int64(h)))
		img.Set(int(x.Int64()), int(y.Int64()), color.RGBA{180, 180, 180, 255})
	}

	// 绘制字符
	spacing := w / (len(code) + 1)
	for i, c := range code {
		x := spacing*(i+1) - 6
		for row := 0; row < 5; row++ {
			for col := 0; col < 5; col++ {
				if charDot(byte(c), row, col) {
					px := x + col
					py := h/2 - 4 + row
					if px >= 0 && px < w && py >= 0 && py < h {
						img.Set(px, py, color.RGBA{40, 40, 120, 255})
					}
				}
			}
		}
	}

	return img
}

// charDot 简单 5x5 像素字体
func charDot(c byte, row, col int) bool {
	fonts := map[byte][]string{
		'A': {"01110", "10001", "11111", "10001", "10001"},
		'B': {"11110", "10001", "11110", "10001", "11110"},
		'C': {"01111", "10000", "10000", "10000", "01111"},
		'D': {"11110", "10001", "10001", "10001", "11110"},
		'E': {"11111", "10000", "11110", "10000", "11111"},
		'F': {"11111", "10000", "11110", "10000", "10000"},
		'G': {"01111", "10000", "10011", "10001", "01110"},
		'H': {"10001", "10001", "11111", "10001", "10001"},
		'J': {"00111", "00010", "00010", "10010", "01100"},
		'K': {"10001", "10010", "11100", "10010", "10001"},
		'L': {"10000", "10000", "10000", "10000", "11111"},
		'M': {"10001", "11011", "10101", "10001", "10001"},
		'N': {"10001", "11001", "10101", "10011", "10001"},
		'P': {"11110", "10001", "11110", "10000", "10000"},
		'Q': {"01110", "10001", "10001", "10010", "01101"},
		'R': {"11110", "10001", "11110", "10010", "10001"},
		'S': {"01111", "10000", "01110", "00001", "11110"},
		'T': {"11111", "00100", "00100", "00100", "00100"},
		'U': {"10001", "10001", "10001", "10001", "01110"},
		'V': {"10001", "10001", "10001", "01010", "00100"},
		'W': {"10001", "10001", "10101", "11011", "10001"},
		'X': {"10001", "01010", "00100", "01010", "10001"},
		'Y': {"10001", "01010", "00100", "00100", "00100"},
		'Z': {"11111", "00010", "00100", "01000", "11111"},
		'2': {"01110", "10001", "00110", "01000", "11111"},
		'3': {"11110", "00001", "01110", "00001", "11110"},
		'4': {"10010", "10010", "11111", "00010", "00010"},
		'5': {"11111", "10000", "11110", "00001", "11110"},
		'6': {"01110", "10000", "11110", "10001", "01110"},
		'7': {"11111", "00010", "00100", "01000", "01000"},
		'8': {"01110", "10001", "01110", "10001", "01110"},
		'9': {"01110", "10001", "01111", "00001", "01110"},
	}
	pattern, ok := fonts[c]
	if !ok || row >= len(pattern) || col >= len(pattern[row]) {
		return false
	}
	return pattern[row][col] == '1'
}

func encodeImageBase64(img image.Image) string {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}
