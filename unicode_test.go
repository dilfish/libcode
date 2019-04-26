// test file for unicode

package libcode

import (
	"testing"
	"unicode/utf8"
)

func TestEncodeUnicode(t *testing.T) {
	table := []int32{
		0x33ff, -1,
		0xa000, -1,
		0x34ff, 0xff,
	}
	for idx, tb := range table {
		if idx%2 == 0 {
			r := EncodeUnicode(tb)
			if r != table[idx+1] {
				t.Error("expect", table[idx+1], "got", r)
			}
		}
	}
}

func TestDecodeUnicode(t *testing.T) {
	table := []int32{
		-1, utf8.RuneError,
		0x9fff, utf8.RuneError,
		0x20, 0x3400 + 0x20,
	}
	for idx, tb := range table {
		if idx%2 == 0 {
			r := DecodeUnicode(tb)
			if r != table[idx+1] {
				t.Error("expect", table[idx+1], "got", r)
			}
		}
	}
}
