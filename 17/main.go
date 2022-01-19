package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type probe struct {
	x    int
	y    int
	velX int
	velY int
}

func (p *probe) tick() {
	p.x += p.velX
	p.y += p.velY

	if p.velX > 0 {
		p.velX--
	} else if p.velX < 0 {
		p.velX++
	}
	p.velY--
}

func (p *probe) missed(target *area) bool {
	return p.x > target.xMax || p.y < target.yMin
}

func (p *probe) inside(target *area) bool {
	return p.x >= target.xMin && p.x <= target.xMax && p.y >= target.yMin && p.y <= target.yMax
}

func (p *probe) launch(target *area) (bool, int) {
	hit := false
	maxY := 0

	for {
		p.tick()

		if p.y > maxY {
			maxY = p.y
		}
		if p.inside(target) {
			hit = true
			break
		}
		if p.missed(target) {
			break
		}
	}

	return hit, maxY
}

type area struct {
	xMin int
	xMax int
	yMin int
	yMax int
}

func parseInt(num string) int {
	n, err := strconv.Atoi(num)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func parseCoords(input string) (int, int) {
	numStrings := strings.Split(input, "..")
	return parseInt(numStrings[0]), parseInt(numStrings[1])
}

func parseArea(desc string) *area {
	res := area{}

	areaDescriptions := strings.Split(desc, ", ")
	res.xMin, res.xMax = parseCoords(areaDescriptions[0])
	res.yMin, res.yMax = parseCoords(areaDescriptions[1])
	return &res
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

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	var input string
	for line := range reader {
		input = line
	}

	input = strings.Replace(
		strings.Replace(
			strings.Replace(input, "target area: ", "", 1),
			"x=",
			"",
			1,
		),
		"y=",
		"",
		1,
	)

	target := parseArea(input)

	res := 0
	// now we try several different probes
	for x := 1; x < 300; x++ {
		for y := -300; y < 300; y++ {
			p := probe{x: 0, y: 0, velX: x, velY: y}
			hit, _ := p.launch(target)
			if hit {
				res++
			}
		}
	}
	fmt.Println(res)
}
