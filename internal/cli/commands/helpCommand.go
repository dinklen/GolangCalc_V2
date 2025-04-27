package commands

import (
	"fmt"

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
		panic(fmt.Errorf("Failed to read help.yaml: %v", err))
	}

	for key := range viper.AllSettings() {
		commandsInfo[key] = viper.GetString(key)
	}

	//fmt.Printf("config: %+v\n", commandsInfo)
}
