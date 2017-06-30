package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func populateCodes() {
	ids, err := getInvitationsWithoutCodes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	codes, err := generateCodes(len(ids))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	for i, id := range ids {
		if err := updateInvitationCode(id, codes[i]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	}
}

func generateCodes(n int) ([]string, error) {
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	words := []string{}
	for scanner.Scan() {
		word := strings.ToUpper(scanner.Text())

		// Only want 4 character words
		if len(word) != 4 {
			continue
		}

		// Don't want words with ambiguous characters
		// Inspired by Crockford's Base32
		if strings.ContainsRune(word, 'I') ||
			strings.ContainsRune(word, 'L') ||
			strings.ContainsRune(word, 'O') ||
			strings.ContainsRune(word, 'U') {
			continue
		}

		words = append(words, word)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	codes := make([]string, n)
	for i := 0; i < n; i++ {
		word := words[rand.Intn(len(words))]
		number := rand.Intn(99)
		codes[i] = fmt.Sprintf("%s%02d", word, number)
	}

	return codes, nil
}
