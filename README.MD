qobuz-sync
==========

A simple script to sync your local music library with your Qobuz account.

```text
Download albums and tracks from Qobuz

Usage:
  qobuz-sync [command]

Available Commands:
  album       Download an album
  completion  Generate the autocompletion script for the specified shell
  favorites   Download all favorite albums and/or tracks
  help        Help about any command
  link        Download an album or track from a URL
  playlist    Download a playlist
  track       Download a track

Flags:
  -h, --help      help for qobuz-sync
  -v, --version   version for qobuz-sync

Use "qobuz-sync [command] --help" for more information about a command.
```

## Requirements

Set the following environment variables:

```bash
export QOBUZ_USERNAME="your-qobuz-username"
export QOBUZ_PASSWORD="your-qobuz-password"
```

additionally you can set the following environment variables to set the download directory (by default it's `./downloads`):

```bash
export QOBUZ_BASEDIR="path-to-your-download-directory"
```

## Debugging

To enable debug logging, set the `DEBUG` environment variable to `true`:

```bash
export DEBUG=true
```

To get debug output from the requests, you can use the hidden debug command: 
```bash
$ qobuz-sync debug --output json <album|track> <id>
```

This will return the response objects from the Qobuz API as JSON (or [spew](https://github.com/davecgh/go-spew) format).

