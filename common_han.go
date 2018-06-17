// use common han 2500 words

package libcode

import (
    "github.com/dilfish/tools"
    "unicode/utf8"
    "errors"
)


var commonHan map[rune]int
var revCommonHan map[int]rune
var hanIdx = 0
var errBadHanFile = errors.New("bad common han file")


func readCommon(w string) error {
    r, _ := utf8.DecodeRune([]byte(w))
    c := utf8.RuneCount([]byte(w))
    if c != 1 {
        return errBadHanFile
    }
    commonHan[r] = hanIdx
    revCommonHan[hanIdx] = r
    hanIdx ++
    return nil
}


func ReadCommon(fn string) error {
    commonHan = make(map[rune]int)
    revCommonHan = make(map[int]rune)
    return tools.ReadLine(fn, readCommon)
}


// prefix: 10
func EncodeCommonHan(code rune) int32 {
    idx, ok := commonHan[code]
    if ok == false {
        return -1
    }
    return int32(idx)
}


func DecodeCommonHan(off int32) rune {
    han, ok := revCommonHan[int(off)]
    if ok == false {
        return rune(-1)
    }
    return rune(han)
}


func InitCommonHan(fn string) error {
    err := RegisterEncoderDecoder(EncodeCommonHan, DecodeCommonHan, "common_han", Ring0)
    if err != nil {
        return err
    }
    err = ReadCommon(fn)
    if err != nil {
        return err
    }
    return nil
}
