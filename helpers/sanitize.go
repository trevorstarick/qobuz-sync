package helpers

import "strings"

func SanitizeStringToPath(str string) string {
	str = strings.ReplaceAll(str, "/", "_")
	str = strings.ReplaceAll(str, "\\", "_")
	str = strings.ReplaceAll(str, ":", "_")
	str = strings.ReplaceAll(str, "*", "_")
	str = strings.ReplaceAll(str, "?", "_")
	str = strings.ReplaceAll(str, "\"", "_")
	str = strings.ReplaceAll(str, "<", "_")
	str = strings.ReplaceAll(str, ">", "_")
	str = strings.ReplaceAll(str, "|", "_")

	str = strings.TrimSpace(str)

	return str
}
