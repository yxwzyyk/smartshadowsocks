package sss

import (
	"strconv"
	"net"
	"errors"
	"fmt"
	"os"
)

//打印系统版本
func PrintVersion() {
	const version = "1.0.0"
	fmt.Printf("SmartShadowSocks Version: %s\n", version)
}

//检测目录是否存在
func IsFileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode()&os.ModeType == 0 {
			return true, nil
		}
		return false, errors.New(path + "exists but is not regular file")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//
func JoinHostPort(host string, port int) (string) {
	return net.JoinHostPort(host, strconv.Itoa(int(port)))
}
