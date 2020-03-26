package security

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

//Md5哈希
func Md5(data string) string {
	m := md5.New()
	io.WriteString(m, data)
	return strings.ToLower(hex.EncodeToString(m.Sum(nil)))
}
func Md5d(data string) string {
	m := md5.New()
	io.WriteString(m, data)
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

//Sha1哈希
func Sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//Sha256哈希
func Sha256(data string) string {
	t := sha256.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//Hmac Sha1哈希
func HmacSha1(data string, key string, args ...bool) string {
	resultString := ""
	isHex := true

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))

	if len(args) > 0 {
		isHex = args[0]
	}

	if isHex {
		resultString = hex.EncodeToString(mac.Sum(nil))
	} else {
		resultString = string(mac.Sum(nil))
	}
	return resultString
}

//Hmac Sha256哈希
func HmacSha256(data string, key string, args ...bool) string {
	resultString := ""
	isHex := true

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))

	if len(args) > 0 {
		isHex = args[0]
	}

	if isHex {
		resultString = hex.EncodeToString(mac.Sum(nil))
	} else {
		resultString = string(mac.Sum(nil))
	}
	return resultString
}
