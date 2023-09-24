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

	IsWorkerDebugOnce() bool
	IsWorkerDebugUser() string
}

type SConfig struct {
	DBTimeoutSec      int    `json:"db_timeout"`
	DBFileName        string `json:"db_filename"`
	BotConfigFileName string `json:"bot_config_filename"`
	BotToken          string `json:"bot_token"`
	YadiskToken       string `json:"yadisk_token"`

	DebugWorkerUser string `json:"debug_worker_user"`
	DebugWorkerOnce bool   `json:"debug_worker_once"`
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

func (sc *SConfig) initConfig(cfgFile string) error {
	file := defaultConfigFile
	if cfgFile != "" {
		file = cfgFile
	}

	fileData, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileData, sc)
	if err != nil {
		return nil
	}

	return nil
}

func (sc *SConfig) IsWorkerDebugOnce() bool {
	return sc.DebugWorkerOnce
}

func (sc *SConfig) IsWorkerDebugUser() string {
	return sc.DebugWorkerUser
}
