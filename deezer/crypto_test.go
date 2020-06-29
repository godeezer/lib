package deezer

import (
	"bytes"
	"testing"
)

// Peroxide by Ecco2k// (Drain Gang)
var song = SongData{
	ID:             "793554582",
	MD5Origin:      "19b61b6fe1faf7a77914a0a5180593af",
	MediaVersion:   "3",
	FilesizeMP3128: 3428100,
}

var expectedfilename = "57cd93a24a5b185e71f5a3e32b991dc672c30161b06b38bd95eea12161c5574f17f07149a52aefa89ccfede56ebadb098d4b589796c63518f728b5895536b480a533d219f8dafaaffba6e0697c0b57e5"

func TestSongFilename(t *testing.T) {
	filename, err := songFilename(song, Quality(0))
	if err != nil {
		t.Errorf(err.Error())
	}
	if filename != expectedfilename {
		t.Errorf("Expected: %s Got: %s", expectedfilename, filename)
	}
}

func TestSongDownloadURL(t *testing.T) {
	url, err := songDownloadURL(song, Quality(0))
	if err != nil {
		t.Errorf(err.Error())
	}
	expected := "https://e-cdns-proxy-1.dzcdn.net/mobile/1/" + expectedfilename
	if url != expected {
		t.Errorf("Expected: %s Got: %s", expected, url)
	}
}

func TestGetBlowfishKey(t *testing.T) {
	key := getBlowfishKey(song.ID)
	expected := []byte("`ic;7h$a3\x7F%29>dc")
	if !bytes.Equal(key, expected) {
		t.Errorf("Expected: %s Got: %s", string(expected), string(key))
	}
}
