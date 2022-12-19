package appconfigs

import (
	"database/sql"
	"fmt"
	"github.com/fsnotify/fsnotify"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type AppConfig struct {
	AppName          string `mapstructure:"APP_NAME"`
	AppDebug         bool   `mapstructure:"APP_DEBUG"`
	AppVersion       string `mapstructure:"APP_VERSION"`
	AppUrl           string `mapstructure:"APP_URL"`
	PokeapiUrl       string `mapstructure:"POKEAPI_URL"`
	Pokedex          string `mapstructure:"POKEDEX"`
	DBDriver         string `mapstructure:"DB_DRIVER"`
	DBDsnUrl         string `mapstructure:"DB_DSN_URL"`
	LastSync         int    `mapstructure:"LAST_SYNC"`
	LastMonsterID    int    `mapstructure:"LAST_MONSTER_ID"`
	TotalMonsterSync int    `mapstructure:"TOTAL_MONSTER_SYNC"`
}

var (
	cfgOnce  sync.Once
	Instance *AppConfig

	dbOnce sync.Once
	DbPool *sql.DB
)

func init() {
	// set config file
	viper.SetConfigFile(".env")
}

func LoadEnv() {
	log.Printf("Load configuration file . . . .")
	// find environment file
	viper.AutomaticEnv()
	// read env handler
	readEnv := func() {
		// error handling for specific case
		if err := viper.ReadInConfig(); err != nil {
			// specified error message
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				panic(".env file not found!, please copy .example.env and paste as .env")
			}
			// general error message
			panic(fmt.Sprintf("ENV_ERROR: %s", err.Error()))
		}
		// extract config to struct
		if err := viper.Unmarshal(&Instance); err != nil {
			panic(fmt.Sprintf("ENV_ERROR: %s", err.Error()))
		}
	}
	// instance
	cfgOnce.Do(func() {
		// read env
		readEnv()
		// subs to event
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("update configuration data . . . .")

			readEnv()
		})
		// watch file update
		viper.WatchConfig()
		// notify that config file is ready
		log.Println("configuration file: ready")
	})
}

func (cfg *AppConfig) UpdateEnv(key, value any) {
	if err := viper.ReadInConfig(); err != nil {
		log.Println("READ", err.Error())
	}

	viper.Set(key.(string), value)

	viper.SetConfigType("dotenv")

	if err := viper.WriteConfig(); err != nil {
		log.Println("WRITE", err.Error())
	}
}

func (cfg *AppConfig) InitDbConn() {
	dbOnce.Do(func() {
		db, err := sql.Open(cfg.DBDriver, cfg.DBDsnUrl)
		if err != nil {
			panic(fmt.Sprintf("DATABASE_ERROR: %s", err.Error()))
		}

		DbPool = db

		if err := DbPool.Ping(); err != nil {
			panic(fmt.Sprintf("DATABASE_ERROR: %s", err.Error()))
		}

		log.Printf("Database connection pool created with %s driver . . . .", cfg.DBDriver)
	})
}
