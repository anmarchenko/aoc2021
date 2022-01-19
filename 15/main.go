package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
)

const DIMENSION = 500
const REPETITION_COUNT = 24

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

type point struct {
	x        int
	y        int
	distance int
	index    int
}

// A PriorityQueue implements heap.Interface and holds points.
type PriorityQueue []*point

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*point)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

var costs [DIMENSION][DIMENSION]byte
var minCosts [DIMENSION][DIMENSION]int
var visited [DIMENSION][DIMENSION]bool

func dijkstra(p point) {
	q := make(PriorityQueue, 0, DIMENSION*DIMENSION)
	heap.Init(&q)

	minCosts[p.x][p.y] = 0
	heap.Push(&q, &p)

	for q.Len() > 0 {
		curPoint := heap.Pop(&q).(*point)
		if visited[curPoint.x][curPoint.y] {
			continue
		}
		visited[curPoint.x][curPoint.y] = true
		curCost := curPoint.distance

		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if (i != 0 && j != 0) || (i == 0 && j == 0) {
					continue
				}
				newX := curPoint.x + i
				newY := curPoint.y + j
				if newX >= 0 && newX < DIMENSION && newY >= 0 && newY < DIMENSION {
					newCost := curCost + int(costs[newX][newY])
					if newCost < minCosts[newX][newY] {
						minCosts[newX][newY] = newCost
						heap.Push(&q, &point{x: newX, y: newY, distance: newCost})
					}
				}
			}
		}
	}
}

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < DIMENSION; i++ {
		for j := 0; j < DIMENSION; j++ {
			minCosts[i][j] = math.MaxInt
		}
	}

	lineIndex := 0
	for line := range reader {
		for i, ch := range line {
			costs[lineIndex][i] = byte(ch) - byte('0')
		}
		lineIndex++
	}

	smallTileDimension := DIMENSION / 5
	for repIndex := 1; repIndex <= REPETITION_COUNT; repIndex++ {
		y := repIndex % 5
		x := repIndex / 5
		for i := 0; i < smallTileDimension; i++ {
			for j := 0; j < smallTileDimension; j++ {
				curX := x*smallTileDimension + i
				curY := y*smallTileDimension + j
				var prevX, prevY int
				if y > 0 {
					prevY = (y-1)*smallTileDimension + j
					prevX = curX
				} else {
					prevY = curY
					prevX = (x-1)*smallTileDimension + i
				}
				costs[curX][curY] = costs[prevX][prevY] + 1
				if costs[curX][curY] == 10 {
					costs[curX][curY] = 1
				}
			}
		}
	}

	dijkstra(point{x: 0, y: 0, distance: 0})
	fmt.Println(minCosts[DIMENSION-1][DIMENSION-1])
}
