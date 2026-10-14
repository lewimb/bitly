package config

import (
	"fmt"
	"os"
	"time"
)

const defaultLinkTTL = 10 * time.Minute

type App struct {
	LinkTTL time.Duration
}

// NewApp reads app settings from the environment.
// LINK_TTL uses Go duration format, e.g. "10m", "1h", "168h". Defaults to 10m when unset.
func NewApp() (*App, error) {
	linkTTL := defaultLinkTTL

	if value := os.Getenv("LINK_TTL"); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil {
			return nil, fmt.Errorf("invalid LINK_TTL %q: %w", value, err)
		}
		linkTTL = parsed
	}

	return &App{LinkTTL: linkTTL}, nil
}
