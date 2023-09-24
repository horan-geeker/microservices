package util

import (
	"crypto/md5"
	"fmt"
)

// MD5 md5 加密
func MD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
