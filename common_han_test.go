// test common han 2500 words

package libcode

import (
    "testing"
)


func TestEncodeCommonHan(t *testing.T) {
    idx := EncodeCommonHan(rune('一'))
    if idx != 0 {
        t.Error("Expect 1, got", idx)
    }
    idx = EncodeCommonHan(rune('屮'))
    if idx != -1 {
        t.Error("Expect -1, got", idx)
    }
}


func TestDecodeCommonHan(t *testing.T) {
    r := DecodeCommonHan(int32(0))
    if r != rune('一') {
        t.Error("expect 一, got", string(r), r)
    }
    r = DecodeCommonHan(2501)
    if r != rune(-1) {
        t.Error("expect -1, got", r, string(r))
    }
}
