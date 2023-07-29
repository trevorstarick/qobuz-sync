package cmds

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func shorten(s string) string {
	if len(s) > 60 {
		return s[:60] + "..."
	}

	return s
}

//nolint:exhaustruct,gochecknoglobals
var Search = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for albums and tracks",
	Args:  cobra.MinimumNArgs(1),
	//Hidden: true,
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
			for _, v := range res.Albums.Items {
				fmt.Println(v.ID, v.Title, "/", v.Artist.Name)
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
				fmt.Println(v.ID, v.Title, "/", v.Album.Title, "/", v.Performer.Name)
			}
		}

		return nil
	},
}
