package pkg

import (
	"github.com/dromara/dongle"
)

func AesEncode(key, iv, str string) string {
	aesEncode := dongle.NewCipher()
	aesEncode.SetMode(dongle.CBC)
	aesEncode.SetKey(key)
	aesEncode.SetIV(iv)
	aseStr := dongle.Encrypt.FromString(str).ByAes(aesEncode).String()
	return aseStr
}

func AesDecode(key, iv, str string) string {
	aesDecode := dongle.NewCipher()
	aesDecode.SetMode(dongle.CBC)
	aesDecode.SetKey(key)
	aesDecode.SetIV(iv)
	aesStr := dongle.Decrypt.FromHexString(str).ByAes(aesDecode).String()
	return aesStr
}
