// Provide functions for the phonetic spelling utilities.
package main

import (
	"bufio"
	"log"
	"os"
	"path"
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
