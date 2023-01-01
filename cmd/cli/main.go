package main

import (
	"fmt"
	"github.com/aasumitro/pokewar/cmd/cli/command"
	"github.com/aasumitro/pokewar/pkg/appconfig"
	"os"
)

func init() {
	appconfig.LoadEnv()

	appconfig.Instance.InitDbConn()
}

func main() {
	if err := command.CliCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
