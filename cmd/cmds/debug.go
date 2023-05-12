package cmds

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals,exhaustruct
var Debug = &cobra.Command{
	Use:    "debug <album|track> <id>",
	Short:  "Debug commands",
	Hidden: true,
	Args:   cobra.MinimumNArgs(2), //nolint:gomnd
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get client from context")
		}

		var res any

		switch args[0] {
		case "album":
			res, err = client.AlbumGet(args[1])
			if err != nil {
				return errors.Wrap(err, "unable to get album")
			}
		case "track":
			res, err = client.TrackGet(args[1])
			if err != nil {
				return errors.Wrap(err, "unable to get track")
			}
		default:
			return errors.Errorf("unknown command %q", args[0])
		}

		switch cmd.Flag("output").Value.String() {
		case "json":
			bytes, err := json.Marshal(res)
			if err != nil {
				return errors.Wrap(err, "unable to marshal response")
			}

			fmt.Fprintf(os.Stdout, "%s\n", bytes)
		default:
			spew.Dump(res)
		}

		return nil
	},
}
