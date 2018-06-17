// use utf8 encode, default implemention

package libcode


func EncodeDefault(word rune) int32 {
    return int32(word)
}


func DecodeDefault(off int32) rune {
    return rune(off)
}


func InitDefault() error {
    return RegisterEncoderDecoder(EncodeDefault, DecodeDefault, "default", Ring10)
}
