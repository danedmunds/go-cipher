package main

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Keyword(keyword string) Cipher {
	// drop all accents and map letter that look the same, keep only alpha runes
	t := transform.Chain(
		norm.NFKD,
		runes.Map(func(r rune) rune {
			return unicode.ToUpper(r)
		}),
		runes.Remove(setFunc(func(r rune) bool {
			for _, o := range alpha {
				if o == r {
					return false
				}
			}
			return true
		})),
	)
	sanitized, _, err := transform.String(t, keyword)
	if err != nil {
		panic(err)
	}

	lookup := make(map[rune]rune)

	toRunes := append([]rune(sanitized), alpha...)
	alphaIndex := 0
	mapped := make(map[rune]bool)
	for _, to := range toRunes {
		if _, there := mapped[to]; !there {
			from := alpha[alphaIndex]
			lookup[from] = to
			mapped[to] = true

			alphaIndex++
			if alphaIndex >= len(alpha) {
				break
			}
		}
	}

	return NewSubstitutionCipher(lookup)
}
