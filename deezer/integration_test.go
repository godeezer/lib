// +build integration

package deezer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"
)

func mustClient(t *testing.T) *Client {
	arl := os.Getenv("ARL")
	if arl == "" {
		t.Fatal("Missing $ARL")
	}
	client, err := NewClient(arl)
	if err != nil {
		t.Fatal("Error creating client:", err)
	}
	return client
}

func TestDownload(t *testing.T) {
	client := mustClient(t)
	song, err := client.Song("793554582")
	if err != nil {
		t.Fatal("Error getting song:", err)
	}
	reader, err := client.Download(*song, MP3128)
	if err != nil {
		t.Fatal("Error creating decrypting reader for song:", err)
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		t.Fatal("Error downloading/decrypting song:", err)
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	if sum != "6044f325dd38bd9ae74e29171918798a09f8d8661907827aa609bcd01e9ca65d" {
		t.Fatal("Got incorrect hash for downloaded song:", sum)
	}
}
