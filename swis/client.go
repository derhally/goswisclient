package swis

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

//Client swis proxy
type Client struct {
	username string
	password string
	endpoint string
	host     string
}

type response struct {
	Results []map[string]interface{} `json:"results"`
	Message string                   `json:"message"`
}

// NewClient creates a new swis proxy
func NewClient(host, user, pass string, port int) *Client {

	endpoint := fmt.Sprintf("https://%s:%d/SolarWinds/InformationService/v3/Json", host, port)

	return &Client{
		username: user,
		password: pass,
		endpoint: endpoint,
		host:     host,
	}
}

// Query executes a swql query
func (p *Client) Query(query string) ([]map[string]interface{}, error) {

	request := fmt.Sprintf("%s/Query?query=%s", p.endpoint, url.QueryEscape(query))

	// create request
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return nil, fmt.Errorf("new request failed: %v", err)
	}

	req.SetBasicAuth(p.username, p.password)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Query failed. Server returned: %v, ", resp.StatusCode)
	}

	defer resp.Body.Close()

	var body response
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %v", err)
	}

	if body.Message != "" {
		return nil, fmt.Errorf("error executing query: %s", body.Message)
	}

	return body.Results, nil
}
