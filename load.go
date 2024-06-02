package main

import (
	"os"
	"strconv"
	"strings"
)

func getCardIDs(filename string) ([]int, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cardIDs []int
	lines := strings.Split(string(fileBytes), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line) // Trim whitespace
		if line == "" {
			continue // Skip empty lines
		}
		cardID, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		cardIDs = append(cardIDs, cardID)
	}
	return cardIDs, nil
}
