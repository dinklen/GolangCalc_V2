package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_calc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
		   examples and usage of using your application. For example:

		   Cobra is a CLI library for Go that empowers applications.
		   This application is a tool to generate the needed files
		   to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, world!")
	},
}

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "This is a subcommand",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Subcommand")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(subCmd)
	subCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
