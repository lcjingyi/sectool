package pkg

import (
	"github.com/dromara/dongle"
)

func MyUrlEncode(str string) string {
	result := dongle.Encode.FromString(str).BySafeURL().ToString()
	return result
}

func MyUrlDncode(str string) string {
	result := dongle.Decode.FromString(str).BySafeURL().ToString()
	return result
}
