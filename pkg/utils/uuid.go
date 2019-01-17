package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func GenerateUUID(digits int, args ...interface{}) string {
	var inputFormat string
	paramCount := len(args)
	for i := 0; i < paramCount; i++ {
		inputFormat += "%v"
	}
	b := []byte(fmt.Sprintf(inputFormat, args...) + strconv.FormatInt(time.Now().Unix(), 10))
	hash := md5.New()
	hash.Write(b)
	return hex.EncodeToString(hash.Sum(nil))[:digits]
}
