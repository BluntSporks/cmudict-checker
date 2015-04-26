// Check that the phonemes in a dictionary file are valid symbols.
package main

import (
	"flag"
	"fmt"
	"github.com/BluntSporks/cmudict"
)

func main() {
	// Parse flags.
	dictFile := flag.String("d", cmudict.DefaultDictPath(), "Name of CMU-formatted file to check")
	symFile := flag.String("s", cmudict.DefaultSymbolPath(), "Name of CMU-formatted symbol file to use for checking")
	flag.Parse()

	// Load cmuDict.
	cmuDict := cmudict.LoadDict(*dictFile)

	// Load symbols.
	cmuSymbols := cmudict.LoadSymbols(*symFile, true)

	// Print any row that has bad phonemes.
	for word, pron := range cmuDict {
		phonemes := cmudict.GetPhonemes(pron, true)
		for _, phoneme := range phonemes {
			if !cmuSymbols[phoneme] {
				fmt.Printf("%s  %s\n", word, pron)
				break
			}
		}
	}
}
