package sender

import (
	"fmt"
	"github.com/levigross/grequests"
)

type Client struct {
	BaseURL string
}

func NewClient(url string) *Client {
	return &Client{url}
}

func (c *Client) SendRequest(method, endpoint string, data interface{}) (
	*grequests.Response,
	error,
) {
	var err error
	var resp *grequests.Response
	ro := &grequests.RequestOptions{JSON: data}
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	switch method {
	case "POST":
		resp, err = grequests.Post(url, ro)
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}
