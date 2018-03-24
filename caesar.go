package main

func Caesar(offset int) Cipher {
	lookup := make(map[rune]rune)
	for i, r := range alpha {
		lookup[r] = alpha[(offset+i)%len(alpha)]
	}

	return NewSubstitutionCipher(lookup)
}
