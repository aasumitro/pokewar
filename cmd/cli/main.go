package main

import (
	"fmt"
	"github.com/aasumitro/pokewar/cmd/cli/command"
	"github.com/aasumitro/pokewar/pkg/configs"
	"os"
)

func init() {
	configs.LoadEnv()

	configs.Instance.InitDbConn()
}

func main() {
	if err := command.CliCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
