package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// 总配文件
type config struct {
	Db    db    `yaml:"db"`
	Redis redis `yaml:"redis"`
}

// 数据库的配置
type db struct {
	Dialects string `yaml:"dialects"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Db       string `yaml:"db"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
	MaxIdle  int    `yaml:"maxIdle"`
	MaxOpen  int    `yaml:"maxOpen"`
}

// redis配置
type redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

var Config *config

// 配置初始化
func init() {
	yamlFile, err := os.ReadFile("./config.yaml")
	//有错就down机
	if err != nil {
		panic(err)
	}
	//绑定值
	fmt.Println("看看打开了没")
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		panic(err)
	}
}
