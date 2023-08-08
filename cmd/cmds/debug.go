package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	qlient "github.com/trevorstarick/qobuz-sync/client"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals,exhaustruct
var Debug = &cobra.Command{
	Use:    "debug <album|tracki|favorites|search> <id|type|search-query>",
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
			if len(args) > 2 {
				var format qlient.TrackFormat

				switch strings.ToUpper(args[2]) {
				case "MP3":
					format = qlient.QualityMP3
				case "FLAC":
					format = qlient.QualityFLAC
				case "HIRES":
					format = qlient.QualityHIRES
				case "MAX":
					format = qlient.QualityMAX
				}

				res, err = client.TrackGetFileURL(args[1], format)
				if err != nil {
					return errors.Wrap(err, "unable to get track url")
				}
			} else {
				res, err = client.TrackGet(args[1])
				if err != nil {
					return errors.Wrap(err, "unable to get track")
				}
			}
		case "favorites":
			if args[1] == "albums-tracks" {
				fres, err := client.FavoriteGetUserFavorites(qlient.ListTypeALBUM, 0)
				if err != nil {
					return errors.Wrap(err, "unable to get favorite")
				}

				var alist []any

				for _, album := range fres.Albums.Items {
					ares, err := client.AlbumGet(album.ID)
					if err != nil {
						return errors.Wrap(err, "unable to get album")
					}

					alist = append(alist, ares)
				}

				res = alist
				break
			}

			res, err = client.FavoriteGetUserFavorites(qlient.ListType(args[1]), 0)
			if err != nil {
				return errors.Wrap(err, "unable to get favorite")
			}
		case "search":
			res, err = client.Search(strings.Join(args[1:], " "))
			if err != nil {
				return errors.Wrap(err, "unable to search")
			}
		default:
			return errors.Errorf("unknown command %q", args[0])
		}

		switch cmd.Flag("output").Value.String() {
		case "spew":
			spew.Dump(res)
		case "json":
			bytes, err := json.Marshal(res)
			if err != nil {
				return errors.Wrap(err, "unable to marshal response")
			}

			fmt.Fprintf(os.Stdout, "%s\n", bytes)
		default: // case "json-pretty":
			bytes, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return errors.Wrap(err, "unable to marshal response")
			}

			fmt.Fprintf(os.Stdout, "%s\n", bytes)
		}

		return nil
	},
}
