//nolint:gochecknoglobals
package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var preRun = func(cmd *cobra.Command, args []string) error {
	username := os.Getenv("QOBUZ_USERNAME")
	password := os.Getenv("QOBUZ_PASSWORD")

	if username == "" || password == "" {
		log.Error().Msg("QOBUZ_USERNAME and QOBUZ_PASSWORD envvars must be set")

		return errors.Wrap(ErrAuthFailed, "missing credentials")
	}

	baseDir := DefaultBaseDir
	if os.Getenv("QOBUZ_BASEDIR") != "" {
		baseDir = os.Getenv("QOBUZ_BASEDIR")
	}

	err := os.MkdirAll(baseDir, dirPerm)
	if err != nil {
		return errors.Wrap(err, "unable to create base dir")
	}

	trackTracker, err = NewTracker(filepath.Join(baseDir, "tracks.txt"))
	if err != nil {
		return errors.Wrap(err, "unable to create track tracker")
	}

	albumTracker, err = NewTracker(filepath.Join(baseDir, "albums.txt"))
	if err != nil {
		return errors.Wrap(err, "unable to create album tracker")
	}

	c, err := NewClient(username, password)
	if err != nil {
		return errors.Wrap(err, "unable to create client")
	}

	handler := NewHandler(c, baseDir)
	ctx := context.WithValue(cmd.Context(), Handler{}, handler) //nolint:exhaustruct
	cmd.SetContext(ctx)

	return nil
}

//nolint:exhaustruct
var albumCommand = &cobra.Command{
	Use:   "album <id> [id...]",
	Short: "Download an album",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := getHandlerFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		for _, id := range args {
			err = handler.Album(id)
			if err != nil {
				return errors.Wrap(err, "unable to download album")
			}
		}

		return nil
	},
}

var postRun = func(cmd *cobra.Command, args []string) error {
	err := trackTracker.Close()
	if err != nil {
		return errors.Wrap(err, "unable to close track tracker")
	}

	err = albumTracker.Close()
	if err != nil {
		return errors.Wrap(err, "unable to close album tracker")
	}

	return nil
}

//nolint:exhaustruct
var trackCommand = &cobra.Command{
	Use:   "track <id> [id...]",
	Short: "Download a track",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := getHandlerFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		for _, id := range args {
			err = handler.Track(id)
			if err != nil {
				return errors.Wrap(err, "unable to download track")
			}
		}

		return nil
	},
}

//nolint:exhaustruct
var searchCommand = (&cobra.Command{
	Use:   "search <query>",
	Short: "Search for albums and tracks",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := getHandlerFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		err = handler.Search(strings.Join(args, " "))
		if err != nil {
			return errors.Wrap(err, "unable to search")
		}

		return nil
	},
})

//nolint:exhaustruct
var playlistCommand = (&cobra.Command{
	Use:   "playlist <id> [id...]",
	Short: "Download a playlist",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := getHandlerFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		for _, id := range args {
			err = handler.Playlist(id)
			if err != nil {
				return errors.Wrap(err, "unable to download playlist")
			}
		}

		return nil
	},
})

//nolint:exhaustruct
var favoritesCommand = &cobra.Command{
	Use:   "favorites <albums|tracks|albums+tracks>",
	Short: "Download all favorite albums and/or tracks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := getHandlerFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		switch args[0] {
		case "albums":
			err = handler.FavoriteAlbums()
			if err != nil {
				return err
			}
		case "tracks":
			err = handler.FavoriteTracks()
			if err != nil {
				return err
			}
		case "albums+tracks", "tracks+albums":
			err = handler.FavoriteTracks()
			if err != nil {
				return err
			}

			err = handler.FavoriteAlbums()
			if err != nil {
				return err
			}
		}

		return nil
	},
}
