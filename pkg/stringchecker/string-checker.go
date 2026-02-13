package stringchecker

import (
	"unicode"
)

func CheckFirstLowercase(s string) (string, bool) {
	runes := []rune(s)
	if len(runes) == 0 {
		return "", true
	}
	for i, r := range runes {
		if unicode.IsLetter(r) {
			if unicode.IsLower(r) {
				return s, true
			}
			runes[i] = unicode.ToLower(r)
			return string(runes), false
		}
	}
	return s, true
}

func CheckEnglish(s string) bool {
	for _, r := range s {

		if !unicode.IsLetter(r) {
			continue
		}

		if !unicode.In(r, unicode.Latin) {
			return false
		}
	}
	return true
}

func CheckSpecial(s string, ignore map[rune]any) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}

		if ignore == nil || ignore[r] == nil {
			return false
		}
	}
	return true
}

func CheckSensitive(s string, keys []string) bool {
	if keys == nil {
		return true
	}

	keysMap := make(map[string]any, len(keys))
	for _, key := range keys {
		for _, word := range splitWord(key) {
			keysMap[word] = true
		}
	}

	for _, word := range splitWord(s) {
		if keysMap[word] != nil {
			return false
		}
	}
	return true
}

func splitWord(s string) []string {
	var words []string
	var current []rune
	runs := []rune(s)
	for i, r := range runs {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			if len(current) > 0 {
				words = append(words, string(current))
				current = nil
			}
			continue
		}

		if i > 0 && unicode.IsLower(runs[i-1]) && unicode.IsUpper(runs[i]) && len(current) > 0 {
			words = append(words, string(current))
			current = nil
		}

		if !unicode.IsLetter(r) {
			continue
		}

		current = append(current, unicode.ToLower(r))
	}

	if len(current) > 0 {
		words = append(words, string(current))
	}

	return words
}
