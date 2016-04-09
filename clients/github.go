package client

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client holds information for pinging the API
type Client struct {
	Token string
	URL   url.URL
}

// GithubUser holds user information from Github
type GithubUser struct {
	Name string
}

// NewClient returns a new instance of Client
func NewClient(host string) *Client {
	c := &Client{}
	c.URL.Host = host
	return c
}

// GithubGetUser gets a user from the API
func (c *Client) GithubGetUser(token string) (string, error) {
	// Build URL
	u := c.URL
	u.Scheme = "https"
	u.Path = "/user"

	// Build the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "token "+token)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	// Parse the response
	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

// GithubGetGists gets a user from the API
func (c *Client) GithubGetGists(token, username string) (string, error) {
	// Build URL
	u := c.URL
	u.Scheme = "https"
	u.Path = "/users/" + username + "/gists"

	// Build the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "token "+token)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	// Parse the response
	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(j), nil
}
