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
    err := libcode.InitLibCode("core_values.txt")
    if err != nil {
        panic(err)
    }
    str := ""
    if *flagD == false {
        str = libcode.Encoder(*flagS)
    } else {
        str, err = libcode.Decoder(*flagS)
    }
    fmt.Println(str, err)
}
