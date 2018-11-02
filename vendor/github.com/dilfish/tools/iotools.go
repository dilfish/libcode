package tools

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

var ErrBadFmt = errors.New("bad format")
var ErrNoSuch = errors.New("no such")
var ErrDupData = errors.New("dup data")

func ReadConfig(fn string, conf interface{}) error {
	bt, err := ReadFile(fn)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, conf)
}

func RandInt(w int) int32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int31n(int32(w))
}

func RandStr(w int) string {
	rand.Seed(time.Now().UnixNano())
	base := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	str := ""
	for i := 0; i < w; i++ {
		idx := rand.Int31n(int32(len(base)))
		str = str + string(base[idx])
	}
	return str
}

func ReadFile(fn string) ([]byte, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}


type LineFunc func(line string) error

func ReadLineArr(fn string, lf LineFunc, split int) error {
	return readLine(fn, lf, split)
}

func ReadLine(fn string, lf LineFunc) error {
	return readLine(fn, lf, 0)
}

func readLine(fn string, lf LineFunc, split int) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		if line == "" {
			continue
		}
		line = line[:len(line)-1]
		if split != 0 {
			arr := strings.Split(line, " ")
			if len(arr) != split {
				return ErrBadFmt
			}
		}
		err = lf(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func FileMd5(fn string) (int, string, error) {
	file, err := os.Open(fn)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()
	bt, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, "", err
	}
	return len(bt), fmt.Sprintf("%x", md5.Sum(bt)), nil
}

func UnixToBJ(unix int64) time.Time {
	return unixTo(unix, "Asia/Shanghai")
}

func UnixToUSPacific(unix int64) time.Time {
	return unixTo(unix, "US/Pacific")
}

func UnixToUTC(unix int64) time.Time {
	return unixTo(unix, "UTC")
}

func unixTo(unix int64, name string) time.Time {
	l, err := time.LoadLocation(name)
	if err != nil {
		panic("bad time name : " + name)
	}
	t := time.Unix(unix, 0)
	return t.In(l)
}

func TimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
