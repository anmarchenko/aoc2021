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

const STEPS = 256
const NEW_FISH_TIMER = 8
const NEXT_FISH_TIMER = 6

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	var lanternfishes [9]uint64

	for line := range reader {
		for _, fish := range strings.Split(line, ",") {
			num, err := strconv.Atoi(fish)
			if err != nil {
				log.Fatal(err)
			}
			lanternfishes[num]++
		}
	}

	for day := 0; day < STEPS; day++ {
		new_fishes := lanternfishes[0]
		for i := 0; i < len(lanternfishes)-1; i++ {
			lanternfishes[i] = lanternfishes[i+1]
		}
		lanternfishes[NEXT_FISH_TIMER] += new_fishes
		lanternfishes[NEW_FISH_TIMER] = new_fishes
	}
	res := 0
	for i := 0; i < len(lanternfishes); i++ {
		res += int(lanternfishes[i])
	}
	fmt.Println(res)
}
