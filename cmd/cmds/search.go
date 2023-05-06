package cmds

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct,gochecknoglobals
var Search = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for albums and tracks",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := GetClientFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		err = handler.Search(strings.Join(args, " "))
		if err != nil {
			return errors.Wrap(err, "unable to search")
		}

		return nil
	},
}
