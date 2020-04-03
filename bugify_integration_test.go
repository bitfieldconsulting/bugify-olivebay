// +build integration

package bugify_test

import (
	"os"
	"testing"

	"github.com/olivebay/bugify-olivebay"
)

func getAPIKey(t *testing.T) string {
	key := os.Getenv("GITHUB_API_KEY")
	if key == "" {
		t.Fatal("'GITHUB_API_KEY' must be set for integration tests")
	}
	return key
}

func TestCreateIssue(t *testing.T) {
	issue := bugify.Issue{
		GitHubAPIKey: getAPIKey(t),
		Repo:         "bitfieldconsulting/bugify-olivebay",
		Title:        "Error in user program",
		Body:         "Bad stuff happened internally",
	}

	_, err := bugify.Open(issue)
	if err != nil {
		t.Fatalf("failed to create GitHub issue: %v", err)
	}
}
