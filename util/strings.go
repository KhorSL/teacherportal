package util

import "strings"

func StringsToSingleQuoteCommaSep(emails []string) string {
	return "'" + strings.Join(emails, "','") + "'"
}
