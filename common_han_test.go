// Copyright 2018 Sean.ZH
// test common han 2500 words

package libcode

import (
	"errors"
	"testing"
	"unicode/utf8"
)

func TestLibcodeEncodeCommonHan(t *testing.T) {
	ch, err := NewCommonHan("app/icved/common_han.txt")
	if err != nil {
		t.Error("expect nil, got", err)
	}
	idx := ch.EncodeCommonHan(rune('一'))
	if idx != 0 {
		t.Error("Expect 1, got", idx)
	}
	idx = ch.EncodeCommonHan(rune('屮'))
	if idx != -1 {
		t.Error("Expect -1, got", idx)
	}
}

func TestLibcodeDecodeCommonHan(t *testing.T) {
	ch, err := NewCommonHan("app/icved/common_han.txt")
	if err != nil {
		t.Error("expect nil, got", err)
	}
	r := ch.DecodeCommonHan(int32(0))
	if r != rune('一') {
		t.Error("expect 一, got", string(r), r)
	}
	r = ch.DecodeCommonHan(2501)
	if r != utf8.RuneError {
		t.Error("expect -1, got", r, string(r))
	}
}

func TestLibcodeBadFile(t *testing.T) {
	_, err := NewCommonHan("common_han.go")
	if !errors.Is(err, errBadHanFile) {
		t.Error("expect errBadHanFile, got", err)
	}
}
