// decoder encoder test 

package libcode

import (
    "testing"
)


func TestDeBaseFunc(t *testing.T) {
    list := []int32{1, 2}
    v := DeBaseFunc(list)
    if v != 11 {
        t.Error("v is", v)
    }
}


func TestCheckPrefix(t *testing.T) {
    v := CheckPrefix(11)
    if v != true {
        t.Error("v is", v)
    }
    v = CheckPrefix(1)
    if v != false {
        t.Error("v is", v)
    }
}


func TestDecoder(t *testing.T) {
    Init()
    Decoder("民")
    Decoder("民主")
    Decoder("文化")
    Decoder("诚信民主民主法治诚信民主民主爱国诚信民主文明富强")
}


func TestUnMapCoreValue(t *testing.T) {
    Init()
    _, err := UnMapCoreValue("民主")
    if err != nil {
        t.Error("err is", err)
    }
    _, err = UnMapCoreValue("民主自由")
    if err != nil {
        t.Error("err is", err)
    }
}


func TestEncoder(t *testing.T) {
    Encoder("abc")
    Encoder("一")
    Encoder("屮")
}


func TestDecodeIndice(t *testing.T) {
    list := []int32{11,1,1,1,1,1,1,1}
    s, err := DecodeIndice(list)
    if err != ErrBadCoreValueStr {
        t.Error("expect error, got", s, err)
    }
}
