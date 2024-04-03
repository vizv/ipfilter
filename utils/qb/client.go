package qb

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	BaseURL string
	http.Client
}

func NewClient(webUIURL *url.URL) (*Client, error) {
	if webUIURL == nil {
		return nil, nil
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %v", err)
	}

	client := http.Client{Jar: jar}

	username := webUIURL.User.Username()
	password, passwordSet := webUIURL.User.Password()
	webUIURL.User = nil
	baseURL := webUIURL.String()
	if passwordSet && username != "" && password != "" {
		log.Infof("credentials found, performing login...")
		loginData := url.Values{}
		loginData.Set("username", username)
		loginData.Set("password", password)

		loginResp, err := client.PostForm(baseURL+"/api/v2/auth/login", loginData)
		if err != nil {
			return nil, fmt.Errorf("error logging in: %v", err)
		}
		defer loginResp.Body.Close()
	} else {
		log.Infof("credentials not found, will not login")
	}

	return &Client{baseURL, client}, nil
}

func (c *Client) GetPreferences() (string, error) {
	prefResp, err := c.Get(c.BaseURL + "/api/v2/app/preferences")
	if err != nil {
		return "", fmt.Errorf("error getting preferences: %v", err)
	}
	defer prefResp.Body.Close()

	prefBody, err := io.ReadAll(prefResp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing preferences: %v", err)
	}
	if prefResp.StatusCode != 200 {
		return "", fmt.Errorf("error getting preferences: %d - %s", prefResp.StatusCode, string(prefBody))
	}
	return string(prefBody), nil
}

func (c *Client) RefreshIPFilter(newFilterPath string) error {
	setPrefData := url.Values{}
	setPrefData.Set("json", fmt.Sprintf(`{"ip_filter_enabled":true,"ip_filter_path":"%s"}`, newFilterPath))
	setPrefResp, err := c.PostForm(c.BaseURL+"/api/v2/app/setPreferences", setPrefData)
	if err != nil {
		return fmt.Errorf("error setting preferences: %v", err)
	}
	defer setPrefResp.Body.Close()

	return nil
}
