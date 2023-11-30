package util

import "unicode"

func HasNonPrintableCharacters(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}

func RemoveNonPrintableCharacters(s string) string {
	if HasNonPrintableCharacters(s) {
		s = removeNonPrintableCharacters(s)
	}
	return s
}

func removeNonPrintableCharacters(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.IsPrint(r) {
			result = append(result, r)
		}
	}
	return string(result)
}
