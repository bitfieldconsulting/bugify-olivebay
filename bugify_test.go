package bugify_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/olivebay/bugify-olivebay"
)

func TestOpenGitHubIssue(t *testing.T) {
	called := false
	issue := bugify.Issue{
		GitHubAPIKey: "dummy",
		Repo:         "olivebay/urlinfo2",
		Title:        "Error in user program",
		Body:         "Bad stuff happened internally",
	}
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		wantMethod := http.MethodPost
		if !cmp.Equal(r.Method, wantMethod) {
			t.Error(cmp.Diff(r.Method, wantMethod))
		}

		wantURL := fmt.Sprintf("/repos/%s/issues", issue.Repo)
		if !cmp.Equal(wantURL, r.URL.String()) {
			t.Error(cmp.Diff(wantURL, r.URL.Path))
		}

		wantData := fmt.Sprintf("{\"title\":%q,\"body\":%q}", issue.Title, issue.Body)
		gotData, _ := ioutil.ReadAll(r.Body) // dont leave errors!
		if !cmp.Equal(wantData, string(gotData)) {
			t.Error(cmp.Diff(wantData, string(gotData)))
		}

		w.WriteHeader(http.StatusCreated)
		file, err := os.Open("testdata/response.json")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		_, err = io.Copy(w, file)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	bugify.HTTPClient = ts.Client()
	bugify.GitHubAPIURL = ts.URL
	wantIssueURL := "https://github.com/olivebay/urlinfo2/issues/6/"
	issueURL, err := bugify.Open(issue)
	if err != nil {
		t.Fatalf("failed to create GitHub issue %v", err)
	}
	if !cmp.Equal(wantIssueURL, issueURL) {
		t.Error(cmp.Diff(wantIssueURL, issueURL))
	}
	if !called {
		t.Error("want API call, but got none")
	}
}
