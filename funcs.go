// Provide functions for the phonetic spelling utilities.
package main

import (
	"bufio"
	"github.com/BluntSporks/cmudict"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

// defaultSpellPath gets the spelling file location from the PHONUTIL_DATA environment variable.
func defaultSpellPath() string {
	dir := os.Getenv("PHONUTIL_DATA")
	if dir == "" {
		log.Fatal("Set PHONUTIL_DATA variable to directory of spelling file")
	}
	return path.Join(dir, "spellings.csv")
}

// fixPhonemes fixes a pronunciation string with custom fixes designed for better phonetic spelling.
func fixPhonemes(phonemes []string) []string {
	n := len(phonemes)
	newPhonemes := make([]string, 0, 2*n)
	for i := 0; i < n; i++ {
		phoneme := phonemes[i]
		if len(phoneme) > 2 && phoneme[:2] == "ER" {
			// Use ahx r instead of erx.
			newPhonemes = append(newPhonemes, "AH"+string(phoneme[2]))
			phoneme = "R"
		} else if i < n-1 {
			if phoneme == "HH" && phonemes[i+1] == "W" {
				// Use wh instead of hw.
				phoneme = "WH"
				i++
			}
		}
		newPhonemes = append(newPhonemes, phoneme)
	}
	n = len(newPhonemes)
	out := make([]string, 0, 2*n)
	for i, ph := range newPhonemes {
		out = append(out, ph)
		if i < n-1 {
			// Use an apostrophe to split up ambiguous combinations of sounds.
			split := false
			if newPhonemes[i+1] == "HH" {
				if ph == "D" || ph == "S" || ph == "T" || ph == "W" || ph == "Z" {
					split = true
				}
			} else if ph == "N" && newPhonemes[i+1] == "G" {
				split = true
			} else if cmudict.IsVowel(ph) && cmudict.IsVowel(newPhonemes[i+1]) {
				split = true
			}
			if split {
				out = append(out, "'")
			}
		}
	}
	return out
}

// loadSpellings loads a file of phonemes to spellings.
func loadSpellings(file string) map[string]string {
	// Open file.
	hdl, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer hdl.Close()

	// Scan file line by line.
	spellings := make(map[string]string)
	scanner := bufio.NewScanner(hdl)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		flds := strings.Split(line, ",")
		phoneme, spelling := flds[0], flds[1]
		spellings[phoneme] = spelling
	}
	return spellings
}

// loadWordList loads a list of words in a file and returns it as an uppercased lookup list.
func loadWordList(file string) map[string]bool {
	// Open file.
	hdl, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer hdl.Close()

	// Scan file line by line.
	words := make(map[string]bool)
	scanner := bufio.NewScanner(hdl)
	for scanner.Scan() {
		word := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		word = strings.ToUpper(word)
		words[word] = true
	}
	return words
}

// matchWords matches the ASCII words in a text, including any embedded apostrophes.
// It also matches wordlike expressions that contain nonword characters, such as email or web addresses.
// Finally, it captures any other character as a single byte.
func matchWords(text string) []string {
	re := regexp.MustCompile(`\pL+(\pP+\pL+)*|.`)
	return re.FindAllString(text, -1)
}
