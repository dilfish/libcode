// encoder and decoder

package libcode

import (
    "unicode/utf8"
    "errors"
    "github.com/dilfish/tools"
)


var ErrBadCoreValueStr = errors.New("bad core value string")
var ErrTooManyEncoder = errors.New("too many encoder")
var ErrReReg = errors.New("re register")


var coreValueMap map[string]int32
var revCoreValueMap map[int32]string
var coreValueIdx = int32(0)

// 12 words, default, unicode, common_han using 1 for each
const CHN_PREFIX = 11
const UNI_PREFIX = 10
const DEF_PREFIX = 9
const modNum = int32(9)
const BadRune = rune(-1)
const BadCode = int32(-1)

func deBaseFunc(list []int32) rune {
    v := rune(0)
    for _, l := range list {
        v = v * rune(modNum) + rune(l)
    }
    return v
}


func decodeWord(list []int32) rune {
    if len(list) == 0 {
        return 0
    }
    t := list[0]
    code := deBaseFunc(list[1:])
    switch t {
        case CHN_PREFIX:
            return DecodeCommonHan(code)
        case UNI_PREFIX:
            return DecodeUnicode(code)
        case DEF_PREFIX:
            return DecodeDefault(code)
    }
    return BadRune
}


func decodeIndice(indice []int32) (string, error) {
    list := make([]int32, 0)
    orig := ""
    for _, index := range indice {
        if index != CHN_PREFIX && index != DEF_PREFIX && index != UNI_PREFIX {
            list = append(list, index)
            continue
        }
        o := decodeWord(list)
        if o == BadRune {
            return "", ErrBadCoreValueStr
        }
        if o != 0 {
            orig = orig + string(o)
        }
        list = make([]int32, 0)
        list = append(list, index)
    }
    o := decodeWord(list)
    if o == BadRune {
        return "", ErrBadCoreValueStr
    }
    if o != 0 {
        orig = orig + string(o)
    }
    return orig, nil
}


func unMapCoreValue(cv string) ([]int32, error) {
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
        idx, ok := coreValueMap[cvWord]
        if ok == false {
            return nil, ErrBadCoreValueStr
        }
        list = append(list, idx)
    }
    return list, nil
}


func Decoder(cv string) (string, error) {
    indice, err := unMapCoreValue(cv)
    if err != nil {
        return "", err
    }
    orig, err := decodeIndice(indice)
    if err != nil {
        return "", err
    }
    return orig, nil
}


func readCoreValue(line string) error {
    coreValueMap[line] = coreValueIdx
    revCoreValueMap[coreValueIdx] = line
    coreValueIdx ++
    return nil
}


func InitLibCode(cv string) error {
    coreValueMap = make(map[string]int32)
    revCoreValueMap = make(map[int32]string)
    return tools.ReadLine(cv, readCoreValue)
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


func getCode(r rune) (int32, int32) {
    code := EncodeCommonHan(r)
    if code != BadCode {
        return code, CHN_PREFIX
    }
    code = EncodeUnicode(r)
    if code != BadCode {
        return code, UNI_PREFIX
    }
    return EncodeDefault(r), DEF_PREFIX
}


// original word to code list
func getList(r rune) []int32 {
    code, prefix := getCode(r)
    if code == BadCode {
        return nil
    }
    list := baseFunc(code)
    list = append([]int32{prefix}, list...)
    return list
}


// transform original message to core value message
func Encoder(orig string) string {
    cv := ""
    for _, o := range orig {
        code := getList(rune(o))
        for _, c := range code {
            cv = cv + revCoreValueMap[c]
        }
    }
    return cv
}
