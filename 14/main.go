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

const STEPS = 40

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	var originalLine string
	pairs := make(map[string]int)
	rules := make(map[string]string)

	readingRules := false
	for line := range reader {
		if line == "" {
			readingRules = true
			continue
		}
		if readingRules {
			elements := strings.Split(line, " -> ")
			rules[elements[0]] = elements[1]
		} else {
			originalLine = line
			for i := 0; i < len(line)-1; i++ {
				var pair string
				if i == len(line)-2 {
					pair = line[i:]
				} else {
					pair = line[i : i+2]
				}
				pairs[pair] = pairs[pair] + 1
			}
		}
	}

	for i := 0; i < STEPS; i++ {
		newPairs := make(map[string]int)
		for pair, num := range pairs {
			insertion := rules[pair]
			if insertion == "" {
				log.Fatal(fmt.Sprintf("Unknown pair %v", pair))
			}

			leftPair := fmt.Sprintf("%v%v", pair[0:1], insertion)
			rightPair := fmt.Sprintf("%v%v", insertion, pair[1:])

			newPairs[leftPair] = newPairs[leftPair] + num
			newPairs[rightPair] = newPairs[rightPair] + num
		}
		pairs = newPairs
	}

	counts := make(map[byte]int)
	for pair, num := range pairs {
		counts[pair[0]] += num
	}
	counts[originalLine[len(originalLine)-1]] += 1

	nums := make([]int, 0, len(counts))
	for _, num := range counts {
		nums = append(nums, num)
	}
	sort.Ints(nums)
	fmt.Println(nums[len(nums)-1] - nums[0])
}
