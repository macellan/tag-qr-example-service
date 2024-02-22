package routes

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func hash(data []string) string {
	s := strings.Join(data, "")

	hasher := sha1.New()
	hasher.Write([]byte(s))

	return hex.EncodeToString(hasher.Sum(nil))
}
