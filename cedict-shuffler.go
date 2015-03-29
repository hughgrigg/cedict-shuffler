package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	cedictFile, err := os.Open(os.Args[1])
	check(err)
	entry := makeCedictEntry(randomLine(cedictFile))
	fmt.Printf(
		"%s (%s): %s\n", entry.traditional, entry.prettyPinyin(), entry.definition,
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

// Convert entry's pinyin to tonemarks and apply ANSI colours
func (e *Entry) prettyPinyin() string {
	syllables := strings.Split(e.pinyin, " ")
	marker := toneMarker()
	var buffer bytes.Buffer
	for i, syllable := range syllables {
		tone, letters := toneAndLetters(syllable)
		buffer.WriteString(marker(tone, letters))
		if i < len(syllables)-1 {
			buffer.WriteString(" ")
		}
	}
	return buffer.String()
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
		pinyin:      vsToUmlaut(parts[3]),
		definition:  parts[4],
	}
}

// Replace "v" with "\u00fc" and "V" with "\u00dc"
func vsToUmlaut(pinyin string) string {
	return strings.Replace(
		strings.Replace(pinyin, "V", "\u00dc", -1), "v", "\u00fc", -1,
	)
}

// Get pinyin tonemarking closure
func toneMarker() func(tone int, letters string) string {
	toneMarks := [4]string{"\u0304", "\u0301", "\u030C", "\u0300"}
	targets := [13]string{"A", "E", "I", "O", "U", "\u00dc", "iu", "a", "e", "i",
		"o", "u", "\u00fc"}
	toneMarker := func(tone int, letters string) string {
		checkTone(tone)
		if tone == 5 {
			return letters
		}
		// Replace first found tonemark target vowel with tonemarked version
		for i := 0; i < 13; i++ {
			if strings.Index(letters, targets[i]) > -1 {
				replaced := strings.Replace(letters, targets[i],
					targets[i]+toneMarks[tone-1], 1)
				return string(norm.NFC.Bytes([]byte(replaced)))
			}
		}
		return letters
	}
	return toneMarker
}

// Bail out if tone not in range 1-5
func checkTone(tone int) {
	if tone < 1 || tone > 5 {
		panic("Invalid tone: " + string(tone))
	}
}

// Get the tone and plain letters of a pinyin syllable
func toneAndLetters(syllable string) (int, string) {
	tone, err := strconv.Atoi(syllable[len(syllable)-1 : len(syllable)])
	check(err)
	checkTone(tone)
	letters := syllable[:len(syllable)-1]
	return tone, letters
}
