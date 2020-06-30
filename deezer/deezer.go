package deezer

import (
	"errors"
	"fmt"
	"regexp"
)

type Quality int

const (
	MP3128 Quality = 1
	MP3320         = 3
	FLAC           = 9
)

type ContentType string

const (
	ContentAlbum  ContentType = "album"
	ContentArtist             = "artist"
	ContentSong               = "track"
)

type Song struct {
	Data SongData `json:"DATA"`
}

type SongData struct {
	ID                string       `json:"SNG_ID"`
	ProductTrackID    string       `json:"PRODUCT_TRACK_ID"`
	UploadID          int          `json:"UPLOAD_ID"`
	Title             string       `json:"SNG_TITLE"`
	ArtistID          string       `json:"ART_ID"`
	ProviderID        string       `json:"PROVIDER_ID"`
	ArtistName        string       `json:"ART_NAME"`
	Artists           []ArtistData `json:"ARTISTS"`
	AlbumID           string       `json:"ALB_ID"`
	AlbumTitle        string       `json:"ALB_TITLE"`
	MD5Origin         string       `json:"MD5_ORIGIN"`
	Video             bool         `json:"VIDEO"`
	Duration          string       `json:"DURATION"`
	AlbumPicture      string       `json:"ALB_PICTURE"`
	ArtistPicture     string       `json:"ART_PICTURE"`
	Rank              string       `json:"RANK_SNG"`
	FilesizeMP3128    int          `json:"FILESIZE_MP3_128,string"`
	FilesizeMP3320    int          `json:"FILESIZE_MP3_320,string"`
	FilesizeFLAC      int          `json:"FILESIZE_FLAC,string"`
	Filesize          string       `json:"FILESIZE"`
	MediaVersion      string       `json:"MEDIA_VERSION"`
	DiskNumber        string       `json:"DISK_NUMBER"`
	TrackNumber       string       `json:"TRACK_NUMBER"`
	Version           string       `json:"VERSION"`
	ExplicitLyrics    string       `json:"EXPLICIT_LYRICS"`
	ISRC              string       `json:"ISRC"`
	HierarchicalTitle string       `json:"HIERARCHICAL_TITLE"`
	LyricsID          int          `json:"LYRICS_ID"`
	Status            int          `json:"STATUS"`
}

type Album struct {
	Data  AlbumData `json:"DATA"`
	Songs struct {
		Data []SongData `json:"data"`
	} `json:"SONGS"`
}

type AlbumData struct {
	ID                  string `json:"ALB_ID"`
	ArtistID            string `json:"ART_ID"`
	ArtistName          string `json:"ART_NAME"`
	LabelName           string `json:"LABEL_NAME"`
	StyleName           string `json:"STYLE_NAME"`
	Title               string `json:"ALB_TITLE"`
	Version             string `json:"VERSION"`
	Picture             string `json:"ALB_PICTURE"`
	DigitalReleaseDate  string `json:"DIGITAL_RELEASE_DATE"`
	PhysicalReleaseDate string `json:"PHYSICAL_RELEASE_DATE"`
	ProviderID          string `json:"PROVIDER_ID"`
	SonyProdID          string `json:"SONY_PROD_ID"`
	UPC                 string `json:"UPC"`
	Status              string `json:"STATUS"`
	Fans                int    `json:"NB_FAN"`
	Available           bool   `json:"AVAILABLE"`
	Songs               *struct {
		Data []SongData `json:"data"`
	} `json:"SONGS"` // empty when fetching an album with client.Album, but present in artist.Albums[n].AlbumData from client.Artist
}

type Artist struct {
	Data   ArtistData `json:"DATA"`
	Albums struct {
		Data []AlbumData `json:"data"`
	} `json:"ALBUMS"`
}

type ArtistData struct {
	ID      string `json:"ART_ID"`
	Name    string `json:"ART_NAME"`
	Picture string `json:"ART_PICTURE"`
	Fans    int    `json:"NB_FAN"`
}

// ValidSongQuality returns the preferred audio quality if it is availabe
// for the song, the highest available quality otherwise, and an error
// if there are no available qualities for the song.
func ValidSongQuality(s SongData, preferred Quality) (Quality, error) {
	var qualities []Quality
	switch {
	case s.FilesizeFLAC != 0:
		qualities = append(qualities, FLAC)
	case s.FilesizeMP3320 != 0:
		qualities = append(qualities, MP3320)
	case s.FilesizeMP3128 != 0:
		qualities = append(qualities, MP3128)
	}
	for _, q := range qualities {
		if q == preferred {
			return q, nil
		}
	}
	if len(qualities) > 0 {
		return qualities[0], nil
	}
	return Quality(0), errors.New("no valid song quality")
}

// LinkContentType returns the content type and id of a given Deezer link.
func ParseLink(link string) (ctype ContentType, id string) {
	re := regexp.MustCompile(`deezer.com(?:\/[a-zA-Z]{2})?\/(album|artist|track)\/(\d+)`)
	m := re.FindAllStringSubmatch(link, -1)
	if len(m) < 1 || len(m[0]) < 3 {
		return "", ""
	}
	switch m[0][1] {
	case "album":
		ctype = ContentAlbum
	case "artist":
		ctype = ContentArtist
	case "track":
		ctype = ContentSong
	}
	return ctype, m[0][2]
}

// Link returns a link from a given content type and content id,
// being essentially the opposite of ParseLink.
func Link(ctype ContentType, id string) (link string) {
	return fmt.Sprintf("https://www.deezer.com/%s/%s", ctype, id)
}
