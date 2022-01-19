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

func sortedString(x string) string {
	strs := strings.Split(x, "")
	sort.Strings(strs)
	return strings.Join(strs, "")
}

func similarities(a string, b string) int {
	res := 0
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				res++
				break
			}
		}
	}
	return res
}

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	sum := 0

	for line := range reader {
		stringToNum := make(map[string]int)
		numToString := make(map[int]string)

		io := strings.Split(line, "|")

		input := strings.TrimSpace(io[0])
		input_digits := strings.Split(input, " ")

		next_round := make([]string, 0)
		for _, input_digit := range input_digits {
			digit := sortedString(input_digit)
			if len(digit) == 2 {
				stringToNum[digit] = 1
				numToString[1] = digit
			} else if len(digit) == 3 {
				stringToNum[digit] = 7
				numToString[7] = digit
			} else if len(digit) == 4 {
				stringToNum[digit] = 4
				numToString[4] = digit
			} else if len(digit) == 7 {
				stringToNum[digit] = 8
				numToString[8] = digit
			} else {
				next_round = append(next_round, digit)
			}
		}

		for _, digit := range next_round {
			if len(digit) == 5 {
				if similarities(numToString[1], digit) == 2 {
					stringToNum[digit] = 3
					numToString[3] = digit
				} else if similarities(numToString[4], digit) == 2 {
					stringToNum[digit] = 2
					numToString[2] = digit
				} else {
					stringToNum[digit] = 5
					numToString[5] = digit
				}
			} else if len(digit) == 6 {
				if similarities(numToString[1], digit) == 1 {
					stringToNum[digit] = 6
					numToString[6] = digit
				} else if similarities(numToString[4], digit) == 4 {
					stringToNum[digit] = 9
					numToString[9] = digit
				} else {
					stringToNum[digit] = 0
					numToString[0] = digit
				}
			}
		}

		output := strings.TrimSpace(io[1])
		digits := strings.Split(output, " ")

		number := 0
		pos := 1000

		for _, digit := range digits {
			number += pos * stringToNum[sortedString(digit)]
			pos /= 10
		}
		sum += number
	}
	fmt.Println(sum)
}
