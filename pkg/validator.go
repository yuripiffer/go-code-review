package pkg

import "unicode"

// IsAlphaNumericOnly returns true if the input contains only letters and digits.
func IsAlphaNumericOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
