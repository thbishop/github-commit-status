package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "runtime"
)

type config struct {
	apiUrl   string
	apiToken string
}

func newConfig() *config {
    c := &config{
		apiToken: os.Getenv("GITHUB_TOKEN"),
        apiUrl: getApiUrl(),
	}

    if c.apiUrl == "" {
        fmt.Printf("Error: Invalid API URL specified '%s'", c.apiUrl)
        os.Exit(1)
    }

    if c.apiUrl == "" || c.apiToken == "" {
        fmt.Printf("Error: GITHUB_TOKEN environment variable not specified")
        os.Exit(1)
    }

    return c
}

func getApiUrl() string {
    if apiUrl := os.Getenv("GITHUB_API"); apiUrl != "" {
        return apiUrl
    }
    return "https://api.github.com"
}

type statusBody struct {
    Context string `json:"context,omitempty"`
    Description string `json:"description,omitempty"`
    State string `json:"state"`
    TargetUrl string `json:"target_url,omitempty"`
}

func statusRequestBody(o *options) ([]byte, error) {
    body := &statusBody{
        Context: o.Context,
        Description: o.Description,
        State: o.State,
        TargetUrl: o.TargetUrl,
    }

    b, err := json.Marshal(body)
    if err != nil {
        return []byte{}, err
    }

    return b, nil
}

func statusUrl(o *options, c *config) string {
    return fmt.Sprintf("%s/repos/%s/%s/statuses/%s", c.apiUrl, o.User, o.Repo, o.Commit)
}

func main() {
	options := parseCliArgs()
    config := newConfig()

    fmt.Printf("Updating status...\n")
    url := statusUrl(options, config)

    body, err := statusRequestBody(options)
    if err != nil {
        os.Stderr.Write([]byte(fmt.Sprintf("Error creating API request body: %s", err.Error())))
        os.Exit(1)
    }

    client := &http.Client{
    }
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    req.Header.Add("Authorization", "token " + config.apiToken)
    req.Header.Add("Content-Type", "application/json")
    agent := "github-commit-status/" + version + " (" + runtime.GOOS + ")"
    req.Header.Add("User-Agent", agent)

    resp, err := client.Do(req)
    if err != nil {
        os.Stderr.Write([]byte(fmt.Sprintf("Error making API call: %s", err.Error())))
        os.Exit(1)
    }

    defer resp.Body.Close()

    if resp.StatusCode < 200 && resp.StatusCode >= 300 {
        os.Stderr.Write([]byte(fmt.Sprintf("Did not receive a successful API response (received '%s')", resp.Status)))

        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            os.Stderr.Write([]byte(fmt.Sprintf("Error reading response body: %s", err)))
            os.Exit(1)
        }
        os.Stderr.Write([]byte(fmt.Sprintf("Response Body:\n%s", contents)))
        os.Exit(1)
    }

    fmt.Printf("Status update complete")
}
