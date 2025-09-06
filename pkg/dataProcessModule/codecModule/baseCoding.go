package codecModule

import (
	"github.com/dromara/dongle"
)

func B64Encode(str string) string {
	//dongle.Encode.FromString(str string).加密类型.输出形式
	result := dongle.Encode.FromString(str).ByBase64().String()
	return result
}

func B64Decode(str string) string {
	result := dongle.Decode.FromString(str).ByBase64().ToString()
	return result
}

func B32Eecode(str string) string {
	result := dongle.Encode.FromString(str).ByBase32().String()
	return result
}

func B32Decode(str string) string {
	result := dongle.Decode.FromString(str).ByBase32().ToString()
	return result
}
