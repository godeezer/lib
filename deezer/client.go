package deezer

import (
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Client struct {
	*http.Client
	Arl string
}

func NewClient(arl string) *Client {
	client := &Client{
		Client: &http.Client{},
		Arl:    arl,
	}
	return client
}

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,en-US;q=0.8,en;q=0.7")
	req.AddCookie(&http.Cookie{
		Name:  "arl",
		Value: c.Arl,
	})
	return req, nil
}

func (c *Client) get(url string) (resp *http.Response, err error) {
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) getDeezerJSON(url string) (string, error) {
	resp, err := c.get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return scrapeJSON(string(body))
}

func scrapeJSON(body string) (string, error) {
	re := regexp.MustCompile(`<script>window\.__DZR_APP_STATE__ = (.*)<\/script>`)
	match := re.FindAllStringSubmatch(body, -1)
	return match[0][1], nil
}
