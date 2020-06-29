package deezer

import "errors"

type Song struct {
	Data SongData `json:"DATA"`
}

type SongData struct {
	ID             string `json:"SNG_ID"`
	Title          string `json:"SNG_TITLE"`
	ArtistName     string `json:"ART_NAME"`
	AlbumTitle     string `json:"ALB_TITLE"`
	AlbumID        string `json:"ALB_ID"`
	MD5Origin      string `json:"MD5_ORIGIN"`
	MediaVersion   string `json:"MEDIA_VERSION"`
	FilesizeMP3128 int    `json:"FILESIZE_MP3_128,string"`
	FilesizeMP3320 int    `json:"FILESIZE_MP3_320,string"`
	FilesizeFLAC   int    `json:"FILESIZE_FLAC,string"`
}

type Quality int

const (
	MP3128 Quality = 1
	MP3320         = 3
	FLAC           = 9
)

type Album struct {
	Data  AlbumData `json:"DATA"`
	Songs struct {
		Data []SongData `json:"data"`
	} `json:"SONGS"`
}

type AlbumData struct {
	ID    string `json:"ALB_ID"`
	Title string `json:"ALB_TITLE"`
}

func GetValidSongQuality(s SongData, preferred Quality) (Quality, error) {
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