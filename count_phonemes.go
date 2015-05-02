// Count phonemes of given length, with or without accents, in dictionary file.
package main

import (
	"flag"
	"fmt"
	"github.com/BluntSporks/cmudict"
	"strings"
)

func main() {
	// Parse flags.
	dictFile := flag.String("dict", cmudict.DefaultDictPath(), "Name of cmu-formatted dictionary file to check")
	length := flag.Int("len", 1, "Length of phoneme sequences to count")
	accent := flag.Bool("accent", false, "Use the accent on vowel phonemes")
	flag.Parse()

	// Load cmuDict.
	cmuDict := cmudict.LoadDict(*dictFile)

	// Count the phonemes.
	counts := make(map[string]int)
	for _, pron := range cmuDict {
		phonemes := cmudict.GetPhonemes(pron, *accent)
		n := len(phonemes)
		for i := 0; i < n-*length+1; i++ {
			seq := strings.Join(phonemes[i:i+*length], " ")
			counts[seq]++
		}
	}

	// Print the sequence counts.
	for seq, count := range counts {
		fmt.Printf("%s,%d\n", seq, count)
	}
}
