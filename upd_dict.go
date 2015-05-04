// Update the dictionary with the replacment or new word/pron pairs in the other dictionary file.
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
	updFile := flag.String("upd", "", "Name of CMU-formatted upd file")
	flag.Parse()

	if len(*updFile) == 0 {
		log.Fatal("Missing -upd argument")
	}

	// Load replacements file.
	repls := cmudict.LoadDict(*updFile)

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
			word := flds[0]
			pron := flds[1]

			// Add replacement exactly once if it is available.
			if repls[word] != "" {
				pron = repls[word]
				delete(repls, word)
			}
			fmt.Printf("%s  %s\n", word, pron)
		}
	}

	// Now append any new terms.
	for word, pron := range repls {
		fmt.Printf("%s  %s\n", word, pron)
	}
}
