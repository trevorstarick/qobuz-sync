package cmds

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct,gochecknoglobals
var Favorites = &cobra.Command{
	Use:   "favorites <albums|tracks|albums+tracks>",
	Short: "Download all favorite albums and/or tracks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get client from context")
		}

		switch args[0] {
		case "albums":
			err = client.FavoriteAlbums()
			if err != nil {
				return errors.Wrap(err, "unable to download favorite albums")
			}
		case "tracks":
			err = client.FavoriteTracks()
			if err != nil {
				return errors.Wrap(err, "unable to download favorite tracks")
			}
		case "albums+tracks", "tracks+albums":
			err = client.FavoriteTracks()
			if err != nil {
				return errors.Wrap(err, "unable to download favorite tracks")
			}

			err = client.FavoriteAlbums()
			if err != nil {
				return errors.Wrap(err, "unable to download favorite albums and tracks")
			}
		}

		return nil
	},
}
