package sqlutil

import "strings"

// IsEmptyOrComment reporta si un bloque de SQL solo contiene
// líneas vacías o comentarios (--), sin sentencias ejecutables.
func IsEmptyOrComment(content string) bool {
	lines := strings.SplitSeq(content, "\n")
	for line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
			return false
		}
	}
	return true
}
