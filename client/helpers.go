package client

import (
	"os"
	"strings"

	"github.com/frolovo22/tag"
	"github.com/pkg/errors"
	"github.com/trevorstarick/qobuz-sync/common"
)

func SetTags(path string, tags common.Metadata) error {
	fileTags, err := tag.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "failed to read part file for tags")
	}

	for _, err := range []error{
		fileTags.SetAlbum(tags.Album),
		fileTags.SetAlbumArtist(tags.AlbumArtist),
		fileTags.SetArtist(tags.Artist),

		fileTags.SetComment(tags.Comment),
		fileTags.SetGenre(tags.Genre),
		fileTags.SetDate(tags.Date),
		fileTags.SetComposer(tags.Composer),

		fileTags.SetTitle(tags.Title),
		fileTags.SetTrackNumber(tags.Track, tags.TrackTotal),
		fileTags.SetDiscNumber(tags.Disc, tags.DiscTotal),
	} {
		if err != nil {
			return errors.Wrap(err, "failed to set tag")
		}
	}

	err = fileTags.SaveFile(strings.TrimSuffix(path, ".part"))
	if err != nil {
		return errors.Wrap(err, "failed to save tags")
	}

	err = os.Remove(path)
	if err != nil {
		return errors.Wrap(err, "failed to remove part file")
	}

	return nil
}
