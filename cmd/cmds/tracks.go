package cmds

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct,gochecknoglobals
var Track = &cobra.Command{
	Use:   "track <id> [id...]",
	Short: "Download a track",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := GetClientFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		for _, id := range args {
			err = handler.DownloadTrack(id)
			if err != nil {
				return errors.Wrap(err, "unable to download track")
			}
		}

		return nil
	},
}
