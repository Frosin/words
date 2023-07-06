package service

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const (
	defaultConfigFile = "./service.json"
)

type ServiceConfig interface {
	DBTimeout() time.Duration
	GetDBFileName() string
	GetBotToken() string
	GetYadiskToken() string
	GetBotConfigFileName() string
}

// type FakeServiceCfg struct {
// }

type SConfig struct {
	DBTimeoutSec      int    `json:"db_timeout"`
	DBFileName        string `json:"db_filename"`
	BotConfigFileName string `json:"bot_config_filename"`
	BotToken          string `json:"bot_token"`
	YadiskToken       string `json:"yadisk_token"`
}

func NewServiceConfig(configFile string) *SConfig {
	sc := &SConfig{}
	if err := sc.initConfig(configFile); err != nil {
		log.Fatal(err)
	}
	return sc
}

func (sc *SConfig) GetDBFileName() string {
	return sc.DBFileName
}
func (sc *SConfig) GetBotToken() string {
	return sc.BotToken
}
func (sc *SConfig) GetYadiskToken() string {
	return sc.YadiskToken
}

func (sc *SConfig) DBTimeout() time.Duration {
	return time.Duration(sc.DBTimeoutSec) * time.Second
}

func (sc *SConfig) GetBotConfigFileName() string {
	return sc.BotConfigFileName
}

// func (f *FakeServiceCfg) DBTimeout() time.Duration {
// 	return time.Second
// }

func (c *SConfig) initConfig(cfgFile string) error {
	file := defaultConfigFile
	if cfgFile != "" {
		file = cfgFile
	}

	fileData, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileData, c)
	if err != nil {
		return nil
	}

	return nil
}
