package appconfigs_test

import (
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	viper.SetConfigFile("../../.example.env")

	appconfigs.LoadEnv()

	tests := []struct {
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
			name:     "Test AppDebug Env",
			value:    appconfigs.Instance.AppName,
			expected: "Pokewar",
		},
		{
			name:     "Test AppVersion Env",
			value:    appconfigs.Instance.AppVersion,
			expected: "0.0.1-dev",
		},
		{
			name:     "Test AppUrl Env",
			value:    appconfigs.Instance.AppUrl,
			expected: "localhost:8000",
		},
		{
			name:     "Test PokeApiUrl Env",
			value:    appconfigs.Instance.PokeapiUrl,
			expected: "https://pokeapi.co/api/v2/",
		},
		{
			name:     "Test Pokedex Env",
			value:    appconfigs.Instance.Pokedex,
			expected: "kanto",
		},
		{
			name:     "Test DbDriver Env",
			value:    appconfigs.Instance.DBDriver,
			expected: "sqlite3",
		},
		{
			name:     "Test DbDsnUrl Env",
			value:    appconfigs.Instance.DBDsnUrl,
			expected: "./db/local-data.db",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.value != test.expected {
				t.Errorf("expected %v, got %v", test.expected, test.value)
			}
		})
	}
}

// LOCAL ITS FINE
func TestUpdateEnv(t *testing.T) {
	t.Skip()
	viper.SetConfigFile("../../.example.env")

	appconfigs.LoadEnv()

	initialValue := appconfigs.Instance.AppDebug

	appconfigs.Instance.UpdateEnv("APP_DEBUG", !initialValue)

	appconfigs.LoadEnv()

	if appconfigs.Instance.AppDebug != false {
		t.Errorf("Expected APP_DEBUG to be false, got %v", appconfigs.Instance.AppDebug)
	}

	appconfigs.Instance.UpdateEnv("APP_DEBUG", initialValue)
}

func TestInitDBConn(t *testing.T) {
	viper.SetConfigFile("../../.example.env")

	appconfigs.LoadEnv()

	appconfigs.Instance.DBDsnUrl = "../../db/local-data.db"
	appconfigs.Instance.InitDbConn()

	if appconfigs.DbPool == nil {
		t.Errorf("Expected DbPool to be non-nil, got %v", appconfigs.DbPool)
	}
}
