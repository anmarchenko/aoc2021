package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

func toNumber(a []string) uint64 {
	i, err := strconv.ParseUint(strings.Join(a, ""), 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func filterNumbers(nums []string, ch byte, pos int) []string {
	res := make([]string, 0)
	for _, num := range nums {
		if num[pos] == ch {
			res = append(res, num)
		}
	}
	return res
}

func constructNumber(frequent bool, pos int, nums []string) uint64 {
	lines_count := len(nums)
	if lines_count == 0 {
		log.Fatal(frequent, pos, nums)
	}
	if lines_count == 1 {
		return toNumber(strings.Split(nums[0], ""))
	}
	sum := 0
	for _, num := range nums {
		if num[pos] == '1' {
			sum += 1
		}
	}
	if frequent {
		char := byte('1')
		if sum < lines_count-sum {
			char = '0'
		}
		return constructNumber(frequent, pos+1, filterNumbers(nums, char, pos))
	} else {
		char := byte('0')
		if sum < lines_count-sum {
			char = '1'
		}
		return constructNumber(frequent, pos+1, filterNumbers(nums, char, pos))
	}
}

func main() {
	reader, err := Readlines("input2")
	if err != nil {
		log.Fatal(err)
	}
	// nums := make([]int, line_length)
	lines := make([]string, 0)
	for line := range reader {
		lines = append(lines, line)
	}
	fmt.Println(constructNumber(false, 0, lines) * constructNumber(true, 0, lines))

	// for line := range reader {
	// 	for i, char := range line {
	// 		if char == '1' {
	// 			nums[i] += 1
	// 		}
	// 	}
	// }
	// gamma := make([]string, line_length)
	// epsilon := make([]string, line_length)
	// for i, num := range nums {
	// 	if num > lines_count/2 {
	// 		gamma[i] = "1"
	// 		epsilon[i] = "0"
	// 	} else {
	// 		gamma[i] = "0"
	// 		epsilon[i] = "1"
	// 	}
	// }
}
