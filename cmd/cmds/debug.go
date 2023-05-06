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
	Use:    "debug",
	Short:  "Debug commands",
	Hidden: true,
	Args:   cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler, err := GetClientFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get handler from context")
		}

		var res any

		switch args[0] {
		case "album":
			res, err = handler.AlbumGet(args[1])
			if err != nil {
				return errors.Wrap(err, "unable to get album")
			}
		case "track":
			res, err = handler.TrackGet(args[1])
			if err != nil {
				return errors.Wrap(err, "unable to get track")
			}

			spew.Dump(res)
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