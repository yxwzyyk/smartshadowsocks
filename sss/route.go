package sss

import (
	"bufio"
	"errors"
	"github.com/ryanuber/go-glob"
	"io"
	"os"
	"strings"
)

var Routebuffer []string

func SetRouteBuffer(fileName string) error {
	exists, err := IsFileExists(fileName)
	if !exists || err != nil {
		return errors.New("FileName not found!")
	}
	Routebuffer = readLine(fileName)
	return nil
}

func RouteMatch(s string) bool {
	for _, str := range Routebuffer {
		if r := glob.Glob(str, s); r == true {
			return true
		}
	}
	return false
}

func readLine(fn string) (b []string) {
	file, err := os.Open(fn)
	if err != nil {
		return nil
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')

		line = strings.Replace(line, " ", "", -1)  //去除空格
		line = strings.Replace(line, "\n", "", -1) //去除换行符
		line = strings.Replace(line, "\t", "", -1) //去除制表符

		if len(line) != 0 {
			b = append(b, line)
		}
		if err == io.EOF {
			break
		}
	}
	return b
}
