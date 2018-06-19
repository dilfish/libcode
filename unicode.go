// use unicode han code block

package libcode


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


func DecodeUnicode(code int32) rune {
    code = code + 0x3400
    if code < 0x3400 || code > 0x9fff {
        return BadRune
    }
    return rune(code)
}
