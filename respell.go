// Respell a text file phonetically by looking up available words.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/BluntSporks/cmudict"
	"log"
	"os"
	"strings"
)

func main() {
	// Parse flags.
	dictFile := flag.String("d", cmudict.DefaultDictPath(), "Name of CMU-formatted file to modify")
	spellFile := flag.String("s", defaultSpellPath(), "Name of spelling file that maps phonemes to spellings")
	textFile := flag.String("t", "", "Name of text file to respell")
	flag.Parse()

	if len(*textFile) == 0 {
		log.Fatal("Missing -t argument")
	}

	// Load CMUDict file.
	cmuDict := cmudict.LoadDict(*dictFile)

	// Load the spelling file.
	spellings := loadSpellings(*spellFile)

	// Open dict file.
	handle, err := os.Open(*textFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Scan file line by line.
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		words := matchWords(line)
		for _, word := range words {
			upword := strings.ToUpper(word)
			pron := cmuDict[upword]
			if pron == "" {
				fmt.Print(word)
			} else {
				phonemes := cmudict.GetPhonemes(pron, true)
				fixed := fixPhonemes(phonemes)
				for _, phoneme := range fixed {
					bare := cmudict.StripAccent(phoneme)
					spelling := spellings[bare]
					output := phoneme
					if spelling != "" {
						output = spelling
					}
					fmt.Print(output)
				}
			}
		}
		fmt.Println()
	}
}
