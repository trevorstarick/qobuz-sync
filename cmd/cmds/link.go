package cmds

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct,gochecknoglobals
var Link = &cobra.Command{
	Use:   "link <url> [url...]",
	Short: "Download an album or track from a URL",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, url := range args {
			client, err := GetClientFromContext(cmd.Context())
			if err != nil {
				return errors.Wrap(err, "unable to get client from context")
			}

			err = client.Link(url)
			if err != nil {
				return errors.Wrap(err, "unable to download link")
			}
		}

		return nil
	},
}
