package main

import (
    "flag"
    "github.com/dilfish/libcode"
    "fmt"
)


var flagD = flag.Bool("d", false, "decode message")
var flagS = flag.String("str", "", "the message")
var flagH = flag.Bool("h", false, "help message")


func Help() {
    flag.PrintDefaults()
}


func main() {
    flag.Parse()
    if *flagS == "" || *flagH == true {
        Help()
        return
    }
    lc, err := libcode.NewLibCode("core_values.txt", "common_han.txt")
    if err != nil {
        panic(err)
    }
    str := ""
    if *flagD == false {
        str = lc.Encoder(*flagS)
    } else {
        str, err = lc.Decoder(*flagS)
    }
    fmt.Println(str, err)
}
