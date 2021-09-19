package mnemonic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// isUnique checks if the passed in mnemonic is unique
func isUnique(input []string) bool {
	wordMap := make(map[string]bool)

	for _, word := range input {
		if _, found := wordMap[word]; found {
			return false
		}

		wordMap[word] = true
	}

	return true
}

func TestMnemonicGenerator_GenerateMnemonic(t *testing.T) {
	testTable := []struct {
		name     string
		numWords int
	}{
		{
			"Generate 5 different words",
			5,
		},
		{
			"Generate 10 different words",
			10,
		},
		{
			"Generate 100 different words",
			100,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mg := &MnemonicGenerator{NumWords: testCase.numWords}

			output, outputErr := mg.GenerateMnemonic()

			outputArr := strings.Split(output, " ")
			assert.NoError(t, outputErr)
			assert.Len(t, outputArr, testCase.numWords)
			assert.True(t, isUnique(outputArr))
		})
	}
}
