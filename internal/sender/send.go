package sender

import (
	"fmt"
	"github.com/levigross/grequests"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
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

func (c *Client) sendRequest(
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

func (c *Client) SendRequest(
	method,
	endpoint string,
	data interface{},
	token string,
) (*grequests.Response, error) {
	retryDelays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	maxRetries := 3

	var resp *grequests.Response
	var err error

	for retry := 0; retry < maxRetries; retry++ {
		fmt.Printf("Attempt %d/%d\n", retry+1, maxRetries)

		if retry > 0 {
			time.Sleep(retryDelays[retry-1])
		}

		resp, err = c.sendRequest(method, endpoint, data, token)

		if err != nil {
			fmt.Printf("Error on attempt %d: %v\n", retry+1, err)

			if retry == maxRetries-1 {
				return nil, fmt.Errorf("error sending request: %v", err)
			}
			continue
		}

		fmt.Printf("Success on attempt %d\n", retry+1)
		break
	}
	return resp, err
}
