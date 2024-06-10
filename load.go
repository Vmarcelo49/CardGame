package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getCardIDs(filename string) ([]int, error) {
	var cardIDs []int
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line) // Trim whitespace
		if line == "" {
			continue // Skip empty lines
		}
		cardID, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Failed to convert '%s' to integer: %v\n", line, err)
			continue // Skip this line and continue with the next
		}
		cardIDs = append(cardIDs, cardID)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return cardIDs, nil
}
