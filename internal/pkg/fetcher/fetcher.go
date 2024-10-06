package fetcher

import (
	"encoding/base64"
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

var (
	resourceNotAvailableError = errors.New("resource not available")
)

type Client struct {
	c *resty.Client
}

func NewClient(timeout int) *Client {
	c := resty.New().SetTimeout(time.Duration(timeout) * time.Second)
	return &Client{c: c}
}
func (c *Client) Fetch(baseUrl string) ([]byte, error) {
	available := checkResourceAvailable(baseUrl)
	if !available {
		return nil, resourceNotAvailableError
	}
	resp, err := c.c.R().Get(baseUrl)
	if err != nil {
		return nil, err
	}
	if !isBase64(resp.Body()) {
		return resp.Body(), nil
	}
	decodeLen := base64.StdEncoding.EncodedLen(len(resp.Body()))
	decoded := make([]byte, decodeLen)
	n, err := base64.StdEncoding.Decode(decoded, resp.Body())
	if err != nil {
		return nil, err
	}
	return decoded[:n], nil
}

func isBase64(data []byte) bool {
	if len(data)%4 != 0 {
		return false
	}
	for _, b := range data {
		if !((b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '+' || b == '/' || b == '=') {
			return false
		}
	}
	return true
}

func checkResourceAvailable(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true

}
