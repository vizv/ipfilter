package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func CalculateMD5(bytes []byte) string {
	md5 := md5.New()
	md5.Write(bytes)
	md5Hex := hex.EncodeToString(md5.Sum(nil))

	return md5Hex
}
