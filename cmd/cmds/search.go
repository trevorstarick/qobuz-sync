package cmds

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func shorten(str string) string {
	//nolint:gomnd
	if len(str) > 60 {
		return str[:60] + "..."
	}

	return str
}

//nolint:exhaustruct,gochecknoglobals,forbidigo
var Search = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for albums and tracks",
	Args:  cobra.MinimumNArgs(1),
	// Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientFromContext(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "unable to get client from context")
		}

		res, err := client.Search(strings.Join(args, " "))
		if err != nil {
			return errors.Wrap(err, "unable to search")
		}

		if len(res.Albums.Items) > 0 {
			fmt.Println()
			fmt.Println("==== Albums ====================")
			for _, item := range res.Albums.Items {
				bits := item.MaximumBitDepth
				sampleRate := item.MaximumSamplingRate
				channels := item.MaximumChannelCount
				details := ""

				if channels == 1 {
					details = fmt.Sprintf("(mono/%v bits/%v kHz)", bits, sampleRate)
				} else if bits != 16 || sampleRate != 44.1 {
					details = fmt.Sprintf("(%v bits/%v kHz)", bits, sampleRate)
				}

				fmt.Println(item.ID, item.Title, "/", item.Artist.Name, details)
			}
		}

		if len(res.Artists.Items) > 0 {
			fmt.Println()
			fmt.Println("==== Artists ====================")
			for _, v := range res.Artists.Items {
				fmt.Println(v.ID, v.Name)
			}
		}

		if len(res.Playlists.Items) > 0 {
			fmt.Println()
			fmt.Println("==== Playlists ====================")
			for _, v := range res.Playlists.Items {
				fmt.Println(v.ID, v.Name, "/", v.TracksCount, "/", shorten(v.Description))
			}
		}

		if len(res.Tracks.Items) > 0 {
			fmt.Println()
			fmt.Println("==== Tracks ====================")
			for _, v := range res.Tracks.Items {
				performer := v.Performers
				if v.Performer != nil {
					performer = v.Performer.Name
				}
				fmt.Println(v.ID, v.Title, "/", v.Album.Title, "/", performer)
			}
		}

		return nil
	},
}
