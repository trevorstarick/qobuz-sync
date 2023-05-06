- `playlist/get`
    - the api returns a rectangular album art but also the first 4 album covers from the first 4 songs, we can combine and make a collage then save to disk and use that

- `downloadAlbumArt()`
    - the `res.Body`'s [magic numbers](https://en.wikipedia.org/wiki/List_of_file_signatures) should be checked to see if it's a valid image file and not a rogue binary

- `downloadFile()`
    - the `res.Body`'s [magic numbers](https://en.wikipedia.org/wiki/List_of_file_signatures) should be checked to see if it's a valid audio file and not a rogue binary
