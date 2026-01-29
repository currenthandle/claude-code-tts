package tts

import (
	"fmt"
	"os"
)

// Provider defines the interface for TTS providers
type Provider interface {
	// Synthesize converts text to speech and returns audio data
	Synthesize(text string, voice Voice) ([]byte, error)
	// Name returns the provider name
	Name() string
}

// ProviderType represents available provider types
type ProviderType string

const (
	ProviderOpenAI ProviderType = "openai"
	ProviderAzure  ProviderType = "azure"
)

// NewProvider creates a TTS provider based on TTS_PROVIDER environment variable
// Returns OpenAI provider by default if TTS_PROVIDER is not set
func NewProvider() (Provider, error) {
	providerType := ProviderType(os.Getenv("TTS_PROVIDER"))
	if providerType == "" {
		providerType = ProviderOpenAI
	}

	switch providerType {
	case ProviderOpenAI:
		return NewClient(), nil
	case ProviderAzure:
		return NewAzureClient()
	default:
		return nil, fmt.Errorf("unknown TTS provider: %s (valid: openai, azure)", providerType)
	}
}
