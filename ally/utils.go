package ally

import (
	"time"

	"github.com/joho/godotenv"
)

func TruncateToMarketStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 13, 0, 0, 0, t.UTC().Location())
}
func TruncateToMarketEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 20, 0, 0, 0, t.UTC().Location())
}
func TrySetEnv() {
	godotenv.Load(".devcontainer/.env")
}
