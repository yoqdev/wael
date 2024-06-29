package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

func main() {
	// Define constants
	const githubToken = "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" // Replace with your GitHub token
	const repoOwner = "yoqdev"                                     // Replace with your GitHub username
	const repoName = "wael"                                        // Replace with your GitHub repository name
	const filePath = "/home/w/GoSend/MyFiles.zip"                  // Replace with the path to your local file
	const branch = "main"                                          // or the branch you want to upload to
	const commitMessage = "Add MyFiles.zip"                        // Commit message

	// Read the file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return
	}

	// Encode the file content to base64
	contentBase64 := base64.StdEncoding.EncodeToString(content)

	// Create the URL for the GitHub API endpoint
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, strings.TrimPrefix(filePath, "/"))

	// Set up the request
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Create the JSON payload
	payload := fmt.Sprintf(`{
		"message": "%s",
		"content": "%s",
		"branch": "%s"
	}`, commitMessage, contentBase64, branch)

	req, err := http.NewRequest("PUT", url, strings.NewReader(payload))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", githubToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Execute the request
	resp, err := tc.Do(req)
	if err != nil {
		fmt.Printf("Failed to make request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("File uploaded successfully.")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to upload file: %s\n", body)
	}
}
