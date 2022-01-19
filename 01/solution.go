package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
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

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	currentSum := math.MaxInt
	answer := 0
	nums := []int{}

	for line := range reader {
		number, err := strconv.Atoi(line)

		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, number)
	}

	for index, _ := range nums {
		if index+2 == len(nums) {
			break
		}
		sum := nums[index] + nums[index+1] + nums[index+2]
		if sum > currentSum {
			answer += 1
		}
		currentSum = sum
	}

	fmt.Printf("%v\n", answer)
}
