package bugify_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/olivebay/bugify"
)

func TestOpenGitHubIssue(t *testing.T) {
	called := false

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		if !cmp.Equal(r.Method, http.MethodPost) {
			t.Error(cmp.Diff(r.Method, http.MethodGet))
		}

		wantURL := "/repos/test_owner/test_repo/issues"
		if !cmp.Equal(wantURL, r.URL.String()) {
			t.Error(cmp.Diff(wantURL, r.URL.Path))
		}

		w.WriteHeader(http.StatusCreated)
		data := `{"id": "1"}`
		fmt.Fprintf(w, data)
	}))
	defer ts.Close()

	client := bugify.NewClient("dummy")
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
	client.Repository = "test_owner/test_repo"

	wantRes := bugify.Issue{
		Title: "[auto-generated]",
		Body:  "I'm having a problem with this.",
	}

	res, err := client.Create(wantRes)
	if err != nil {
		t.Fatal(err)
	}

	wantID := "1"
	if !cmp.Equal(res, wantID) {
		t.Error(cmp.Diff(res, wantID))
	}

	if !called {
		t.Error("want API call, but got none")
	}
}
