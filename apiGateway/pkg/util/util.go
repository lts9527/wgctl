package util

import (
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

// GetBetweenStr 截取字符串
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

// RandString 生成随机字符串
func RandString(len int) string {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().Unix()))
	bs := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bs[i] = byte(b)
	}
	return string(bs)
}

// 方法二
func RandStr2(n int) string {
	result := make([]byte, n/2)
	rand.Seed(time.Now().UnixNano())
	rand.Read(result)
	return hex.EncodeToString(result)
}

// RandStringlowercase 生成随机小写字符串
func RandStringlowercase(len int) string {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().Unix()))
	bs := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bs[i] = byte(b)
	}
	return strings.ToLower(string(bs))
}
