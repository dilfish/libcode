// +build gofuzz

package libcode


var lc *LibCode


func Fuzz(data []byte) int {
	lc.Decoder(string(data))
	return 0
}


func init() {
	var err error
	lc, err = NewLibCode("app/icved/core_values.txt", "app/icved/common_han.txt", 0)
	if err != nil {
		panic("new lib code error:" + err.Error())
	}
}
