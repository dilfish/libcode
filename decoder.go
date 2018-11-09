// Copyright 2018 Sean.ZH

package libcode

import (
    "unicode/utf8"
    "errors"
    "github.com/dilfish/tools"
    "math"
)


// ErrBadCoreValueStr indicate an invalid string which could
// not resolved as list of code
var ErrBadCoreValueStr = errors.New("bad core value string")
// ErrTooManyEncoder indicate that you registerd more than 12 encoder
var ErrTooManyEncoder = errors.New("too many encoder")
// ErrReReg indicate that you register a registerd encoder
var ErrReReg = errors.New("re register")


// 12 words, default, unicode, common_han using 1 for each
// TOTAL_CV is 12 because Mr Xi only uses 12 words as core value
const TOTAL_CV = 12
// CHN_PREFIX used as chinese word
const CHN_PREFIX = 11
// UNI_PREFIX used as unicode point
const UNI_PREFIX = 10
// DEF_PREFIX used as default encoder
const DEF_PREFIX = 9
// BAD_PREFIX used as others
const BAD_PREFIX = -1
const modNum = int32(9)
// BadRune is defines in utf8 lib
const BadRune = utf8.RuneError
// BadCode is defined negative value
const BadCode = int32(-1)

func deBaseFunc(list []int32) rune {
    if list == nil || len(list) == 0 {
        return BadRune
    }
    v := int64(0)
    for _, l := range list {
        v = v * int64(modNum) + int64(l)
        if v > math.MaxInt32 {
            return BadRune
        }
    }
    return rune(v)
}


func (lc *LibCode) decodeWord(list []int32) rune {
    if len(list) == 0 {
        return 0
    }
    t := list[0]
    code := deBaseFunc(list[1:])
    if code == BadRune {
        return BadRune
    }
    switch t {
        case CHN_PREFIX:
            return lc.ch.DecodeCommonHan(code)
        case UNI_PREFIX:
            return DecodeUnicode(code)
        case DEF_PREFIX:
            return DecodeDefault(code)
    }
    return BadRune
}


func (lc *LibCode) decodeIndice(indice []int32) (string, error) {
    list := make([]int32, 0)
    orig := ""
    for _, index := range indice {
        if index != CHN_PREFIX && index != DEF_PREFIX && index != UNI_PREFIX {
            list = append(list, index)
            continue
        }
        o := lc.decodeWord(list)
        if o == BadRune {
            return "", ErrBadCoreValueStr
        }
        if o != 0 {
            orig = orig + string(o)
        }
        list = make([]int32, 0)
        list = append(list, index)
    }
    o := lc.decodeWord(list)
    if o == BadRune {
        return "", ErrBadCoreValueStr
    }
    if o != 0 {
        orig = orig + string(o)
    }
    return orig, nil
}


func (lc *LibCode) unMapCoreValue(cv string) ([]int32, error) {
    cvs := make([]string, 0)
    list := make([]int32, 0)
    for len(cv) > 0 {
        r, size := utf8.DecodeLastRuneInString(cv)
        cv = cv[:len(cv) - size]
        cvs = append([]string{string(r)}, cvs...)
    }
    if len(cvs) % 2 != 0 {
        return nil, ErrBadCoreValueStr
    }
    for i := 0;i < len(cvs); i = i + 2 {
        cvWord := cvs[i] + cvs[i + 1]
        idx, ok := lc.coreValueMap[cvWord]
        if ok == false {
            return nil, ErrBadCoreValueStr
        }
        list = append(list, idx)
    }
    return list, nil
}


// Decoder read a list of encoded string
// and send them to right decoder
func (lc *LibCode) Decoder(cv string) (string, error) {
    indice, err := lc.unMapCoreValue(cv)
    if err != nil {
        return "", err
    }
    orig, err := lc.decodeIndice(indice)
    if err != nil {
        return "", err
    }
    return orig, nil
}


func (lc *LibCode) readCoreValue(line string) error {
    if utf8.RuneCountInString(line) != 2 {
        return ErrBadCoreValueStr
    }
    lc.coreValueMap[line] = lc.idx
    lc.revCoreValueMap[lc.idx] = line
    lc.idx ++
    return nil
}


func baseFunc(index rune) []int32 {
    off := make([]int32, 0)
    for index > modNum {
        num := index % modNum
        index = index / modNum
        off = append([]int32{num}, off...)
    }
    num := index % modNum
    off = append([]int32{num}, off...)
    return off
}


func (lc *LibCode) getCode(r rune) (int32, int32) {
    code := lc.ch.EncodeCommonHan(r)
    if code != BadCode {
        return code, CHN_PREFIX
    }
    code = EncodeUnicode(r)
    if code != BadCode {
        return code, UNI_PREFIX
    }
    return EncodeDefault(r), DEF_PREFIX
}


func (lc *LibCode) getList(r rune) []int32 {
    code, prefix := lc.getCode(r)
    list := baseFunc(code)
    list = append([]int32{prefix}, list...)
    return list
}


// Encoder transform original message to core value message
func (lc *LibCode) Encoder(orig string) string {
    cv := ""
    for _, o := range orig {
        code := lc.getList(rune(o))
        for _, c := range code {
            cv = cv + lc.revCoreValueMap[c]
        }
    }
    return cv
}


// LibCode defines a list of encoder decoder pair
type LibCode struct {
    coreValueMap map[string]int32
    revCoreValueMap map[int32]string
    idx int32
    ch *CommonHanEncoder
}


// NewLibCode get an object
func NewLibCode(cv, ch string) (*LibCode, error) {
    lc := new(LibCode)
    lc.coreValueMap = make(map[string]int32)
    lc.revCoreValueMap = make(map[int32]string)
    err := tools.ReadLine(cv, lc.readCoreValue)
    if err != nil {
        return nil, err
    }
    if len(lc.coreValueMap) != TOTAL_CV {
        return nil, ErrBadCoreValueStr
    }
    lc.ch, err = NewCommonHan(ch)
    if err != nil {
        return nil, err
    }
    return lc, nil
}
