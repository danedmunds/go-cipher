package main

import (
	"fmt"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var alpha = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

type setFunc func(rune) bool

func (s setFunc) Contains(r rune) bool {
	return s(r)
}

func Caesar(offset int) transform.Transformer {
	lookup := make(map[rune]rune)
	for i, r := range alpha {
		lookup[r] = alpha[(offset+i)%len(alpha)]
	}

	return lookupMapper(lookup)
}

func Rot13() transform.Transformer {
	lookup := make(map[rune]rune)
	for i := 0; i < 13; i++ {
		lookup[alpha[i]] = alpha[i+13]
		lookup[alpha[i+13]] = alpha[i]
	}

	return lookupMapper(lookup)
}

func Keyword(keyword string) transform.Transformer {
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

	return lookupMapper(lookup)
}

func lookupMapper(lookup map[rune]rune) transform.Transformer {
	return runes.Map(func(r rune) rune {
		if mapped, ok := lookup[r]; ok {
			return mapped
		}

		return r
	})
}

func PrintLookup(lookup map[rune]rune) {
	for k, v := range lookup {
		fmt.Printf("%s -> %s\n", string(k), string(v))
	}
}
