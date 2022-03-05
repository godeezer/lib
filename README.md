# godeezer/lib

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/godeezer/lib)
[![Report Card](https://goreportcard.com/badge/github.com/godeezer/lib)](https://goreportcard.com/report/github.com/godeezer/lib)

A Go library for interacting with Deezer

## Example usage
These are just a few examples, view the godoc for full documentation.
(error handling omitted for brevity)
```go
// create client with arl token stored in $ARL
client, _ := deezer.NewClient(os.Getenv("ARL"))

// fetch a song
song, _ := client.Song("1297748632")
// download that song
r, _ := client.Download(song, deezer.MP3128)
// r is an io.ReadCloser, you can copy it into an io.Writer etc

// fetch an album
album, _ := client.Album("219026842")
// fetch an album's songs
songs, _ := client.SongsByAlbum("219026842", -1)
```

## Why can't I download high-quality music?
Deezer implemented a server-side restriction which only allows premium/hi-fi
Deezer accounts to download songs at a quality that isn't 128 kbps MP3. If you
need to download higher-quality music then you must use the ARL of a premium
account.

## Contributing
Pull requests and issues are welcome.

## License
This library is free software. See the ISC license included in LICENSE.
