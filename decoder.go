// encoder and decoder

package libcode

import (
    "unicode/utf8"
    "errors"
    "github.com/dilfish/tools"
)


var ErrBadCoreValueStr = errors.New("bad core value string")
var ErrTooManyEncoder = errors.New("too many encoder")


var CoreValueMap map[string]int32
var RevCoreValueMap map[int32]string
var coreValueIdx = int32(0)
var baseNum = int32(12)

type fEncoder func(rune)int32
type fDecoder func(int32)rune
type EncoderDecoder struct {
    En fEncoder
    De fDecoder
    Priority int
    Prefix int32
}


var enDeMap map[string]EncoderDecoder

const (
    BadIndex = -1
    BadPrefix = -1
    Ring0 = 0
    Ring5 = 5
    Ring10 = 10
    RingNobody = 11
)


func DeBaseFunc(list []int32) rune {
    v := rune(0)
    for _, l := range list {
        v = v * rune(baseNum) + rune(l)
    }
    return v
}


func DecodeWord(list []int32) rune {
    if len(list) == 0 {
        return 0
    }
    t := list[0]
    idx := DeBaseFunc(list[1:])
    currPri := RingNobody
    currRune := rune(-1)
    for _, ed := range enDeMap {
        if ed.Prefix == t {
            if ed.Priority < currPri {
                currRune = ed.De(idx)
            }
        }
    }
    return currRune
}


func CheckPrefix(prefix int32) bool {
    for _, ed := range enDeMap {
        if ed.Prefix == prefix {
            return true
        }
    }
    return false
}


func DecodeIndice(indice []int32) (string, error) {
    list := make([]int32, 0)
    orig := ""
    for _, index := range indice {
        if CheckPrefix(index) == false {
            list = append(list, index)
            continue
        }
        o := DecodeWord(list)
        if o == -1 {
            return "", ErrBadCoreValueStr
        }
        orig = orig + string(o)
        list = make([]int32, 0)
        list = append(list, index)
    }
    o := DecodeWord(list)
    if o == -1 {
        return "", ErrBadCoreValueStr
    }
    orig = orig + string(o)
    return orig, nil
}


func UnMapCoreValue(cv string) ([]int32, error) {
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
        idx, ok := CoreValueMap[cvWord]
        if ok == false {
            return nil, ErrBadCoreValueStr
        }
        list = append(list, idx)
    }
    return list, nil
}


func Decoder(cv string) (string, error) {
    indice, err := UnMapCoreValue(cv)
    if err != nil {
        return "", err
    }
    orig, err := DecodeIndice(indice)
    if err != nil {
        return "", err
    }
    return orig, nil
}


func RegisterEncoderDecoder(e fEncoder, d fDecoder, name string, pri int) error {
    if enDeMap == nil {
        enDeMap = make(map[string]EncoderDecoder)
    }
    baseNum = baseNum - 1
    if baseNum < 0 {
        return ErrTooManyEncoder
    }
    var ed EncoderDecoder
    ed.En = e
    ed.De = d
    ed.Priority = pri
    ed.Prefix = baseNum
    enDeMap[name] = ed
    return nil
}


func readCoreValue(line string) error {
    CoreValueMap[line] = coreValueIdx
    RevCoreValueMap[coreValueIdx] = line
    coreValueIdx ++
    return nil
}


func Init(cv, ch string) error {
    CoreValueMap = make(map[string]int32)
    RevCoreValueMap = make(map[int32]string)
    err := InitDefault()
    if err != nil {
        return err
    }
    err = InitCommonHan(ch)
    if err != nil {
        return err
    }
    return tools.ReadLine(cv, readCoreValue)
}


func BaseFunc(index rune) []int32 {
    off := make([]int32, 0)
    for index > baseNum {
        num := index % baseNum
        index = index / baseNum
        off = append([]int32{num}, off...)
    }
    num := index % baseNum
    off = append([]int32{num}, off...)
    return off
}


// original word to code list
func GetCode(o rune) []int32 {
    currIdx := int32(BadIndex)
    currPriority := RingNobody
    currPrefix := int32(BadPrefix)
    for _, ed := range enDeMap {
        idx := ed.En(o)
        if idx < 0 {
            continue
        }
        if ed.Priority < currPriority {
            currIdx = idx
            currPriority = ed.Priority
            currPrefix = ed.Prefix
        }
    }
    list := BaseFunc(currIdx)
    list = append([]int32{currPrefix}, list...)
    return list
}


// transform original message to core value message
func Encoder(orig string) string {
    cv := ""
    for _, o := range orig {
        code := GetCode(rune(o))
        for _, c := range code {
            cv = cv + RevCoreValueMap[c]
        }
    }
    return cv
}
