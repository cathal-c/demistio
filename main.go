package main

import (
	"context"
	"github.com/cathal-c/demistio/internal/app"
	"github.com/rs/zerolog"
	"os"
)

const (
	version = "v1.25.1"
)

func main() {
	appConfig := app.ParseFlagsToConfig()

	log := setupLogs(appConfig.LogLevel)
	ctx := log.WithContext(context.Background())

	if err := app.Generate(ctx, appConfig); err != nil {
		log.Fatal().Err(err).Msg("failed to generate configuration")
	}
}

func setupLogs(level zerolog.Level) zerolog.Logger {
	log := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(level)

	return log
}
