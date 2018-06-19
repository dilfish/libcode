// decoder encoder test 

package libcode

import (
    "testing"
)


func TestDecoder(t *testing.T) {
    err := InitLibCode("icved/core_values.txt")
    if err != nil {
        t.Error("expect init good, got", err)
    }
    r, err := Decoder("友善法治爱国爱国")
    if err != nil {
        t.Error("exptect nil, got", err)
    }
    if r != "我" {
        t.Error("expect 我, got", []byte(r), []byte("我"))
    }
    r, err = Decoder("诚信民主平等民主法治富强")
    if err != nil {
        t.Error("expect nil, got", err)
    }
    if r != "屮" {
        t.Error("expect 屮, got", []byte(r), []byte("屮"))
    }
    r, err = Decoder("敬业民主民主法治敬业民主民主爱国敬业民主文明富强")
    if err != nil {
        t.Error("expect nil, got", err)
    }
    if r != "abc" {
        t.Error("expect abc, got", r)
    }
    _, err = Decoder("真")
    if err != ErrBadCoreValueStr {
        t.Error("expect ErrBadCoreValueStr, got", err)
    }
    _, err = Decoder("测试")
    if err != ErrBadCoreValueStr {
        t.Error("expect ErrBadCoreValueStr, got", err)
    }
}


func TestEncoder(t *testing.T) {
    err := InitLibCode("icved/core_values.txt")
    if err != nil {
        t.Error("expect init good, got", err)
    }
    Encoder("abc")
    Encoder("一")
    Encoder("屮")
}
