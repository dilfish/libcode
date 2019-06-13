// Copyright 2018 Sean.ZH

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
	"net/http"
)

// ErrBadFmt bad format in file
var ErrBadFmt = errors.New("bad format")
// ErrNoSuch not exists
var ErrNoSuch = errors.New("no such")
// ErrDupData indicate duplicate data
var ErrDupData = errors.New("dup data")

// ReadConfig read file to interface
func ReadConfig(fn string, conf interface{}) error {
	bt, err := ReadFile(fn)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, conf)
}

// RandInt generate a int limited to w
func RandInt(w int) int32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int31n(int32(w))
}

// RandStr return random string length w
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

// ReadFile read all file content as byte
func ReadFile(fn string) ([]byte, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}


// LineFunc read a file into lines
type LineFunc func(line string) error
// ReadLine read line and calls back lf
func ReadLine(fn string, lf LineFunc) error {
    file, err := os.Open(fn)
    if err != nil {
        return err
    }
    defer file.Close()
	return readLine(file, lf, 0)
}


// GetLine get http content as file and calls lf on every line
func GetLine(url string, lf LineFunc) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return readLine(resp.Body, lf, 0)
}

func readLine(reader io.Reader, lf LineFunc, split int) error {
	rd := bufio.NewReader(reader)
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

// FileMd5 calc file's md5
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

// UnixToBJ unix timestamp to beijing time
func UnixToBJ(unix int64) time.Time {
	return unixTo(unix, "Asia/Shanghai")
}

// UnixToUSPacific unix timestamp to the U.S time
func UnixToUSPacific(unix int64) time.Time {
	return unixTo(unix, "US/Pacific")
}

// UnixToUTC unix timestamp to utc time
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

// TimeStr return standard time string
func TimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}


func dfsCallback(dir string, ret map[string]error, cb DFSCallback) {
	file, err := os.Open(dir)
	if err != nil {
		ret[dir] = err
		return
	}
	defer file.Close()
	fi, err := file.Readdir(-1)
	if err != nil {
		ret[dir] = err
		return
	}
	for _, f := range fi {
		if f.IsDir() == false {
			err = cb(dir + "/" + f.Name())
			if err != nil {
				ret[dir + "/" + f.Name()] = err
			}
		} else {
			dfsCallback(dir + "/" + f.Name(), ret, cb)
		}
	}
}


// DFSIter find every file at a dir
// and calls cb for every file
type DFSCallback func (string) error
func DFSIter(dir string, cb DFSCallback) map[string]error {
	ret := make(map[string]error)
	dfsCallback(dir, ret, cb)
	return ret
}
