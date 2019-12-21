package config

import (
	"gopkg.in/gcfg.v1"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/ebar-go/ego/helper"
	"os"
)


// Load 加载ini配置
func LoadIni(conf interface{}, filePath string) error {
	return  gcfg.ReadFileInto(conf, filePath)
}

// Load 加载yaml配置
func LoadYaml(conf interface{}, filePath string) error {
	yamlFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, conf)
}

// Load 加载json配置
func LoadJson(conf interface{}, filePath string) error {
	jsonBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}

	return helper.JsonDecode(jsonBytes, conf)
}

// 获取环境变量
func Getenv(name string) string {
	return os.Getenv(name)
}