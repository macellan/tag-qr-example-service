package routes

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func hash(sb []string) string {
	s := strings.Join(sb, "")

	hasher := sha1.New()
	hasher.Write([]byte(s))

	return hex.EncodeToString(hasher.Sum(nil))
}
