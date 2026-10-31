package main

import (
	"strings"
)

func censorBadWord(chirp string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	splitSentence := strings.Split(chirp, " ")
	for i, word := range splitSentence {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			splitSentence[i] = censor
		}
	}
	return strings.Join(splitSentence, " ")
}
