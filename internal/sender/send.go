package sender

import (
	"fmt"
	"github.com/levigross/grequests"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

func NewClient(baseURL string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Jar: jar,
		},
	}
}

func (c *Client) SendRequest(
	method,
	endpoint string,
	data interface{},
	token string,
) (*grequests.Response, error) {
	var err error
	var resp *grequests.Response
	ro := &grequests.RequestOptions{
		JSON:       data,
		HTTPClient: c.HTTPClient,
	}
	if token != "" {
		ro.Headers = map[string]string{
			"Authorization": "Bearer " + token,
		}
	}

	urlApp := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	switch method {
	case "POST":
		resp, err = grequests.Post(urlApp, ro)
	case "GET":
		resp, err = grequests.Get(urlApp, ro)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, err
	}
	h := resp.Header.Get("Authorization")
	c.AuthToken = strings.TrimPrefix(h, "Bearer ")
	return resp, nil
}
