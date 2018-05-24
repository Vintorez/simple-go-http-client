package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type client struct {
	log         ILogger
	httpClient  *http.Client
	hostUrl     *url.URL
	baseUrlPath string
	userAgent   string
	username    string
	password    string
}

func (c *client) GET(urlStr string, v interface{}) error {
	return c.doRequest("GET", urlStr, nil, v)
}

func (c *client) PUT(urlStr string, body io.Reader, v interface{}) error {
	return c.doRequest("PUT", urlStr, body, v)
}

func (c *client) POST(urlStr string, body io.Reader, v interface{}) error {
	return c.doRequest("POST", urlStr, body, v)
}

func (c *client) DELETE(urlStr string, body io.Reader) error {
	return c.doRequest("DELETE", urlStr, nil, nil)
}

func (c *client) doRequest(method, urlPath string, body io.Reader, v interface{}) error {
	r, err := c.newRequest(method, fmt.Sprintf("%s%s", c.baseUrlPath, urlPath), body)
	if err != nil {
		return err
	}

	return c.do(r, v)
}

func (c *client) newRequest(method, ref string, body io.Reader) (*http.Request, error) {
	fullUrl, err := c.hostUrl.Parse(ref)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method, fullUrl.String(), body)
	if err != nil {
		return nil, err
	}

	r.SetBasicAuth(c.username, c.password)

	userAgent := "GoHttpClientAPI"
	if len(c.userAgent) != 0 {
		userAgent = c.userAgent
	}
	r.Header.Add("User-Agent", userAgent)
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Accept-Charset", "utf-8")
	if method == "POST" || method == "PUT" {
		r.Header.Add("Content-Type", "application/json; charset=utf-8")
	}

	return r, nil
}

func (c *client) do(r *http.Request, v interface{}) error {
	if c.log != nil {
		dumpRequest, _ := httputil.DumpRequest(r, true)
		c.log.Print(string(dumpRequest))
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.log != nil {
		dumpResponse, _ := httputil.DumpResponse(resp, true)
		c.log.Print(string(dumpResponse))
	}

	if c := resp.StatusCode; c < 200 || 299 < c {
		return fmt.Errorf("%d", resp.StatusCode)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return err
}
