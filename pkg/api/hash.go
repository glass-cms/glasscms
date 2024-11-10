package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// HashItem generates a SHA-256 hash for the given Item. It hashes the content,
// properties, and metadata of the item. The function returns the hexadecimal
// representation of the hash or an error if any step of the hashing process fails.
func HashItem(content string, properties map[string]interface{}, metadata map[string]interface{}) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(content))

	propertiesJSON, err := json.Marshal(properties)
	if err != nil {
		return "", fmt.Errorf("failed to serialize properties: %w", err)
	}
	hasher.Write(propertiesJSON)

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to serialize metadata: %w", err)
	}
	hasher.Write(metadataJSON)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
