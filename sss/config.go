package sss

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type SSSConfig struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Local      string `json:"local"`
	LocalPort  int    `json:"local_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	Timeout    int    `json:"timeout"`
	ListFile   string `json:"list_file"`
}

var Config *SSSConfig

//解析配置文件
func ParseConfig(path string) (config *SSSConfig, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	config = &SSSConfig{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return
}