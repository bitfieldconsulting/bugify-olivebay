package bugify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client represents a GitHub client.
// You can specify your API key and the repository in which to create the issue.
type Client struct {
	apiKey     string
	URL        string
	Repository string
	HTTPClient *http.Client
}

// Issue represents information about an GitHub issue.
type Issue struct {
	ID    int `json:"id,omitempty"`
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url,omitempty"`
}

// NewClient takes GitHub API key and the repository name in the form of :owner/:repo and returns a Client.
func NewClient(apiKey string, repo string) Client {
	return Client{
		apiKey:     apiKey,
		URL:        "https://api.github.com",
		Repository: repo,
		HTTPClient: http.DefaultClient,
	}
}
func NewIssue(title string, body string) Issue{
	return Issue{
		Title: title,
		Body: body,
	}
}

// Create creates a new GitHub issue
func (c *Client) Create(title string, body string) (Issue, error) {
	issue := NewIssue(title, body)
	data, err := json.Marshal(issue)
	if err != nil {
		return Issue{}, err
	}

	requestURL := c.URL + "/repos/" + c.Repository + "/issues"
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(data))
	if err != nil {
		return Issue{}, fmt.Errorf("failed to create GitHub issue request: %v", err)
	}
	
	req.Header.Add("Authorization", "Token "+c.apiKey)
	req.Header.Add("content-type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Issue{}, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Issue{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return Issue{}, fmt.Errorf("unexpected response: %d %q", resp.StatusCode, string(res))
	}

	var result Issue
	err = json.NewDecoder(bytes.NewReader(res)).Decode(&result)
	if err != nil {
		return Issue{}, fmt.Errorf("error decoding: %s %v", res, err)
	}

	return result, nil
}
