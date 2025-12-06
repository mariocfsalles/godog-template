package snapshot

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func AssertSnapshot[T any](fixturePath string, gotJSON []byte, normalize func(*T)) error {
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		return fmt.Errorf("error reading fixture %s: %w", fixturePath, err)
	}

	var expected, got T

	if err := json.Unmarshal(data, &expected); err != nil {
		return fmt.Errorf("error unmarshaling fixture %s: %w", fixturePath, err)
	}
	if err := json.Unmarshal(gotJSON, &got); err != nil {
		return fmt.Errorf("error unmarshaling API response: %w", err)
	}

	if normalize != nil {
		normalize(&expected)
		normalize(&got)
	}

	diff := cmp.Diff(
		expected,
		got,
		cmpopts.EquateEmpty(),
	)

	if diff != "" {
		return fmt.Errorf("snapshot mismatch (-expected +got):\n%s", diff)
	}

	return nil
}