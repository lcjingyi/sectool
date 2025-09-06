package hashCalcuation

import (
	"github.com/dromara/dongle"
)

func MyMd5(str string) string {
	hashObject := dongle.Encrypt.FromString(str).ByMd5().ToHexString()
	return hashObject
}
