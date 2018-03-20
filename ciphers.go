package main

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

var alpha = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func Caesar(offset int) transform.Transformer {
	lookup := make(map[rune]rune)
	for i, r := range alpha {
		lookup[r] = alpha[(offset+i)%len(alpha)]
	}

	return runes.Map(func(r rune) rune {
		if res, ok := lookup[r]; ok {
			return res
		}

		return r
	})
}
