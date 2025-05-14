// Temporarily stopped
package cli

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/cli/commands"
	"github.com/dinklen/GolangCalc_V2/internal/config"
	redisc "github.com/dinklen/GolangCalc_V2/internal/redis"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var rootCmd *cobra.Command

func createHTTPClient() *http.Client {
	caCert, _ := os.ReadFile("../../certs/cert.pem")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}

	client := &http.Client{Transport: transport}

	return client
}

func initLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("\x1b[38;5;87m" + t.Format("02.01.2006 15:04:05.000") + "\x1b[0m")
	}

	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		color := ""
		text := l.CapitalString()

		switch l {
		case zapcore.DebugLevel:
			color = "\x1b[1;35m"
			text = "BYE"
		case zapcore.InfoLevel:
			color = "\x1b[1;34m"
		case zapcore.WarnLevel:
			color = "\x1b[1;33m"
		case zapcore.ErrorLevel:
			color = "\x1b[1;31m"
		case zapcore.FatalLevel:
			color = "\x1b[1;38;5;88m"
		}

		enc.AppendString(color + "[" + text + "]\x1b[0m")
	}

	config.DisableCaller = true
	config.DisableStacktrace = true

	logger, _ := config.Build()

	return logger
}

func initConfig(logger *zap.Logger) *config.Config {
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return cfg
}

func initRedisClient(cfg *config.Config, logger *zap.Logger) *redis.Client {
	redisClient, err := redisc.NewClient(cfg, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return redisClient
}

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

	httpClient := createHTTPClient()

	logger := initLogger()
	cfg := initConfig(logger)
	rc := initRedisClient(cfg, logger)

	root.AddCommand(
		commands.NewSignUpCommand(httpClient, url),
		commands.NewStartSessionCommand(httpClient, url),
		commands.NewLogoutCommand(httpClient, url),
		commands.NewHelpCommand(),
	)

	return root
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd = setupCommands()
}
