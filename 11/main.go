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

const STEPS = 100

func flash(octopuses *[10][10]int, x int, y int) {
	octopuses[x][y] = -1
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			nX := x + i
			nY := y + j
			if nX < 0 || nX > 9 || nY < 0 || nY > 9 {
				continue
			}
			if octopuses[nX][nY] == -1 {
				continue
			}
			octopuses[nX][nY]++
			if octopuses[nX][nY] > 9 {
				flash(octopuses, nX, nY)
			}
		}
	}
}

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	var octopuses [10][10]int

	row_num := 0
	for line := range reader {
		for i, ch := range strings.Split(line, "") {
			n, err := strconv.Atoi(ch)
			if err != nil {
				log.Fatal(err)
			}
			octopuses[row_num][i] = n
		}
		row_num++
	}

	finish := false
	step := 0
	for !finish {
		step++
		// inrement all by 1
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				octopuses[x][y]++
			}
		}
		// flashes
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				if octopuses[x][y] > 9 {
					flash(&octopuses, x, y)
				}
			}
		}

		// count flashes
		flashes := 0
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				if octopuses[x][y] == -1 {
					flashes++
					octopuses[x][y] = 0
				}
			}
		}
		if flashes == 100 {
			finish = true
		}
	}
	fmt.Println(step)
}
