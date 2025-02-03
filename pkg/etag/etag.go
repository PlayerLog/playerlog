package etag

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func CalculateETagForAvatar(url string) string {
	// Hash the URL using SHA-1
	hash := sha1.New()
	hash.Write([]byte(url))
	etag := hex.EncodeToString(hash.Sum(nil))

	// Return the ETag as a string
	return fmt.Sprintf(`"%s"`, etag)
}
