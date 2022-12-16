package configs

import (
	"fmt"
	"log"
	"sync"
)

type AppConfig struct {
	AppName  string `mapstructure:"APP_NAME"`
	ApiUrl   string `mapstructure:"API_URL"`
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBDsnUrl string `mapstructure:"DB_DSN_URL"`
}

var (
	cfgOnce     sync.Once
	CfgInstance *AppConfig

	dbOnce sync.Once
	//dbInstance
)

func LoadEnv() {
	log.Printf("Load configuration file . . . .")

	cfgOnce.Do(func() {
		CfgInstance = nil
		// TODO LOAD ENV
	})

	fmt.Println(CfgInstance)
}

func (cfg *AppConfig) InitDb() {
	dbOnce.Do(func() {
		// TODO
	})
}