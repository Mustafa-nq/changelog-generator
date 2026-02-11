package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// AIClient handles communication with the AI API
type AIClient struct {
	APIKey string
	Model  string
}

func NewAIClient() (*AIClient, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable not set")
	}
	return &AIClient{
		APIKey: apiKey,
		Model:  "claude",
	}, nil
}

// ImproveCommitMessage uses AI to make a commit message more descriptive
func (c *AIClient) ImproveCommitMessage(commit *Commit) (string, error) {
	// Create the prompt
	prompt := fmt.Sprintf(`You are helping to create a changelog. 

Given this git commit message: "%s"

Please improve it to be:
1. Clear and user-friendly (for non-technical users)
2. Focused on WHAT changed, not HOW
3. One sentence, under 80 characters
4. Start with a capital letter

Just respond with the improved message, nothing else.`, commit.Message)

	// Make API request
	response, err := c.callClaude(prompt)
	if err != nil {
		// If API fails, return original message
		return cleanCommitMessageAI(commit.Message), nil
	}

	return response, nil
}

// callClaude makes a request to Claude API
func (c *AIClient) callClaude(prompt string) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	requestBody := map[string]interface{}{
		"model":      c.Model,
		"max_tokens": 100,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract text from response
	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("unexpected response format")
	}

	firstContent, ok := content[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return text, nil
}

// ImproveAllCommits improves all commit messages using AI
func (c *AIClient) ImproveAllCommits(commits []*Commit) {
	fmt.Println("Using AI to improve commit messages...")
	fmt.Println()

	for i, commit := range commits {
		fmt.Printf("  Processing %d/%d: %s\n", i+1, len(commits), commit.Hash)

		improved, err := c.ImproveCommitMessage(commit)
		if err != nil {
			fmt.Printf("  Error: %v (using original)\n", err)
			continue
		}

		// Update the commit message
		commit.Message = improved
	}

	fmt.Println()
	fmt.Println(" AI processing complete!")
	fmt.Println()
}

// cleanCommitMessageAI is a simple version for when API fails
func cleanCommitMessageAI(msg string) string {
	// Remove trailing newlines
	cleaned := ""
	for _, c := range msg {
		if c == '\n' {
			break
		}
		cleaned += string(c)
	}
	return cleaned
}
