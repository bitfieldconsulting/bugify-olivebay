package bugify_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/bugify-olivebay"
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

		wantData := `{"title":"[auto-generated]","body":"something went wrong: error"}`
		gotData, _ := ioutil.ReadAll(r.Body)
		if !cmp.Equal(wantData, string(gotData)) {
			t.Error(cmp.Diff(wantData, string(gotData)))
		}

		w.WriteHeader(http.StatusCreated)
		data := `{"id": 1, "title": "[auto-generated]", "body":"something went wrong: error"}`
		fmt.Fprintf(w, data)
	}))
	defer ts.Close()

	client := bugify.NewClient("dummy", "test_owner/test_repo")
	client.HTTPClient = ts.Client()
	client.URL = ts.URL

	wantRes := bugify.Issue{
		ID:    1,
		Title: "[auto-generated]",
		Body:  "something went wrong: error",
	}

	exampleError := errors.New("error")
	res, err := client.Create("[auto-generated]", fmt.Sprintf("something went wrong: %v", exampleError))
	if err != nil {
		t.Fatal("failed to create GitHub issue", err)
	}

	if !cmp.Equal(res, wantRes) {
		t.Error(cmp.Diff(res, wantRes))
	}

	if !called {
		t.Error("want API call, but got none")
	}
}
