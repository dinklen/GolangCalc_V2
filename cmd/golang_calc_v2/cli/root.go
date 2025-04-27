package cli

import (
	"fmt"
	"os"

	"github.com/dinklen/GolangCalc_V2/internal/cli/commands"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func setupCommands() *cobra.Command {
	var root = &cobra.Command{
		Use:   "go_calc",
		Short: "Root command",
		Long: `The utility that allows you to calculate mathematical expressions
			by sending them to a local server, breaking them down into logical
			parts there and, if possible, calculating each of them in parallel
			on a microservice, while storing the history of calculations in a
			database.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use [go_calc help] for more info")
		},
	}

	root.AddCommand(
		commands.NewSignUpCommand(),
		commands.NewStartSessionCommand(),
		commands.NewLogoutCommand(),
		commands.NewHelpCommand(),
	)

	return root
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// uber log: print err
		os.Exit(1)
	}
}

func init() {
	rootCmd = setupCommands()
}
