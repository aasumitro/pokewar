package command

import (
	"github.com/spf13/cobra"
)

var pokemonCmd = &cobra.Command{
	Use:  "pokemon",
	Long: `pokemon cmd is used for main app: pokemon < list | rank | match >`,
}

var pokemonListCmd = &cobra.Command{
	Use:  "list",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

var pokemonRankCmd = &cobra.Command{
	Use:  "rank",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

var pokemonBattleCmd = &cobra.Command{
	Use:  "match",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

func init() {
	CliCmd.AddCommand(pokemonCmd)
	pokemonCmd.AddCommand(pokemonListCmd)
	pokemonCmd.AddCommand(pokemonRankCmd)
	pokemonCmd.AddCommand(pokemonBattleCmd)
}
