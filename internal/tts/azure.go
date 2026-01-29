package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// AzureClient handles Azure OpenAI TTS API requests
type AzureClient struct {
	endpoint   string // e.g., "https://my-resource.openai.azure.com"
	apiKey     string
	deployment string
	apiVersion string
	httpClient *http.Client
}

// NewAzureClient creates a new Azure TTS client from environment variables
// Required: AZURE_OPENAI_ENDPOINT, AZURE_OPENAI_API_KEY, AZURE_OPENAI_DEPLOYMENT
// Optional: AZURE_OPENAI_API_VERSION (defaults to "2024-02-15-preview")
func NewAzureClient() (*AzureClient, error) {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	deployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT")

	if endpoint == "" {
		return nil, fmt.Errorf("AZURE_OPENAI_ENDPOINT is required for Azure provider")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("AZURE_OPENAI_API_KEY is required for Azure provider")
	}
	if deployment == "" {
		return nil, fmt.Errorf("AZURE_OPENAI_DEPLOYMENT is required for Azure provider")
	}

	apiVersion := os.Getenv("AZURE_OPENAI_API_VERSION")
	if apiVersion == "" {
		apiVersion = "2024-02-15-preview"
	}

	return &AzureClient{
		endpoint:   endpoint,
		apiKey:     apiKey,
		deployment: deployment,
		apiVersion: apiVersion,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Name returns the provider name
func (c *AzureClient) Name() string {
	return "azure"
}

// azureTTSRequest represents the Azure API request payload
type azureTTSRequest struct {
	Model string  `json:"model"`
	Input string  `json:"input"`
	Voice string  `json:"voice"`
	Speed float64 `json:"speed,omitempty"`
}

// Synthesize converts text to speech using Azure OpenAI
// If speed is 0, DefaultSpeed (1.0) is used
func (c *AzureClient) Synthesize(text string, voice Voice, speed float64) ([]byte, error) {
	effectiveSpeed := speed
	if effectiveSpeed == 0 {
		effectiveSpeed = DefaultSpeed
	}

	reqBody := azureTTSRequest{
		Model: "tts-1",
		Input: text,
		Voice: string(voice),
		Speed: effectiveSpeed,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Azure endpoint format: {endpoint}/openai/deployments/{deployment}/audio/speech?api-version={version}
	url := fmt.Sprintf("%s/openai/deployments/%s/audio/speech?api-version=%s",
		c.endpoint, c.deployment, c.apiVersion)

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Azure uses api-key header instead of Authorization Bearer
	req.Header.Set("api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Azure API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Azure API error (status %d): %s", resp.StatusCode, string(body))
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return audioData, nil
}
