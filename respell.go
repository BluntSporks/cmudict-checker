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
	dictFile := flag.String("dict", cmudict.DefaultDictPath(), "Name of CMU-formatted file to modify")
	spellFile := flag.String("spell", defaultSpellPath(), "Name of spelling file that maps phonemes to spellings")
	textFile := flag.String("text", "", "Name of text file to respell")
	flag.Parse()

	if len(*textFile) == 0 {
		log.Fatal("Missing -text argument")
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
			isUpper := word != strings.ToLower(word)
			pron := cmuDict[upword]
			if pron == "" {
				fmt.Print(word)
			} else {
				phonemes := cmudict.GetPhonemes(pron, true)
				fixed := fixPhonemes(phonemes)
				for i, phoneme := range fixed {
					bare := cmudict.StripAccent(phoneme)
					spelling := spellings[bare]
					output := phoneme
					if spelling != "" {
						output = spelling
					}

					// Capitalize the first letter of output if the original word was not
					// lowercase.
					if i == 0 && isUpper {
						capd := strings.ToUpper(output[0:1])
						if len(output) > 1 {
							capd += output[1:]
						}
						output = capd
					}
					fmt.Print(output)
				}
			}
		}
		fmt.Println()
	}
}
