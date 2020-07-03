package deezer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const apiURL = "https://www.deezer.com/ajax/gw-light.php"

type apiMethod string

const (
	getUserData         apiMethod = "deezer.getUserData"
	songGetData                   = "song.getData"
	songListByAlbum               = "song.getListByAlbum"
	albumGetData                  = "album.getData"
	artistGetData                 = "artist.getData"
	albumGetDiscography           = "album.getDiscography"
)

type songGetDataBody struct {
	ID string `json:"sng_id"`
}

type albumGetDataBody struct {
	ID string `json:"alb_id"`
}

type artistGetDataBody struct {
	ID string `json:"art_id"`
}

type songListByAlbumBody struct {
	ID    string `json:"alb_id"`
	Limit int    `json:"nb"`
}

type albumGetDiscographyBody struct {
	ArtistID   string `json:"art_id"`
	Language   string `json:"lang"`
	FilterRole []int  `json:"filter_role_id"`
	Limit      int    `json:"nb"`
	LimitSongs int    `json:"nb_songs"`
	Start      int    `json:"start"`
}

type userData struct {
	CheckForm string `json:"checkForm"`
}

type response struct {
	Results json.RawMessage `json:"results"`
}

type multiSongResponse struct {
	Data []Song `json:"data"`
}

type multiAlbumResponse struct {
	Data []Album `json:"data"`
}

type Client struct {
	*http.Client
	Arl string
}

// NewClient returns a Deezer client with
// the given arl (used for authentication)
// this arl can be gotten by following these instructions:
// https://notabug.org/RemixDevs/DeezloaderRemix/wiki/Login+via+userToken
func NewClient(arl string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(url,
		[]*http.Cookie{
			&http.Cookie{
				Name:  "arl",
				Value: arl,
			},
		},
	)
	client := &Client{
		&http.Client{
			Jar: jar,
		}, arl,
	}
	return client, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,en-US;q=0.8,en;q=0.7")
	return c.Client.Do(req)
}

func (c *Client) apiDo(method apiMethod, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return nil, err
	}
	var token string
	if method != getUserData {
		t, err := c.csrfToken()
		if err != nil {
			return nil, err
		}
		token = t
	} else {
		token = "null"
	}
	qs := url.Values{}
	qs.Add("api_version", "1.0")
	qs.Add("api_token", token)
	qs.Add("input", "3")
	qs.Add("method", string(method))
	req.URL.RawQuery = qs.Encode()
	req.AddCookie(&http.Cookie{Name: "arl", Value: c.Arl})
	r, e := c.Do(req)
	return r, e
}

func (c *Client) apiDoJSON(method apiMethod, body interface{}, v interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r := bytes.NewBuffer(b)
	resp, err := c.apiDo(method, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var jresp response
	err = dec.Decode(&jresp)
	if err != nil {
		return err
	}
	return json.Unmarshal(jresp.Results, &v)
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) csrfToken() (string, error) {
	var udata userData
	err := c.apiDoJSON(getUserData, nil, &udata)
	return udata.CheckForm, err
}

// Song fetches a Song.
func (c *Client) Song(id string) (Song, error) {
	var song Song
	body := songGetDataBody{id}
	err := c.apiDoJSON(songGetData, body, &song)
	return song, err
}

// Album fetches an Album.
func (c *Client) Album(id string) (*Album, error) {
	var album Album
	body := albumGetDataBody{id}
	err := c.apiDoJSON(albumGetData, body, &album)
	return &album, err
}

// Artist fetches an Artist.
func (c *Client) Artist(id string) (*Artist, error) {
	var artist Artist
	body := artistGetDataBody{id}
	err := c.apiDoJSON(artistGetData, body, &artist)
	return &artist, err
}

// SongsByAlbum fetches up to songLimit songs on an album.
// If you want to fetch all of the songs on an album,
// use a songLimit of -1.
func (c *Client) SongsByAlbum(id string, songLimit int) ([]Song, error) {
	var songs multiSongResponse
	body := songListByAlbumBody{id, songLimit}
	err := c.apiDoJSON(songListByAlbum, body, &songs)
	return songs.Data, err
}

// AlbumsBy fetches albums in an artist's discography.
func (c *Client) AlbumsByArtist(id string) ([]Album, error) {
	var albums multiAlbumResponse
	body := albumGetDiscographyBody{id, "us", []int{0}, 500, 300, 0}
	err := c.apiDoJSON(albumGetDiscography, body, &albums)
	return albums.Data, err
}

// AvailableQualities returns the available qualities for download
// of a song.
func (c *Client) AvailableQualities(song Song) []Quality {
	var qualities []Quality
	if c.IsQualityAvailable(song, MP3128) {
		qualities = append(qualities, MP3128)
	}
	if c.IsQualityAvailable(song, MP3320) {
		qualities = append(qualities, MP3320)
	}
	if c.IsQualityAvailable(song, FLAC) {
		qualities = append(qualities, FLAC)
	}
	return qualities
}

// IsQualityAvailable returns whether or not a song is available
// to download for a song.
func (c *Client) IsQualityAvailable(song Song, quality Quality) bool {
	url := SongDownloadURL(song, quality)
	if url == "" {
		return false
	}
	resp, err := c.Get(url)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}
