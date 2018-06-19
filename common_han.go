// use common han 2500 words

package libcode

import (
    "github.com/dilfish/tools"
    "unicode/utf8"
    "errors"
)


type CommonHanEncoder struct {
    commonHan map[rune]int
    revCommonHan map[int]rune
    idx int
}


var errBadHanFile = errors.New("bad common han file")


func (chen *CommonHanEncoder)_readCommon(w string) error {
    r, _ := utf8.DecodeRune([]byte(w))
    c := utf8.RuneCount([]byte(w))
    if c != 1 {
        return errBadHanFile
    }
    chen.commonHan[r] = chen.idx
    chen.revCommonHan[chen.idx] = r
    chen.idx ++
    return nil
}


func (chen *CommonHanEncoder) readCommon(fn string) error {
    chen.commonHan = make(map[rune]int)
    chen.revCommonHan = make(map[int]rune)
    return tools.ReadLine(fn, chen._readCommon)
}


func (chen *CommonHanEncoder) EncodeCommonHan(code rune) int32 {
    idx, ok := chen.commonHan[code]
    if ok == false {
        return BadCode
    }
    return int32(idx)
}


func (chen *CommonHanEncoder) DecodeCommonHan(off int32) rune {
    han, ok := chen.revCommonHan[int(off)]
    if ok == false {
        return BadRune
    }
    return rune(han)
}


func NewCommonHan(fn string) (*CommonHanEncoder, error) {
    chen := new(CommonHanEncoder)
    err := chen.readCommon(fn)
    if err != nil {
        return nil, err
    }
    return chen, nil
}
