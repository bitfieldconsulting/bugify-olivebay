// +build integration

package bugify

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func getAPIKey(t *testing.T) string {
	key := os.Getenv("GITHUB_API_KEY")
	if key == "" {
		t.Fatal("'GITHUB_API_KEY' must be set for integration tests")
	}
	return key
}

func TestCreateIssue(t *testing.T) {
	client := NewClient(getAPIKey(t), "olivebay/urlinfo2")

	exampleError := errors.New("error")
	_, err := client.Create("[auto-generated]", fmt.Sprintf("something went wrong: %s", exampleError))
	if err != nil {
		t.Fatal("failed to create GitHub issue:", err)
	}
	// add a delete endpoint to delete the test issue
}
