package cmds

import (
	"fmt"
	"strings"

	//"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//nolint:exhaustruct,gochecknoglobals
var Search = &cobra.Command{
	Use:    "search <query>",
	Short:  "Search for albums and tracks",
	Args:   cobra.MinimumNArgs(1),
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

		//spew.Dump(res)

		if len(res.Albums.Items) > 0 {
			fmt.Println("==== Albums ====================")
			for _, v := range res.Albums.Items {
				fmt.Println(v.Title, "/", v.Artist.Name)
				fmt.Println(v.URL)
				fmt.Println()
			}
		}

		if len(res.Artists.Items) > 0 {
			fmt.Println("==== Artists ====================")
			for _, v := range res.Artists.Items {
				fmt.Println(v.Name)
			}
		}

		if len(res.Playlists.Items) > 0 {
			fmt.Println("==== Playlists ====================")
			for _, v := range res.Playlists.Items {
				fmt.Println(v.Name)
			}
		}

		if len(res.Tracks.Items) > 0 {
			fmt.Println("==== Tracks ====================")
			for _, v := range res.Tracks.Items {
				fmt.Println(v.Title, "/", v.Album.Title, "/", v.Performer.Name)
			}
		}

		return nil
	},
}
