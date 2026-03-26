package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// TraderConfig 实盘配置结构体
type TraderConfig struct {
	CTP struct {
		FrontAddr  string `yaml:"front_addr"`
		BrokerID   string `yaml:"broker_id"`
		InvestorID string `yaml:"investor_id"`
		UserID     string `yaml:"user_id"`
		Password   string `yaml:"password"`
		AccountID  string `yaml:"account_id"`
	} `yaml:"ctp"`

	Account struct {
		InitialCash float64 `yaml:"initial_cash"`
	} `yaml:"account"`

	Strategies []struct {
		Name   string                 `yaml:"name"`
		Type   string                 `yaml:"type"`
		Params map[string]interface{} `yaml:"params"`
	} `yaml:"strategies"`

	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`
}

// LoadTraderConfig 加载 trader.yaml（支持环境变量覆盖）
func LoadTraderConfig() (*TraderConfig, error) {
	data, err := os.ReadFile("configs/trader.yaml")
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%v", err)
	}

	cfg := &TraderConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败：%v", err)
	}
	// 环境变量覆盖（优先级最高）
	if v := os.Getenv("CTP_FRONT_ADDR"); v != "" {
		cfg.CTP.FrontAddr = v
	}
	if v := os.Getenv("CTP_BROKER_ID"); v != "" {
		cfg.CTP.BrokerID = v
	}
	if v := os.Getenv("CTP_USER_ID"); v != "" {
		cfg.CTP.UserID = v
	}
	if v := os.Getenv("CTP_INVESTOR_ID"); v != "" {
		cfg.CTP.InvestorID = v
	}
	if v := os.Getenv("CTP_PASSWORD"); v != "" {
		cfg.CTP.Password = v
	}
	if v := os.Getenv("CTP_ACCOUNT_ID"); v != "" {
		cfg.CTP.AccountID = v
	}

	return cfg, nil
}
