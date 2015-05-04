// Limit the dictionary to the given word list.
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
	wordFile := flag.String("word", "", "Name of word list file")
	flag.Parse()

	if len(*wordFile) == 0 {
		log.Fatal("Missing -word argument")
	}

	// Load word file.
	words := loadWordList(*wordFile)

	// Open dict file.
	hdl, err := os.Open(*dictFile)
	if err != nil {
		log.Fatal(err)
	}
	defer hdl.Close()

	// Scan file line by line.
	scanner := bufio.NewScanner(hdl)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		if line[:3] == ";;;" {
			// Print all comment lines
			fmt.Println(line)
		} else {
			flds := strings.Split(line, "  ")
			word, pron := flds[0], flds[1]
			bare := cmudict.StripIndex(word)

			// Add each version of that word that is available.
			if words[bare] {
				fmt.Printf("%s  %s\n", word, pron)
			}
		}
	}
}
