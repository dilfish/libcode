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
var errNotInited = errors.New("not inited")


func _readCommon(w string) error {
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


func readCommon(fn string) error {
    commonHan = make(map[rune]int)
    revCommonHan = make(map[int]rune)
    return tools.ReadLine(fn, _readCommon)
}


func EncodeCommonHan(code rune) int32 {
    if len(commonHan) == 0 {
        panic(errNotInited)
    }
    idx, ok := commonHan[code]
    if ok == false {
        return BadCode
    }
    return int32(idx)
}


func DecodeCommonHan(off int32) rune {
    if len(revCommonHan) == 0 {
        panic(errNotInited)
    }
    han, ok := revCommonHan[int(off)]
    if ok == false {
        return BadRune
    }
    return rune(han)
}


func init() {
    err := readCommon("common_han.txt")
    if err != nil {
        panic(err)
    }
}
