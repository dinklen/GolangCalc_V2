// Temporarily stopped
package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commandsInfo = make(map[string]string)

func NewHelpCommand() *cobra.Command {
	var help string

	command := &cobra.Command{
		Use: "help",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(commandsInfo[help])
		},
	}

	command.Flags().StringVarP(&help, "command", "c", "go_calc", "Help command")

	return command
}

func init() {
	viper.SetConfigName("help")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../internal/cli/commands")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("\033[1;38;5;88m[FATAL]\033[0m failed to read help.yaml")
		os.Exit(1)
	}

	for key := range viper.AllSettings() {
		commandsInfo[key] = viper.GetString(key)
	}
}
