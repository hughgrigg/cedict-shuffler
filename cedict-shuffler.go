package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	cedictFile, err := os.Open(os.Args[1])
	check(err)
	entry := makeCedictEntry(randomLine(cedictFile))
	fmt.Printf(
		"%s (%s): %s\n", entry.traditional, entry.pinyin, entry.definition,
	)
}

// Bail out on an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Get a random line from a file with even distribution
func randomLine(file *os.File) (chosenLine string) {
	scanner := bufio.NewScanner(file)
	count := float64(0)
	rand.Seed(time.Now().UnixNano())
	for scanner.Scan() {
		count++
		if rand.Float64() <= (1 / count) {
			chosenLine = scanner.Text()
		}
	}
	return
}

// A CEDICT entry
type Entry struct {
	simplified  string
	traditional string
	pinyin      string
	definition  string
}

// Make an Entry out of a CEDICT line string
func makeCedictEntry(entry string) Entry {
	entryPattern := regexp.MustCompile(`^(\p{Han}+)\s(\p{Han}+)\s\[(.+)\]\s\/(.+)\/$`)
	parts := entryPattern.FindStringSubmatch(entry)
	if parts == nil || len(parts) != 5 {
		panic("Failed to parse CEDICT line: " + entry)
	}
	return Entry{
		simplified:  parts[1],
		traditional: parts[2],
		pinyin:      parts[3],
		definition:  parts[4],
	}
}
