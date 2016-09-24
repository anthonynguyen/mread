package manga

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

// -----------------------------------------------------------------------------

func stripNonAlphanumeric(r rune) rune {
	if !unicode.IsLetter(r) {
		return -1
	}
	return unicode.ToLower(r)
}

func fuzzy(str1 string, str2 string) bool {
	str1 = strings.Map(stripNonAlphanumeric, str1)
	str2 = strings.Map(stripNonAlphanumeric, str2)

	return strings.Contains(str2, str1)
}

// -----------------------------------------------------------------------------

func relTime(t float64) string {
	diff := time.Now().Unix() - int64(t)
	days := diff / 60 / 60 / 24

	if days < 1 {
		return "Today (" + getDate(t) + ")"
	} else if days < 2 {
		return "Yesterday (" + getDate(t) + ")"
	} else {
		return fmt.Sprintf("%v days ago (%v)", days, getDate(t))
	}
}

func getDate(t float64) string {
	return time.Unix(int64(t), 0).Format("2006-01-02")
}
