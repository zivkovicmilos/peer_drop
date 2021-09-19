package mnemonic

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type MnemonicGenerator struct {
	NumWords int // Number of words in the mnemonic
	wordList []string
}

// GenerateMnemonic generates an x word mnemonic based on the present wordlist
func (mg *MnemonicGenerator) GenerateMnemonic() (string, error) {
	// Load the wordlist into the MG object
	if loadErr := mg.loadWordlist(); loadErr != nil {
		return "", loadErr
	}

	rand.Seed(time.Now().UnixNano())

	generatedWords := make([]string, mg.NumWords)
	seenWords := make(map[string]bool)
	cnt := 0

	max := len(mg.wordList)

	for cnt < mg.NumWords {
		randIndex := rand.Intn(max)
		generatedWord := mg.wordList[randIndex]

		if _, seen := seenWords[generatedWord]; !seen {
			generatedWords[cnt] = generatedWord
			cnt++
			seenWords[generatedWord] = true
		}
	}

	return strings.Join(generatedWords[:], " "), nil
}

// loadWordlist loads the wordlist from disk
func (mg *MnemonicGenerator) loadWordlist() error {
	mg.wordList = make([]string, 0)

	wordlistFile, err := os.Open("wordList.txt")
	if err != nil {
		return fmt.Errorf("unable to open wordlist file, %v", err)
	}
	defer wordlistFile.Close()

	scanner := bufio.NewScanner(wordlistFile)
	for scanner.Scan() {
		mg.wordList = append(mg.wordList, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("scanner error, %v", err)
	}

	return nil
}
