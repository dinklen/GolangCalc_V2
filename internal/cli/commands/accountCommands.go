// Temporarily stopped
package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dinklen/GolangCalc_V2/internal/models"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Response struct {
	Message interface{} `json:"message"`
	Error   string      `json:"error"`
}

func setData(logic func(string, string)) *cobra.Command {
	var login, password = "", ""

	command := &cobra.Command{
		Use:   "sign_up",
		Short: "Register a new account",
		Run: func(cmd *cobra.Command, args []string) {
			if login == "" || password == "" {
				log.Println("\033[1;31m[ERROR]\033[0m login or password is empty")
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

func NewAccountPOSTRequest(client *http.Client, url string) *cobra.Command {
	command := setData(func(login string, password string) {
		// send request to server
		// get answer
		// fmt.Println(answer)

		data := &models.AccountData{
			Login:        login,
			PasswordHash: password,
		} //unique func

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("\033[1;33[ERROR]\033[0m " + err.Error())
			return
		}

		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("\033[1;33[ERROR]\033[0m " + err.Error())
			return
		}

		request.Header.Set("Content-Type", "application/json")
		// if

		resp, err := client.Do(request)
		if err != nil {
			log.Println("\033[1;33[ERROR]\033[0m " + err.Error())
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("\033[1;33[ERROR]\033[0m " + err.Error())
			return
		}

		response := new(Response) // unique
		err = json.Unmarshal(body, response)
		if err != nil {
			log.Println("\033[1;33[ERROR]\033[0m " + err.Error())
		}

		if response.Error != "" { //func!
			fmt.Println(response.Error)
			return
		}
		fmt.Println(response.Message) //!func
	})

	return command
}

func NewPOSTRequest(
	start *interface{},
	auth bool,
	outData interface{}, output func(),
	httpClient *http.Client,
	serverURL, redisURL string,
	logger *zap.Logger,
	rc *redis.Client,
) *cobra.Command {
	command := setData(func(login string, password string) {
		data := start

		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		request, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Error(err.Error())
			return
		}

		request.Header.Set("Content-Type", "application/json")
		if auth {
			token, err := rc.Get(context.Background(), "access")
			request.Header.Set("Authorization", "redis")
		}
	})

	return command
}

func NewStartSessionCommand(client *http.Client, url string) *cobra.Command {
	command := setData(func(string, string) {
		// ...
	})

	return command
}

func NewLogoutCommand(client *http.Client, url string) *cobra.Command {
	command := &cobra.Command{
		Use:   "logout",
		Short: "Logout",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("logouted")
		},
	}

	return command
}
