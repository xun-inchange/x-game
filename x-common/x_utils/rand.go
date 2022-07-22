package x_utils

import "math/rand"

func RandUint32(min, max uint32) uint32 {
	return min + uint32(rand.Int31n(int32(max-min+1)))
}

func RandString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		x := rand.Intn(3)
		switch x {
		case 0:
			bytes[i] = byte(RandUint32(65, 90)) //大写字母
		case 1:
			bytes[i] = byte(RandUint32(97, 122)) //小写字母
		case 2:
			bytes[i] = byte(RandUint32(0, 9)) //数字
		}
	}
	return string(bytes)
}
