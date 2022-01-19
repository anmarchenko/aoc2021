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

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	position := 0
	depth := 0
	aim := 0

	for line := range reader {
		command := strings.Split(line, " ")
		number, err := strconv.Atoi(command[1])
		if err != nil {
			log.Fatal(err)
		}

		switch command[0] {
		case "forward":
			position += number
			depth += aim * number
		case "up":
			aim -= number
		case "down":
			aim += number
		}
	}

	fmt.Println(position * depth)
}
