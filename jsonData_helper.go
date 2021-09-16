package litjson

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func (jd *JSONData) GetFloat32ByString() float32 {
	if jd.ensure(TypeString) {
		f, err := strconv.ParseFloat(jd.data.(string), 32)
		if err != nil {
			return 0
		}

		return float32(f)
	}

	return 0
}

func (jd *JSONData) GetFloat64ByString() float64 {
	if jd.ensure(TypeString) {
		f, err := strconv.ParseFloat(jd.data.(string), 64)
		if err != nil {
			return 0
		}

		return f
	}

	return 0
}

func (jd *JSONData) GetInt32ByString() int32 {
	if jd.ensure(TypeString) {
		i, err := strconv.ParseInt(jd.data.(string), 10, 32)
		if err != nil {
			return 0
		}

		return int32(i)
	}

	return 0
}

func (jd *JSONData) GetIntByString() int {
	if jd.ensure(TypeString) {
		i, err := strconv.ParseInt(jd.data.(string), 10, 64)
		if err != nil {
			return 0
		}

		return int(i)
	}

	return 0
}

func (jd *JSONData) GetInt64ByString() int64 {
	if jd.ensure(TypeString) {
		i, err := strconv.ParseInt(jd.data.(string), 10, 64)
		if err != nil {
			return 0
		}

		return i
	}

	return 0
}

func (jd *JSONData) GetUInt32ByString() uint32 {
	if jd.ensure(TypeString) {
		i, err := strconv.ParseInt(jd.data.(string), 10, 32)
		if err != nil {
			return 0
		}

		return uint32(i)
	}

	return 0
}

func (jd *JSONData) GetUIntByString() uint {
	if jd.ensure(TypeString) {
		i, err := strconv.ParseUint(jd.data.(string), 10, 64)
		if err != nil {
			return 0
		}

		return uint(i)
	}

	return 0
}

func (jd *JSONData) GetUInt64ByString() uint64 {
	if jd.ensure(TypeString) {
		i, err := strconv.ParseUint(jd.data.(string), 10, 64)
		if err != nil {
			return 0
		}

		return i
	}

	return 0
}

func (jd *JSONData) SetNumberByString(key string, number interface{}) {
	jd.SetKey(key, fmt.Sprint(number))
}

func (jd *JSONData) SetBytes(key string, b []byte) {
	jd.SetKey(key, base64.StdEncoding.EncodeToString(b))
}

func (jd *JSONData) GetBytes(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(jd.Get(key).GetString())
}
