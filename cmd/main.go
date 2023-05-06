package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/trevorstarick/qobuz-sync/client"
	"github.com/trevorstarick/qobuz-sync/cmd/cmds"
	"github.com/trevorstarick/qobuz-sync/common"
)

//nolint:gochecknoglobals
var (
	DefaultBaseDir = "downloads"
	Version        = "dev"
	Revision       = ""
)

//nolint:gochecknoglobals
var preRun = func(cmd *cobra.Command, args []string) error {
	username := os.Getenv("QOBUZ_USERNAME")
	password := os.Getenv("QOBUZ_PASSWORD")

	if username == "" || password == "" {
		log.Error().Msg("QOBUZ_USERNAME and QOBUZ_PASSWORD envvars must be set")

		return errors.Wrap(common.ErrAuthFailed, "missing credentials")
	}

	baseDir := DefaultBaseDir
	if os.Getenv("QOBUZ_BASEDIR") != "" {
		baseDir = os.Getenv("QOBUZ_BASEDIR")
	}

	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "unable to create base dir")
	}

	c, err := client.NewClient(username, password, baseDir)
	if err != nil {
		return errors.Wrap(err, "unable to create client")
	}

	ctx := context.WithValue(cmd.Context(), client.Key{}, c)
	cmd.SetContext(ctx)

	return nil
}

//nolint:gochecknoglobals
var postRun = func(cmd *cobra.Command, args []string) error {
	client, err := cmds.GetClientFromContext(cmd.Context())
	if err != nil {
		return errors.Wrap(err, "unable to get handler from context")
	}

	err = client.Close()
	if err != nil {
		return errors.Wrap(err, "unable to close client")
	}

	return nil
}

func envOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

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

	cmds.Debug.PersistentFlags().String("output", "spew", "output format (json, spew)")
	cmd.AddCommand(cmds.Debug)

	cmd.AddCommand(
		cmds.Album,
		cmds.Track,
		cmds.Search,
		cmds.Playlist,
		cmds.Favorites,
		cmds.Link,
	)

	err = cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
