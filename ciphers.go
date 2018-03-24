package main

import (
	"fmt"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

var alpha = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

type Cipher interface {
	Encipher() transform.Transformer
	Decipher() transform.Transformer
}

type (
	setFunc func(rune) bool
)

func (s setFunc) Contains(r rune) bool {
	return s(r)
}

type SubstitutionCipher struct {
	lookup        map[rune]rune
	reverseLookup map[rune]rune
}

func (sc *SubstitutionCipher) Encipher() transform.Transformer {
	return lookupMapper(sc.lookup)
}

func (sc *SubstitutionCipher) Decipher() transform.Transformer {
	if sc.reverseLookup == nil {
		sc.reverseLookup = make(map[rune]rune)
		for k, v := range sc.lookup {
			sc.reverseLookup[v] = k
		}
	}
	return lookupMapper(sc.reverseLookup)
}

func lookupMapper(lookup map[rune]rune) transform.Transformer {
	return runes.Map(func(r rune) rune {
		if mapped, ok := lookup[r]; ok {
			return mapped
		}

		return r
	})
}

func NewSubstitutionCipher(lookup map[rune]rune) *SubstitutionCipher {
	return &SubstitutionCipher{
		lookup: lookup,
	}
}

func PrintLookup(lookup map[rune]rune) {
	for k, v := range lookup {
		fmt.Printf("%s -> %s\n", string(k), string(v))
	}
}
