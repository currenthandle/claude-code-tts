package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ybouhjira/claude-code-tts/internal/logging"
	"github.com/ybouhjira/claude-code-tts/internal/server"
)

func main() {
	// Initialize file logging
	if err := logging.Init(); err != nil {
		log.Printf("Warning: failed to initialize file logging: %v", err)
	}

	logging.Info("========================================")
	logging.Info("TTS Server Starting")
	logging.Info("Go version: %s", runtime.Version())
	logging.Info("OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logging.Info("PID: %d", os.Getpid())
	logging.Info("Log file: %s", logging.GetLogPath())
	logging.Info("========================================")

	// Determine provider and validate required environment variables
	provider := os.Getenv("TTS_PROVIDER")
	if provider == "" || provider == "openai" {
		if os.Getenv("OPENAI_API_KEY") == "" {
			logging.Fatal("OPENAI_API_KEY environment variable is required for OpenAI provider")
		}
		logging.Info("Using OpenAI TTS provider (API key length: %d)", len(os.Getenv("OPENAI_API_KEY")))
	} else if provider == "azure" {
		var missing []string
		if os.Getenv("AZURE_OPENAI_ENDPOINT") == "" {
			missing = append(missing, "AZURE_OPENAI_ENDPOINT")
		}
		if os.Getenv("AZURE_OPENAI_API_KEY") == "" {
			missing = append(missing, "AZURE_OPENAI_API_KEY")
		}
		if os.Getenv("AZURE_OPENAI_DEPLOYMENT") == "" {
			missing = append(missing, "AZURE_OPENAI_DEPLOYMENT")
		}
		if len(missing) > 0 {
			logging.Fatal("Missing required Azure environment variables: %v", missing)
		}
		logging.Info("Using Azure OpenAI TTS provider (endpoint: %s, deployment: %s)",
			os.Getenv("AZURE_OPENAI_ENDPOINT"), os.Getenv("AZURE_OPENAI_DEPLOYMENT"))
	} else {
		logging.Fatal("Unknown TTS_PROVIDER: %s (valid: openai, azure)", provider)
	}

	// Create and start the MCP server
	srv, err := server.New()
	if err != nil {
		logging.Fatal("Failed to create server: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE)

	go func() {
		sig := <-sigChan
		logging.Info("Received signal: %v", sig)
		logging.Info("Shutting down TTS server...")
		srv.Shutdown()
		logging.Info("TTS Server stopped gracefully")
		os.Exit(0)
	}()

	// Start serving
	logging.Info("Starting MCP stdio server...")
	if err := srv.Start(); err != nil {
		logging.Error("Server error: %v", err)
		logging.Fatal("Server stopped unexpectedly")
	}

	logging.Info("Server ended normally")
}
