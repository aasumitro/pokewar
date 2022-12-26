package appconfigs_test

import (
	"fmt"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppConfig(t *testing.T) {
	viper.Reset()
	viper.SetConfigFile("../../.example.env")
	viper.SetConfigType("dotenv")
	appconfigs.LoadEnv()

	tt := []struct {
		name     string
		value    any
		expected any
		reflect  any
	}{
		{
			name:     "Test AppName Env",
			value:    appconfigs.Instance.AppName,
			expected: "Pokewar",
		},
		{
			name:     "Test AppVersion Env",
			value:    appconfigs.Instance.AppVersion,
			expected: "0.0.2-dev",
		},
		{
			name:     "Test AppUrl Env",
			value:    appconfigs.Instance.AppURL,
			expected: "localhost:8000",
		},
		{
			name:     "Test PokeApiUrl Env",
			value:    appconfigs.Instance.PokeapiURL,
			expected: "https://pokeapi.co/api/v2/",
		},
		{
			name:     "Test DbDriver Env",
			value:    appconfigs.Instance.DBDriver,
			expected: "sqlite3",
		},
		{
			name:     "Test DbDsnUrl Env",
			value:    appconfigs.Instance.DBDsnURL,
			expected: "./db/local-data.db",
		},
		{
			name:     "TestUpdateEnv Function",
			expected: "UPDATE_SUCCESS",
		},
		{
			name:     "TestInitDBConn",
			expected: "DB_CONN",
		},
		{
			name:     "TestUpdateEnv Function ShouldError ReadWrite",
			expected: "UPDATE_ERROR",
		},
	}

	for _, test := range tt {
		fmt.Println(test.name)
		t.Run(test.name, func(t *testing.T) {
			switch test.expected {
			case "UPDATE_SUCCESS":
				initialValue := appconfigs.Instance.AppDebug
				appconfigs.Instance.UpdateEnv("APP_DEBUG", !initialValue)
				assert.Equal(t, appconfigs.Instance.AppDebug, false)
				appconfigs.Instance.UpdateEnv("APP_DEBUG", initialValue)
			case "DB_CONN":
				appconfigs.Instance.DBDsnURL = "../../db/local-data.db"
				appconfigs.Instance.InitDbConn()
				assert.NotEqual(t, appconfigs.DbPool, nil)
			case "UPDATE_ERROR":
				viper.Reset()
				initialValue := appconfigs.Instance.AppDebug
				appconfigs.Instance.UpdateEnv("APP_DEBUG", !initialValue)
				assert.Equal(t, appconfigs.Instance.AppDebug, initialValue)
			default:
				assert.Equal(t, test.expected, test.value)
			}
		})
	}
}
