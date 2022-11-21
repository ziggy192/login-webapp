package util

import "strings"

// BuildKey builds keys
func BuildKey(keys ...string) string {
	b := strings.Builder{}
	for i, key := range keys {
		b.WriteString(key)
		if i < len(keys)-1 {
			b.WriteString(":")
		}
	}
	return b.String()
}
