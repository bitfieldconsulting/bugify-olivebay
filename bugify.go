package bugify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var GitHubAPIURL = "https://api.github.com"
var HTTPClient = &http.Client{Timeout: 10 * time.Second}

// Issue represents information about an GitHub issue.
type Issue struct {
	Title        string
	Body         string
	GitHubAPIKey string
	Repo         string
}

type apiRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type apiResponse struct {
	ID  int64  `json:"id"`
	URL string `json:"html_url"`
}

// Create creates a new GitHub issue.
func Open(issue Issue) (URL string, err error) { // always add a variable on what it returns. Its more clear
	data, err := json.Marshal(apiRequest{
		Title: issue.Title,
		Body:  issue.Body,
	})
	if err != nil {
		return "", err
	}
	requestURL := GitHubAPIURL + "/repos/" + issue.Repo + "/issues"
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(data))
	if err != nil {
		return "", err // no need for extra text
	}

	req.Header.Add("Authorization", "Token "+issue.GitHubAPIKey)
	req.Header.Add("content-type", "application/json")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected response: %d %q", resp.StatusCode, string(res))
	}

	var result apiResponse
	err = json.NewDecoder(bytes.NewReader(res)).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error decoding: %s %v", res, err)
	}
	return result.URL, nil
}
