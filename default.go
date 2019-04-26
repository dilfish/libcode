// Copyright 2018 Sean.ZH

package libcode

// EncodeDefault just see a word as it's unicode point number
func EncodeDefault(word rune) int32 {
	return int32(word)
}

// DecodeDefault just used the number as unicode point of a word
func DecodeDefault(off int32) rune {
	return rune(off)
}
