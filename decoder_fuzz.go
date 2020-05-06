// +build gofuzz

package libcode


var lc *LibCode


func Fuzz(input string) int {
	lc.Decoder(input)
	return 0
}


func init() {
	var err error
	lc, err = NewLibCode("app/icved/core_values.txt", "app/icved/common_han.txt", 0)
	if err != nil {
		panic("new lib code error:" + err.Error())
	}
}
