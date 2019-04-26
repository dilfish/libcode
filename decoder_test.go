// Copyright 2018 Sean.ZH
// decoder encoder test

package libcode

import (
	"testing"
)

func TestDecoder(t *testing.T) {
	lc, err := NewLibCode("app/icved/core_values.txt", "app/icved/common_han.txt")
	if err != nil {
		t.Error("expect init good, got", err)
	}
	r, err := lc.Decoder("友善法治爱国爱国")
	if err != nil {
		t.Error("exptect nil, got", err)
	}
	if r != "我" {
		t.Error("expect 我, got", []byte(r), []byte("我"))
	}
	r, err = lc.Decoder("诚信民主平等民主法治富强")
	if err != nil {
		t.Error("expect nil, got", err)
	}
	if r != "屮" {
		t.Error("expect 屮, got", []byte(r), []byte("屮"))
	}
	r, err = lc.Decoder("敬业民主民主法治敬业民主民主爱国敬业民主文明富强")
	if err != nil {
		t.Error("expect nil, got", err)
	}
	if r != "abc" {
		t.Error("expect abc, got", r)
	}
	_, err = lc.Decoder("真")
	if err != ErrBadCoreValueStr {
		t.Error("expect ErrBadCoreValueStr, got", err)
	}
	_, err = lc.Decoder("测试")
	if err != ErrBadCoreValueStr {
		t.Error("expect ErrBadCoreValueStr, got", err)
	}
	_, err = lc.Decoder("民主民主民主")
	if err != ErrBadCoreValueStr {
		t.Error("expect ErrBadCoreValueStr, got", err)
	}
	r, err = lc.Decoder("诚信民主民主民主诚信")
	if err != ErrBadCoreValueStr {
		t.Error("expect ErrBadCoreValueStr, got", err, r)
	}
	r, err = lc.Decoder("友善")
	if err != ErrBadCoreValueStr {
		t.Error("expect ErrBadCoreValueStr, got", err, r)
	}
	str := "友善法治法治法治法治法治法治法治法治法治"
	str = str + "法治法治法治法治法治法治法治法治法治法治法治"
	str = str + "法治法治法治法治法治法治法治法治法治法治法治"
	str = str + "法治法治法治法治法治"
	r, err = lc.Decoder(str)
	if err != ErrBadCoreValueStr {
		t.Error("expect ErrBadCoreValueStr, got", err, r)
	}
	_, err = lc.Decoder("友善友善")
	if err != ErrBadCoreValueStr {
		t.Error("exptect nil, got", err)
	}
}

func TestEncoder(t *testing.T) {
	lc, err := NewLibCode("app/icved/core_values.txt", "app/icved/common_han.txt")
	if err != nil {
		t.Error("expect init good, got", err)
	}
	c := lc.Encoder("abc")
	if c != "敬业民主民主法治敬业民主民主爱国敬业民主文明富强" {
		t.Error("exptect long.., got", c)
	}
	c = lc.Encoder("屮")
	if c != "诚信民主平等民主法治富强" {
		t.Error("expect 诚信民主平等民主法治富强, got", c)
	}
	c = lc.Encoder("我")
	if c != "友善法治爱国爱国" {
		t.Error("expect 友善法治爱国爱国, got", c)
	}
}

func TestInit(t *testing.T) {
	_, err := NewLibCode("decoder.go", "")
	if err != ErrBadCoreValueStr {
		t.Error("expect ErrBadCoreValueStr, got", err)
	}
	_, err = NewLibCode("app/icved/core_values.txt", "common_han.go")
	if err != errBadHanFile {
		t.Error("expect errBadHanFile, got", err)
	}
	_, err = NewLibCode("testdata/core_values.txt", "")
	if err != ErrBadCoreValueStr {
		t.Error("expect ErrBadCoreValueStr, got", err)
	}
}
