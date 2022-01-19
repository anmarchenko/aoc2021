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

	var board [2000][2000]bool
	maxX := -1
	maxY := -1

	folding := false
	for line := range reader {
		if line == "" {
			folding = true
			continue
		}

		if folding {
			fold := strings.Replace(line, "fold along ", "", 1)
			instructions := strings.Split(fold, "=")
			if len(instructions) != 2 {
				log.Fatal("Wrong fold instructions")
			}
			dimension := instructions[0]
			foldLine, err := strconv.Atoi(instructions[1])
			if err != nil {
				log.Fatal(err)
			}
			if dimension == "x" {
				for i := 1; foldLine+i <= maxX; i++ {
					for j := 0; j <= maxY; j++ {
						board[j][foldLine-i] = board[j][foldLine-i] || board[j][foldLine+i]
					}
				}
				maxX = foldLine - 1
			} else {
				for j := 1; foldLine+j <= maxY; j++ {
					for i := 0; i <= maxX; i++ {
						board[foldLine-j][i] = board[foldLine-j][i] || board[foldLine+j][i]
					}
				}
				maxY = foldLine - 1
			}
		} else {
			coords := strings.Split(line, ",")
			var coordNums [2]int
			for i, coord := range coords {
				n, err := strconv.Atoi(coord)
				if err != nil {
					log.Fatal(err)
				}
				coordNums[i] = n
			}
			x := coordNums[0]
			y := coordNums[1]
			board[y][x] = true
			if y > maxY {
				maxY = y
			}
			if x > maxX {
				maxX = x
			}
		}
	}

	sum := 0
	for i := 0; i <= maxX; i++ {
		for j := 0; j <= maxY; j++ {
			if board[j][i] {
				sum++
			}
		}
	}
	fmt.Println(sum)

	for i := 0; i <= maxY; i++ {
		for j := 0; j <= maxX; j++ {
			var ch string
			if board[i][j] {
				ch = "#"
			} else {
				ch = "."
			}
			fmt.Print(ch)
		}
		fmt.Print("\n")
	}
}
