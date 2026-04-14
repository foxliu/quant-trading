package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// TraderConfig 实盘配置结构体
type TraderConfig struct {
	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`

	Account struct {
		Name string `yaml:"name"`
	} `yaml:"account"`

	Strategies []struct {
		Name   string                 `yaml:"name"`
		Type   string                 `yaml:"type"`
		Params map[string]interface{} `yaml:"params"`
	} `yaml:"strategies"`

	Logging struct {
		Level      string `yaml:"level"`
		Filename   string `yaml:"filename"`
		MaxSize    int    `yaml:"max_size"`
		MaxAge     int    `yaml:"max_age"`
		MaxBackups int    `yaml:"max_backups"`
	} `yaml:"logging"`

	Mode string `yaml:"mode"`
}

// LoadTraderConfig 加载 trader.yaml（支持环境变量覆盖）
func LoadTraderConfig(filePath string) (*TraderConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%v", err)
	}

	cfg := &TraderConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败：%v", err)
	}

	return cfg, nil
}
