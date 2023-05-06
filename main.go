package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	DefaultBaseDir = "downloads"
	Version        = "dev"
	Revision       = ""
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "timestamp"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //nolint:exhaustruct

	logLevel, err := zerolog.ParseLevel(envOrDefault("LOG_LEVEL", "info"))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse loglevel")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}

	if envOrDefault("DEBUG", "false") == "true" && logLevel > zerolog.DebugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				Revision = setting.Value
			}
		}
	}

	//nolint:exhaustruct
	cmd := &cobra.Command{
		Use:                "qobuz-sync",
		Short:              "Download albums and tracks from Qobuz",
		Version:            fmt.Sprintf("%v %v", Version, Revision),
		PersistentPreRunE:  preRun,
		PersistentPostRunE: postRun,
	}

	cmd.AddCommand(
		albumCommand,
		trackCommand,
		searchCommand,
		playlistCommand,
		favoritesCommand,
	)

	err = cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
