# Claude Code TTS Plugin

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/currenthandle/claude-code-tts/actions/workflows/ci.yml/badge.svg)](https://github.com/currenthandle/claude-code-tts/actions/workflows/ci.yml)
[![MCP](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://modelcontextprotocol.io)

A Text-to-Speech MCP server plugin for Claude Code that converts text to speech using OpenAI or Azure OpenAI TTS APIs. Get audio feedback from Claude as you work!

## Features

- **Multiple Providers**: OpenAI and Azure OpenAI TTS support
- **Deterministic Auto-Speak**: Every Claude response is automatically spoken (via Stop hook)
- **6 High-Quality Voices**: alloy, echo, fable, onyx, nova, shimmer
- **Adjustable Speed**: 0.25x (slow) to 4.0x (fast)
- **Worker Pool Architecture**: Non-blocking queue with concurrent processing
- **Cross-Platform**: macOS (afplay), Linux (mpv/ffplay/mpg123), Windows (PowerShell)
- **Standalone CLI**: `speak-text` binary for direct TTS without MCP

## Quick Install

No Go or git required — downloads a prebuilt binary:

```bash
curl -fsSL https://raw.githubusercontent.com/currenthandle/claude-code-tts/main/install.sh | bash
```

Then configure your provider and register with Claude Code:

```bash
# Set your API key (see Configuration below)
export OPENAI_API_KEY="sk-..."

# Register with Claude Code
claude mcp add tts ~/.claude/plugins/claude-code-tts/bin/tts-server -s user
```

## Configuration

### OpenAI (default)

```bash
export OPENAI_API_KEY="sk-..."
```

### Azure OpenAI

```bash
export TTS_PROVIDER=azure
export AZURE_OPENAI_API_KEY="your-key"
export AZURE_OPENAI_ENDPOINT="https://your-resource.openai.azure.com"
export AZURE_OPENAI_DEPLOYMENT="gpt-4o-mini-tts"
export AZURE_OPENAI_API_VERSION="2025-03-01-preview"
```

### Speed Control

Set a global default speed (0.25–4.0, default 1.0):

```bash
export CLAUDE_TTS_SPEED=1.5
```

Or pass `speed` per call in the `speak` tool.

### Requirements

- **Audio Player**:
  - macOS: `afplay` (built-in)
  - Linux: `mpv`, `ffplay`, or `mpg123`
  - Windows: PowerShell (built-in)

## Usage

### speak(text, voice, speed)

Convert text to speech and play it aloud.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `text` | string | Yes | Text to speak (max 4096 chars) |
| `voice` | string | No | alloy, echo, fable, onyx, nova, shimmer (default: alloy) |
| `speed` | number | No | 0.25–4.0 (default: 1.0 or `CLAUDE_TTS_SPEED`) |

### tts_status()

Get the current status of the TTS system (queue size, processed count, recent jobs).

### speak-text CLI

Standalone binary for direct TTS without MCP:

```bash
speak-text "Hello world"
speak-text -voice onyx "Error occurred"
```

## Building from Source

Requires Go 1.21+:

```bash
git clone https://github.com/currenthandle/claude-code-tts.git
cd claude-code-tts
make build       # Creates bin/tts-server and bin/speak-text
make test        # Run tests
make install     # Install to ~/.claude/plugins/claude-code-tts/
```

## Architecture

```
Claude Code → MCP stdio → TTS Server → Worker Pool → OpenAI/Azure API → Audio Player
```

- **Worker Pool**: 2 workers, 50-slot queue. `speak()` returns immediately after queuing.
- **Mutex-Protected Playback**: One audio plays at a time.
- **Provider Interface**: Pluggable TTS backends (OpenAI, Azure OpenAI).

## License

MIT License - see [LICENSE](LICENSE) for details.
