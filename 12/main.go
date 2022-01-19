package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type Graph struct {
	nodes map[string]*Node
	start *Node
	end   *Node
}

type Node struct {
	neighbours    []*Node
	visited       bool
	big           bool
	name          string
	numVisits     int
	currentVisits int
}

func (g *Graph) maybeAddNode(nodeName string) {
	if g.nodes[nodeName] != nil {
		return
	}

	node := &Node{visited: false, big: false, neighbours: make([]*Node, 0, 20), name: nodeName, numVisits: 0, currentVisits: 0}
	g.nodes[nodeName] = node

	if nodeName == "start" {
		g.start = node
	} else if nodeName == "end" {
		g.end = node
	} else if unicode.IsUpper(rune(nodeName[0])) {
		node.big = true
	}
}

func (g *Graph) addNeighbours(left string, right string) {
	leftNode := g.nodes[left]
	rightNode := g.nodes[right]

	leftNode.neighbours = append(leftNode.neighbours, rightNode)
	rightNode.neighbours = append(rightNode.neighbours, leftNode)
}

func (g *Graph) dfs(node *Node) {
	node.visited = true
	node.numVisits++

	for _, next := range node.neighbours {
		if !next.visited || next.big {
			g.dfs(next)
		}
	}
	node.visited = false
}

func (g *Graph) dfsTwice(node *Node, twice bool) {
	node.currentVisits++
	node.numVisits++

	if node.name == g.end.name {
		node.currentVisits = 0
		return
	}

	for _, next := range node.neighbours {
		if next.name != g.start.name && next.name != g.end.name && !twice && next.currentVisits == 1 && !next.big {
			g.dfsTwice(next, true)
		}
		if next.currentVisits == 0 || next.big {
			g.dfsTwice(next, twice)
		}
	}
	node.currentVisits--
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

	g := &Graph{nodes: make(map[string]*Node)}

	for line := range reader {
		nodeNames := strings.Split(line, "-")
		for _, nodeName := range nodeNames {
			g.maybeAddNode(nodeName)
		}
		g.addNeighbours(nodeNames[0], nodeNames[1])
	}
	g.dfsTwice(g.start, false)
	fmt.Println(g.end.numVisits)
}
