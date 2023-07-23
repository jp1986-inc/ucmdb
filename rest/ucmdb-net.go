package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"crypto/tls"
)

type Client struct {
	client   *http.Client
	headers  *http.Header
	address  *url.URL
	username string
	password string

	Ucmdb
}

func NewClient(address string, username string, password string) (*Client, error) {
	baseURL, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	headers.Add("Content-Type", "application/json")
	headers.Add("Accept", "application/json")

	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
        client_insecure := &http.Client{Transport: tr}

	client := &Client{
		client:   client_insecure,
		headers:  &headers,
		address:  baseURL,
		username: username,
		password: password,
	}
	client.Ucmdb = &ucmdb{client}
	err = client.getToken()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) newRequest(method string, path string, v interface{}) (*http.Request, error) {
	var body io.Reader

	u, err := url.Parse(fmt.Sprintf("%s/%s", c.address.String(), path))
	if err != nil {
		return nil, err
	}
	switch method {
	case "GET", "DELETE":
		if v != nil {
			qs := ""
			for k, v := range v.(map[string]string) {
				if qs == "" {
					qs = fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v))
				} else {
					qs = fmt.Sprintf("%s&%s=%s", qs, url.QueryEscape(k), url.QueryEscape(v))
				}
			}
			u.RawQuery = qs
		}
		LogMe("DEBUG", fmt.Sprintf("ucmdbsdk.newRequest()|%s|%s", method, u.String()), "")
	case "POST", "PUT":
		if v != nil {
			b, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				LogMe("ERROR", fmt.Sprintf("ucmdbsdk.newRequest()|%s|%s", method, u.String()), v)
				return nil, fmt.Errorf("failed to marshall interface to json")
			}
			LogMe("DEBUG", fmt.Sprintf("ucmdbsdk.newRequest()|%s|%s", method, u.String()), string(b))
			body = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for hk, hv := range *c.headers {
		req.Header[hk] = hv
	}

	return req, nil
}

func (c *Client) sendRequest(req *http.Request) (interface{}, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// read body to buffer
	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		text := buffer.String()
		LogMe("ERROR", fmt.Sprintf("ucmdbsdk.sendRequest()|http status code: %d", res.StatusCode), text)
		match, _ := regexp.MatchString("(?i)Token has expired", text)
		if match && res.StatusCode == 401 {
			// refresh token
			err := c.getToken()
			if err != nil {
				return nil, err
			}
			req.Header["Authorization"] = []string{c.headers.Get("Authorization")}
			res, err = c.client.Do(req)
			if err != nil {
				return nil, err
			}
			if res.StatusCode < 200 || res.StatusCode >= 300 {
				return nil, fmt.Errorf(res.Status)
			}
			buffer.Reset()
			_, err = buffer.ReadFrom(res.Body)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf(res.Status)
		}
	}

	//defer res.Body.Close()
	res.Body.Close()
	return buffer.Bytes(), nil

}

func (c *Client) getToken() error {

	// create HTTP Request
	rt := RequestToken{
		Username:      c.username,
		Password:      c.password,
		ClientContext: "1",
	}
	req, err := c.newRequest("POST", "authenticate", rt)
	if err != nil {
		return err
	}
	b, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	data := &Token{}
	err = json.Unmarshal(b.([]byte), data)
	if err != nil {
		return err
	}

	token := data.Token
	c.headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

func (c *Client) SetAuthorizationManually(token string) error {
	c.headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}
