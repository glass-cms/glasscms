package auth_test

import (
	"testing"

	"github.com/glass-cms/glasscms/internal/auth"
)

func TestGenerateRandomString(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := auth.GenerateRandomString()

		expectedLen := 20
		if len(result) != expectedLen {
			t.Errorf("Expected length of %d, but got %d", expectedLen, len(result))
		}

		for _, char := range result {
			if char < 33 || char > 126 {
				t.Errorf("Generated string contains invalid character: %c", char)
			}
		}

		other := auth.GenerateRandomString()
		if result == other {
			t.Error("Two consecutive generations produced identical strings, which is highly improbable")
		}
	}
}
