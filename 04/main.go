package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type board struct {
	numbers  [][]int
	selected [][]bool
	Won      bool
}

func (b board) markNumber(number int) {
	for i, row := range b.numbers {
		for j, n := range row {
			if n == number {
				b.selected[i][j] = true
			}
		}
	}
}
func (b *board) checkWinner() bool {
	winner := false
	for i := 0; i < len(b.numbers); i++ {
		row_winner := true
		column_winner := true
		for j := 0; j < len(b.numbers); j++ {
			row_winner = row_winner && b.selected[i][j]
			column_winner = column_winner && b.selected[j][i]
		}
		winner = winner || row_winner || column_winner
	}
	b.Won = winner
	return winner
}

func (b board) winningScore(number int) int {
	sum := 0
	for i := 0; i < len(b.numbers); i++ {
		for j := 0; j < len(b.numbers); j++ {
			if !b.selected[i][j] {
				sum += b.numbers[i][j]
			}
		}
	}
	return sum * number
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

func stringToNumbers(str string, sep string) []int {
	nums := make([]int, 0)
	var strs []string
	if sep == "whitespaces" {
		strs = strings.Fields(str)
	} else {
		strs = strings.Split(str, sep)
	}
	for _, num_str := range strs {
		num, err := strconv.Atoi(num_str)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, num)
	}
	return nums
}

func stringsToNumbers(strs []string) [][]int {
	res := make([][]int, 0)
	for _, str := range strs {
		res = append(res, stringToNumbers(str, "whitespaces"))
	}
	return res
}

func makeBoard(lines []string) board {
	return board{
		numbers: stringsToNumbers(lines),
		selected: [][]bool{
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, false, false, false},
		},
		Won: false,
	}
}

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	input_numbers := make([]int, 0)
	boards := make([]board, 0)

	lines := make([]string, 0)

	for line := range reader {
		if len(input_numbers) == 0 {
			input_numbers = stringToNumbers(line, ",")
			continue
		}
		if len(line) == 0 {
			continue
		}
		lines = append(lines, line)
		if len(lines) == 5 {
			boards = append(boards, makeBoard(lines))
			lines = make([]string, 0)
		}
	}
	for _, number := range input_numbers {
		for _, b := range boards {
			b.markNumber(number)
		}
		for i := 0; i < len(boards); i++ {
			b := &boards[i]
			if !b.Won && b.checkWinner() {
				not_won := 0
				for _, bb := range boards {
					if !bb.Won {
						not_won++
					}
				}
				if not_won == 0 {
					fmt.Println(b.winningScore(number))
					os.Exit(0)
				}
			}
		}
	}
}
