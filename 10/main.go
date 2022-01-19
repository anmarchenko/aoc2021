package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func Readlines(path string) (<-chan string, error) {
	fobj, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(fobj)
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	chnl := make(chan string)
	go func() {
		for scanner.Scan() {
			chnl <- scanner.Text()
		}
		close(chnl)
	}()

	return chnl, nil
}

var errorCosts = map[string]int{
	")": 3,
	"]": 57,
	"}": 1197,
	">": 25137,
}

var openings = "([<{"
var correctClosings = map[string]string{
	")": "(",
	"]": "[",
	"}": "{",
	">": "<",
}

var openingToClosings = map[string]string{
	"(": ")",
	"[": "]",
	"{": "}",
	"<": ">",
}
var autocompleteCosts = map[string]int{
	")": 1,
	"]": 2,
	"}": 3,
	">": 4,
}

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	errorScore := 0

	aScores := make([]int, 0)

	for line := range reader {
		stack := make([]string, 0)
		error := false
		for _, ch := range strings.Split(line, "") {
			if strings.Contains(openings, ch) {
				stack = append(stack, ch)
			} else {
				if correctClosings[ch] == stack[len(stack)-1] {
					stack = stack[:len(stack)-1]
				} else {
					errorScore += errorCosts[ch]
					error = true
					break
				}
			}
		}

		if !error && len(stack) > 0 {
			autocompleteScore := 0
			for i := len(stack) - 1; i >= 0; i-- {
				autocompleteScore = autocompleteScore*5 + autocompleteCosts[openingToClosings[stack[i]]]
			}
			aScores = append(aScores, autocompleteScore)
		}
	}
	fmt.Println(errorScore)
	sort.Ints(aScores)
	fmt.Println(aScores[len(aScores)/2])
}
