// use unicode han code block

package libcode

// ref: www.unicode.org/versions/Uncide5.0.0/ch12.pdf
// range:
// 0x4e00-0x9fff
// 0x3400-0x4dff
// combined: 0x3400-0x9fff
// prefix: 11
func EncodeUnicode(code rune) int32 {
    if code < 0x3400 || code > 0x9fff {
        return -1
    }
    offset := code - 0x3400
    return int32(offset)
}


func DecodeUnicode(offset int32) rune {
    offset = offset + 0x3400
    if offset < 0x3400 || offset > 0x9fff {
        return rune(-1)
    }
    return rune(offset)
}


func init() {
    err := RegisterEncoderDecoder(EncodeUnicode, DecodeUnicode, "unicode", 5)
    if err != nil {
        panic(err)
    }
}
