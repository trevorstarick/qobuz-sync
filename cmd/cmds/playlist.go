package cmds

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct,gochecknoglobals
var Playlist = &cobra.Command{
	Use:   "playlist <id> [id...]",
	Short: "Download a playlist",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get client from context")
		}

		for _, id := range args {
			err = client.DownloadPlaylist(id)
			if err != nil {
				return errors.Wrap(err, "unable to download playlist")
			}
		}

		return nil
	},
}
