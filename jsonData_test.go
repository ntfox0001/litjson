package litjson_test

import (
	"fmt"
	"testing"

	"plutus/utils/litjson"
)

type binaryStr struct {
	BB []byte
}

func Test1(t *testing.T) {
	bs := binaryStr{
		BB: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
	}
	buf, _ := litjson.Marshal(bs)
	fmt.Printf("js: %s", string(buf))
}

func Test2(t *testing.T) {
	m := make(map[string]string)
	m["aaa"] = "aaaa"

	maptest(m)

	fmt.Print(m, "\n")
}

func maptest(m map[string]string) {
	m["fff"] = "ffff"
}
