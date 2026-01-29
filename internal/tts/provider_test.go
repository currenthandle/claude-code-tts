package tts

import (
	"testing"
)

func TestNewProvider_DefaultOpenAI(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	// Ensure TTS_PROVIDER is not set (use empty to clear it)
	t.Setenv("TTS_PROVIDER", "")

	provider, err := NewProvider()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.Name() != "openai" {
		t.Errorf("expected openai provider, got %s", provider.Name())
	}
}

func TestNewProvider_ExplicitOpenAI(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("TTS_PROVIDER", "openai")

	provider, err := NewProvider()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.Name() != "openai" {
		t.Errorf("expected openai provider, got %s", provider.Name())
	}
}

func TestNewProvider_Azure(t *testing.T) {
	t.Setenv("TTS_PROVIDER", "azure")
	t.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com")
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "tts-deployment")

	provider, err := NewProvider()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.Name() != "azure" {
		t.Errorf("expected azure provider, got %s", provider.Name())
	}
}

func TestNewProvider_AzureMissingEndpoint(t *testing.T) {
	t.Setenv("TTS_PROVIDER", "azure")
	t.Setenv("AZURE_OPENAI_ENDPOINT", "")
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "tts-deployment")

	_, err := NewProvider()
	if err == nil {
		t.Error("expected error for missing AZURE_OPENAI_ENDPOINT")
	}
}

func TestNewProvider_AzureMissingAPIKey(t *testing.T) {
	t.Setenv("TTS_PROVIDER", "azure")
	t.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com")
	t.Setenv("AZURE_OPENAI_API_KEY", "")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "tts-deployment")

	_, err := NewProvider()
	if err == nil {
		t.Error("expected error for missing AZURE_OPENAI_API_KEY")
	}
}

func TestNewProvider_AzureMissingDeployment(t *testing.T) {
	t.Setenv("TTS_PROVIDER", "azure")
	t.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com")
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT", "")

	_, err := NewProvider()
	if err == nil {
		t.Error("expected error for missing AZURE_OPENAI_DEPLOYMENT")
	}
}

func TestNewProvider_UnknownProvider(t *testing.T) {
	t.Setenv("TTS_PROVIDER", "unknown")

	_, err := NewProvider()
	if err == nil {
		t.Error("expected error for unknown provider")
	}
}
