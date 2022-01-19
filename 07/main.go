package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sum_arithmetic_progression(n int) int {
	return (n * (1 + n)) / 2
}

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	positions := make([]int, 0)
	max_position := -1

	for line := range reader {
		for _, pos := range strings.Split(line, ",") {
			num, err := strconv.Atoi(pos)
			if err != nil {
				log.Fatal(err)
			}
			positions = append(positions, num)
			if num > max_position {
				max_position = num
			}
		}
	}

	min_fuel := math.MaxInt

	for pos := 0; pos <= max_position; pos++ {
		fuel := 0
		for _, num := range positions {
			fuel += sum_arithmetic_progression(abs(num - pos))
		}
		if fuel < min_fuel {
			min_fuel = fuel
		}
	}
	fmt.Println(min_fuel)
}
