package deezer

import "testing"

func TestParseLink(t *testing.T) {
	ctype, id := ParseLink("https://www.deezer.com/album/13357219")
	if ctype != "album" || id != "13357219" {
		t.Errorf("Expected album and 13357219 Got %s and %s", ctype, id)
	}
}

func TestLink(t *testing.T) {
	link := Link("album", "13357219")
	if link != "https://www.deezer.com/album/13357219" {
		t.Errorf("Expected https://www.deezer.com/album/13357219 Got %s", link)
	}
}