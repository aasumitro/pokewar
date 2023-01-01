package appconfig_test

import (
	"github.com/aasumitro/pokewar/pkg/appconfig"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppConfig(t *testing.T) {
	viper.Reset()
	viper.SetConfigFile("../../.example.env")
	viper.SetConfigType("dotenv")
	appconfig.LoadEnv()

	tt := []struct {
		name     string
		value    any
		expected any
		reflect  any
	}{
		{
			name:     "Test AppName Env",
			value:    appconfig.Instance.AppName,
			expected: "Pokewar",
		},
		{
			name:     "Test AppVersion Env",
			value:    appconfig.Instance.AppVersion,
			expected: "0.0.2-dev",
		},
		{
			name:     "Test AppUrl Env",
			value:    appconfig.Instance.AppURL,
			expected: "localhost:8000",
		},
		{
			name:     "Test PokeApiUrl Env",
			value:    appconfig.Instance.PokeapiURL,
			expected: "https://pokeapi.co/api/v2/",
		},
		{
			name:     "Test DbDriver Env",
			value:    appconfig.Instance.DBDriver,
			expected: "sqlite3",
		},
		{
			name:     "Test DbDsnUrl Env",
			value:    appconfig.Instance.DBDsnURL,
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
		t.Run(test.name, func(t *testing.T) {
			switch test.expected {
			case "UPDATE_SUCCESS":
				initialValue := appconfig.Instance.AppDebug
				appconfig.Instance.UpdateEnv("APP_DEBUG", !initialValue)
				assert.Equal(t, appconfig.Instance.AppDebug, false)
				appconfig.Instance.UpdateEnv("APP_DEBUG", initialValue)
			case "DB_CONN":
				appconfig.Instance.DBDsnURL = "../../db/local-data.db"
				appconfig.Instance.InitDbConn()
				assert.NotEqual(t, appconfig.DbPool, nil)
			case "UPDATE_ERROR":
				viper.Reset()
				initialValue := appconfig.Instance.AppDebug
				appconfig.Instance.UpdateEnv("APP_DEBUG", !initialValue)
				assert.Equal(t, appconfig.Instance.AppDebug, initialValue)
			default:
				assert.Equal(t, test.expected, test.value)
			}
		})
	}
}
