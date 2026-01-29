package tts

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewAzureClient_Success(t *testing.T) {
	t.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com")
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "tts-deployment")

	client, err := NewAzureClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.endpoint != "https://test.openai.azure.com" {
		t.Errorf("unexpected endpoint: %s", client.endpoint)
	}
	if client.apiKey != "test-key" {
		t.Error("unexpected api key")
	}
	if client.deployment != "tts-deployment" {
		t.Errorf("unexpected deployment: %s", client.deployment)
	}
	if client.apiVersion != "2024-02-15-preview" {
		t.Errorf("unexpected default api version: %s", client.apiVersion)
	}
}

func TestNewAzureClient_CustomAPIVersion(t *testing.T) {
	t.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com")
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "tts-deployment")
	t.Setenv("AZURE_OPENAI_API_VERSION", "2024-05-01")

	client, err := NewAzureClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.apiVersion != "2024-05-01" {
		t.Errorf("expected custom api version, got %s", client.apiVersion)
	}
}

func TestNewAzureClient_MissingEndpoint(t *testing.T) {
	t.Setenv("AZURE_OPENAI_ENDPOINT", "")
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "tts-deployment")

	_, err := NewAzureClient()
	if err == nil {
		t.Error("expected error for missing endpoint")
	}
	if !strings.Contains(err.Error(), "AZURE_OPENAI_ENDPOINT") {
		t.Errorf("expected error to mention AZURE_OPENAI_ENDPOINT, got: %v", err)
	}
}

func TestNewAzureClient_MissingAPIKey(t *testing.T) {
	t.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com")
	t.Setenv("AZURE_OPENAI_API_KEY", "")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "tts-deployment")

	_, err := NewAzureClient()
	if err == nil {
		t.Error("expected error for missing api key")
	}
	if !strings.Contains(err.Error(), "AZURE_OPENAI_API_KEY") {
		t.Errorf("expected error to mention AZURE_OPENAI_API_KEY, got: %v", err)
	}
}

func TestNewAzureClient_MissingDeployment(t *testing.T) {
	t.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com")
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "")

	_, err := NewAzureClient()
	if err == nil {
		t.Error("expected error for missing deployment")
	}
	if !strings.Contains(err.Error(), "AZURE_OPENAI_DEPLOYMENT") {
		t.Errorf("expected error to mention AZURE_OPENAI_DEPLOYMENT, got: %v", err)
	}
}

func TestAzureClient_Name(t *testing.T) {
	client := &AzureClient{}
	if client.Name() != "azure" {
		t.Errorf("expected 'azure', got %s", client.Name())
	}
}

func TestAzureClient_Synthesize_Success(t *testing.T) {
	expectedAudio := []byte("fake-mp3-audio-data")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Azure-specific headers
		if r.Header.Get("api-key") != "test-api-key" {
			t.Errorf("expected api-key header, got %s", r.Header.Get("api-key"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type, got %s", r.Header.Get("Content-Type"))
		}

		// Verify URL contains deployment and api-version
		if !strings.Contains(r.URL.Path, "/openai/deployments/tts-deployment/audio/speech") {
			t.Errorf("unexpected URL path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("api-version") != "2024-02-15-preview" {
			t.Errorf("unexpected api-version: %s", r.URL.Query().Get("api-version"))
		}

		// Verify request body
		body, _ := io.ReadAll(r.Body)
		var req azureTTSRequest
		json.Unmarshal(body, &req)
		if req.Voice != "nova" {
			t.Errorf("expected voice nova, got %s", req.Voice)
		}
		if req.Input != "Hello, Azure!" {
			t.Errorf("expected input 'Hello, Azure!', got %s", req.Input)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(expectedAudio)
	}))
	defer server.Close()

	client := &AzureClient{
		endpoint:   server.URL,
		apiKey:     "test-api-key",
		deployment: "tts-deployment",
		apiVersion: "2024-02-15-preview",
		httpClient: server.Client(),
	}

	audio, err := client.Synthesize("Hello, Azure!", VoiceNova)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(audio) != string(expectedAudio) {
		t.Errorf("unexpected audio data")
	}
}

func TestAzureClient_Synthesize_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid api key"}`))
	}))
	defer server.Close()

	client := &AzureClient{
		endpoint:   server.URL,
		apiKey:     "invalid-key",
		deployment: "tts-deployment",
		apiVersion: "2024-02-15-preview",
		httpClient: server.Client(),
	}

	_, err := client.Synthesize("Hello", VoiceAlloy)
	if err == nil {
		t.Error("expected error for API failure")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected error to contain status code, got: %v", err)
	}
}
