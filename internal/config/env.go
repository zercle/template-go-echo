package config

import (
	"os"
	"strings"
)

// getEnv returns an environment variable value, or a default if not set.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseStringSlice parses a comma-separated string into a slice of strings.
func parseStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	for _, v := range strings.Split(s, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}
