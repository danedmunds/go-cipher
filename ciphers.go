package main

import (
	"bufio"
	"fmt"
	"io"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

var alpha = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

type (
	InOutFunc func(from io.Reader, to io.Writer) error
	Cipher    interface {
		Encipher(from io.Reader, to io.Writer) error
		Decipher(from io.Reader, to io.Writer) error
	}
	FullMessageBuffer interface {
		Write(rune) error
		DumpTo(io.Writer) error
		CleanUp() error
	}
)

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

func NewSubstitutionCipher(lookup map[rune]rune) *SubstitutionCipher {
	return &SubstitutionCipher{
		lookup: lookup,
	}
}

func (sc *SubstitutionCipher) Encipher(from io.Reader, to io.Writer) error {
	reader := transform.NewReader(from, lookupMapper(sc.lookup))
	_, err := io.Copy(to, reader)
	return err
}

func (sc *SubstitutionCipher) Decipher(from io.Reader, to io.Writer) error {
	if sc.reverseLookup == nil {
		sc.reverseLookup = make(map[rune]rune)
		for k, v := range sc.lookup {
			sc.reverseLookup[v] = k
		}
	}
	reader := transform.NewReader(from, lookupMapper(sc.reverseLookup))
	_, err := io.Copy(to, reader)
	return err
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

type FullMessageCipher struct {
	cipherBuffer FullMessageBuffer
}

func NewFullMessageCipher(cipher FullMessageBuffer) Cipher {
	return &FullMessageCipher{
		cipherBuffer: cipher,
	}
}

func (fmc *FullMessageCipher) Encipher(from io.Reader, to io.Writer) error {
	defer func() {
		err := fmc.cipherBuffer.CleanUp()
		if err != nil {
			panic(err)
		}
	}()

	reader := bufio.NewReader(from)

	var err error
	var r rune
	for r, _, err = reader.ReadRune(); err == nil; {
		err = fmc.cipherBuffer.Write(r)
		if err == nil {
			r, _, err = reader.ReadRune()
		}
	}

	if err != io.EOF {
		return err
	}

	err = fmc.cipherBuffer.DumpTo(to)
	if err != nil {
		return err
	}

	return nil
}

func (fmc *FullMessageCipher) Decipher(from io.Reader, to io.Writer) error {
	_, err := to.Write([]byte("Not implemented yet"))
	return err
}
