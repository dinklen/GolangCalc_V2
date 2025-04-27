package commands

import (
	"fmt"
	// errors
	// config

	"github.com/spf13/cobra"
)

func setData(logic func(string, string)) *cobra.Command {
	var login, password = "", ""

	command := &cobra.Command{
		Use:   "sign_up",
		Short: "Register a new account",
		Run: func(cmd *cobra.Command, args []string) {
			if login == "" || password == "" {
				fmt.Println("login or password is empty")
			}
			logic(login, password)
		},
	}

	command.Flags().StringVarP(&login, "login", "l", "", "User's login")
	command.Flags().StringVarP(&password, "password", "p", "", "User's password")

	command.MarkFlagRequired("login")
	command.MarkFlagRequired("password")

	return command
}

func NewSignUpCommand() *cobra.Command {
	command := setData(func(login string, password string) {
		fmt.Printf("login: %v; password: %v\n", login, password)
	})

	return command
}

func NewStartSessionCommand() *cobra.Command {
	command := setData(func(string, string) {
		fmt.Println("session was started")
		// cfg.update->status
	})

	return command
}

func NewLogoutCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "logout",
		Short: "Logout",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("logouted")
			// check online status
		},
	}

	return command
}
