package action_system

import (
	"os"
)

// GetEnv returns the value of the environment variable named by the key.
func GetEnv(key string) string {
	return os.Getenv(key)
}
