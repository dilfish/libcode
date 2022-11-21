// Copyright 2018 Sean.ZH

package libcode

import (
	"errors"
	"log"
	"unicode/utf8"

	dio "github.com/dilfish/tools/io"
)

// CommonHanEncoder is the least priority encoder for libcode
// it just map a word's unicode point to a int32 value
type CommonHanEncoder struct {
	commonHan    map[rune]int
	revCommonHan map[int]rune
	idx          int
}

var errBadHanFile = errors.New("bad common han file")

func (chen *CommonHanEncoder) _readCommon(w string) error {
	r, _ := utf8.DecodeRune([]byte(w))
	c := utf8.RuneCount([]byte(w))
	if c != 1 {
		log.Println("read common error:", c)
		return errBadHanFile
	}
	chen.commonHan[r] = chen.idx
	chen.revCommonHan[chen.idx] = r
	chen.idx++
	return nil
}

func (chen *CommonHanEncoder) readCommon(fn string) error {
	chen.commonHan = make(map[rune]int)
	chen.revCommonHan = make(map[int]rune)
	return dio.ReadLine(fn, chen._readCommon)
}

// EncodeCommonHan encodes a word to int32 as it's unicode point
// using the comon 2500 Chinese characters
func (chen *CommonHanEncoder) EncodeCommonHan(code rune) int32 {
	idx, ok := chen.commonHan[code]
	if ok == false {
		log.Println("not a valid common han:", code)
		return BadCode
	}
	return int32(idx)
}

// DecodeCommonHan read a int32 as unicode point and map it
// to a char
func (chen *CommonHanEncoder) DecodeCommonHan(off int32) rune {
	han, ok := chen.revCommonHan[int(off)]
	if ok == false {
		log.Println("not a valid unicode:", off)
		return BadRune
	}
	return rune(han)
}

// NewCommonHan provide a new service
func NewCommonHan(fn string) (*CommonHanEncoder, error) {
	chen := new(CommonHanEncoder)
	err := chen.readCommon(fn)
	if err != nil {
		log.Println("read common error:", err)
		return nil, err
	}
	return chen, nil
}
