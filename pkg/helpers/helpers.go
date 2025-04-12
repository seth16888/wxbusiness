package helpers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"
	"time"
)

var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// ValidatePhone 验证手机号
func ValidatePhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone number")
	}
	return nil
}

// EncryptPassword 加密密码
//
//	使用MD5加密，然后加上盐，十六进制字符串
func EncryptPassword(password, salt string) (string, error) {
	h := md5.New()
	_, e := io.WriteString(h, password+salt)
	if e != nil {
		return "", e
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, encryptedPassword, salt string) bool {
	encrypted, err := EncryptPassword(password, salt)
	if err != nil {
		return false
	}
	return encrypted == encryptedPassword
}

// TimenowInTimezone 获取当前时间
func TimenowInTimezone(tz *time.Location) time.Time {
	return time.Now().In(tz)
}

// IsImageFile 判断文件是否为图片
func IsImageFile(fileName string) bool {
	ext := fileName[strings.LastIndex(fileName, ".")+1:]
	ext = strings.ToLower(ext)
	imageExtensions := []string{"jpg", "jpeg", "png", "gif", "svg", "webp"}

	return slices.Contains(imageExtensions, ext)
}

// 从 数字 切片中获取最小值，返回错误处理
func FindMin[T int](numbers []T) (T, error) {
  if len(numbers) == 0 {
      return 0, errors.New("slice is empty")
  }

  var min T = numbers[0]
  for _, num := range numbers[1:] {
      if num < min {
          min = num
      }
  }
  return min, nil
}

// IsEmail 验证邮箱格式
func IsEmail(email string) bool {
  // 正则表达式匹配邮箱格式
  emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
  match, _ := regexp.MatchString(emailRegex, email)
  return match
}