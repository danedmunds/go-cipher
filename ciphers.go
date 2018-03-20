package main

import (
	"fmt"

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
