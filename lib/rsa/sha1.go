package rsa

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

// HmacSha1对字符串进行加密
func HmacSHA1(keyStr, value string) string {
	mac := hmac.New(sha1.New, []byte(keyStr))
	mac.Write([]byte(value))
	res := hex.EncodeToString(mac.Sum(nil))
	return res
}
