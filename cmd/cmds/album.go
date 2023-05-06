package cmds

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct,gochecknoglobals
var Album = &cobra.Command{
	Use:   "album <id> [id...]",
	Short: "Download an album",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := GetClientFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		for _, id := range args {
			err = handler.DownloadAlbum(id)
			if err != nil {
				return errors.Wrap(err, "unable to download album")
			}
		}

		return nil
	},
}
