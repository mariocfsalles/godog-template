package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var ulidRe = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)

const PlaceholderAnyULID = "$ANY_ULID"
const PlaceholderAnyTimestamp = "$ANY_TIMESTAMP"

// Converte qualquer número (float64 do JSON) ou string para ID textual.
func AnyToStringID(v any) (string, bool) {
	switch x := v.(type) {
	case float64:
		return fmt.Sprintf("%.0f", x), true
	case string:
		return x, true
	default:
		return "", false
	}
}

func PrettyJSON(b []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, b, "", "  "); err != nil {
		return string(b)
	}
	return buf.String()
}

// compara expected vs actual suportando placeholders como $ANY_ULID
func MatchWithPlaceholders(expected, actual any, path string) error {
	// nil
	if expected == nil || actual == nil {
		if expected == actual {
			return nil
		}
		return fmt.Errorf("mismatch at %s: expected %v, got %v", path, expected, actual)
	}

	// placeholders baseados em string
	if s, ok := expected.(string); ok {
		switch s {
		case PlaceholderAnyTimestamp:
			as, ok := actual.(string)
			if !ok {
				return fmt.Errorf("mismatch at %s: expected timestamp string, got %T (%v)", path, actual, actual)
			}
			// aceita RFC3339 ou RFC3339Nano
			if _, err := time.Parse(time.RFC3339, as); err != nil {
				if _, err2 := time.Parse(time.RFC3339Nano, as); err2 != nil {
					return fmt.Errorf("mismatch at %s: value %q is not a valid RFC3339/RFC3339Nano timestamp", path, as)
				}
			}
			return nil

		default:
			// string normal → comparação literal
			if !reflect.DeepEqual(s, actual) {
				return fmt.Errorf("mismatch at %s: expected %v, got %v", path, s, actual)
			}
			return nil
		}
	}

	// tipos básicos (float64, bool etc.)
	switch exp := expected.(type) {
	case float64, bool:
		if !reflect.DeepEqual(exp, actual) {
			return fmt.Errorf("mismatch at %s: expected %v, got %v", path, exp, actual)
		}
		return nil
	}

	// objetos (map[string]any)
	if em, ok := expected.(map[string]any); ok {
		am, ok := actual.(map[string]any)
		if !ok {
			return fmt.Errorf("mismatch at %s: expected object, got %T", path, actual)
		}

		if len(em) != len(am) {
			return fmt.Errorf("mismatch at %s: expected %d keys, got %d", path, len(em), len(am))
		}

		for k, ev := range em {
			av, ok := am[k]
			if !ok {
				return fmt.Errorf("mismatch at %s: missing field %q", path, k)
			}
			subPath := path + "." + k
			if path == "" {
				subPath = k
			}
			if err := MatchWithPlaceholders(ev, av, subPath); err != nil {
				return err
			}
		}
		return nil
	}

	// arrays
	if es, ok := expected.([]any); ok {
		as, ok := actual.([]any)
		if !ok {
			return fmt.Errorf("mismatch at %s: expected array, got %T", path, actual)
		}
		if len(es) != len(as) {
			return fmt.Errorf("mismatch at %s: expected array len %d, got %d", path, len(es), len(as))
		}
		for i := range es {
			subPath := fmt.Sprintf("%s[%d]", path, i)
			if path == "" {
				subPath = fmt.Sprintf("[%d]", i)
			}
			if err := MatchWithPlaceholders(es[i], as[i], subPath); err != nil {
				return err
			}
		}
		return nil
	}

	// fallback: comparação simples
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("mismatch at %s: expected %v, got %v", path, expected, actual)
	}
	return nil
}

func TrimBaseURL(s string) string {
	return strings.TrimRight(s, "/")
}
