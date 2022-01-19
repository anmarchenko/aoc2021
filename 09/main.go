package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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

func borderRow(a []int) {
	for i := 0; i < len(a); i++ {
		a[i] = math.MaxInt
	}
}

type point struct {
	x int
	y int
}

func basinSize(p point, board [][]int) int {
	res := 1
	board[p.x][p.y] = math.MaxInt
	queue := []point{p}

	for len(queue) != 0 {
		head := queue[0]
		queue = queue[1:]

		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if i != 0 && j != 0 {
					continue
				}
				curX := head.x + i
				curY := head.y + j
				if board[curX][curY] < 9 {
					queue = append(queue, point{x: curX, y: curY})
					board[curX][curY] = math.MaxInt
					res++
				}
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

	l := -1
	board := make([][]int, 0)

	for line := range reader {
		l = len(line)
		if len(board) == 0 {
			board = append(board, make([]int, l+2))
			borderRow(board[0])
		}
		row := make([]int, l+2)
		row[0] = math.MaxInt
		row[l+1] = math.MaxInt
		for i, ch := range strings.Split(line, "") {
			n, err := strconv.Atoi(ch)
			if err != nil {
				log.Fatal(err)
			}
			row[i+1] = n
		}
		board = append(board, row)
	}
	board = append(board, make([]int, l+2))
	borderRow(board[len(board)-1])

	low_points := make([]point, 0)
	for i := 1; i < len(board)-1; i++ {
		for j := 1; j <= l; j++ {
			if board[i][j] < board[i-1][j] && board[i][j] < board[i+1][j] && board[i][j] < board[i][j-1] && board[i][j] < board[i][j+1] {
				low_points = append(low_points, point{x: i, y: j})
			}
		}
	}

	basins := make([]int, 0)
	for _, p := range low_points {
		basins = append(basins, basinSize(p, board))
	}

	sort.Ints(basins)
	basinsCount := len(basins)
	if basinsCount < 3 {
		log.Fatal("less than 3 basins found")
	}
	fmt.Println(basins[basinsCount-1] * basins[basinsCount-2] * basins[basinsCount-3])
}
