package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

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

func parsePair(pair string) point {
	nums := strings.Split(strings.TrimSpace(pair), ",")
	x, err := strconv.Atoi(nums[0])
	if err != nil {
		log.Fatal("incorrect number")
	}
	y, err := strconv.Atoi(nums[1])
	if err != nil {
		log.Fatal("incorrect number")
	}
	return point{
		x: x,
		y: y,
	}
}

const MAX_DIMENSION = 1000

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	var board [1000][1000]int

	for line := range reader {
		coord_pairs := strings.Split(line, "->")
		start := parsePair(coord_pairs[0])
		end := parsePair(coord_pairs[1])

		if start.y == end.y {
			for i := min(start.x, end.x); i <= max(start.x, end.x); i++ {
				board[i][start.y]++
			}
		} else if start.x == end.x {
			for j := min(start.y, end.y); j <= max(start.y, end.y); j++ {
				board[start.x][j]++
			}
		} else {
			dirx := (end.x - start.x) / abs(end.x-start.x)
			diry := (end.y - start.y) / abs(end.y-start.y)
			for i := 0; i <= abs(start.x-end.x); i++ {
				board[start.x+i*dirx][start.y+i*diry]++
			}
		}
	}

	res := 0
	for i := 0; i < MAX_DIMENSION; i++ {
		for j := 0; j < MAX_DIMENSION; j++ {
			if board[i][j] > 1 {
				res++
			}
		}
	}
	fmt.Println(res)
}
