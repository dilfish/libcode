// Copyright 2018 Sean.ZH

package libcode

// EncodeUnicde used the chinese unicode map as
// offset to map a list of word
// ref: www.unicode.org/versions/Uncide5.0.0/ch12.pdf
// range:
// 0x4e00-0x9fff
// 0x3400-0x4dff
// combined: 0x3400-0x9fff
func EncodeUnicode(r rune) int32 {
    if r < 0x3400 || r > 0x9fff {
        return BadCode
    }
    offset := r - 0x3400
    return int32(offset)
}


// DecodeUnicode read a list of code and 
// add offset to get it's unicode point
func DecodeUnicode(code int32) rune {
    code = code + 0x3400
    if code < 0x3400 || code > 0x9fff {
        return BadRune
    }
    return rune(code)
}
