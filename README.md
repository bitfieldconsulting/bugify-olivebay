# bugify-olivebay
bugify-olivebay created by GitHub Classroom

`bugify` is a Go library that allows you to auto-generate a GitHub issue to report a bug using the GitHub API.

## Using the Go library

Import the library using:

```
github.com/olivebay/bugify-olivebay
```

## Creating a client:

Create a new Client object by calling `bugify.NewClient()` with your API key and repository: 

```
apiKey := "f4c2d78215eb95a52bed3cbcd0fbcda49aa2a1a7"
client := bugify.NewClient(apikey, "test_owner/test_repo")
```

Or read the key from an environment variable:

```
client := bugify.NewClient(os.Getenv("GITHUB_API_KEY"), "test_owner/test_repo")
```

## Creating an issue:

Will create a GitHub issue containg the error reported.

```
if err != nil {
  client.Create("something went wrong: %w", err)
}
```


