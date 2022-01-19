package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func orderAllLeaves(num *snailnumber) []*snailnumber {
	if num.leaf {
		return []*snailnumber{num}
	} else {
		return append(orderAllLeaves(num.left), orderAllLeaves(num.right)...)
	}
}

func findPairByDepth(node *snailnumber, depth int) *snailnumber {
	if depth > 4 && !node.leaf {
		return node
	}
	if node.leaf {
		return nil
	}

	leftResult := findPairByDepth(node.left, depth+1)
	if leftResult != nil {
		return leftResult
	}

	return findPairByDepth(node.right, depth+1)
}

type snailnumber struct {
	leaf   bool
	value  int
	parent *snailnumber
	left   *snailnumber
	right  *snailnumber
}

func parseSnailNumber(s string) *snailnumber {
	if s[0] == '[' {
		res := &snailnumber{
			leaf:  false,
			value: 0,
		}
		s = s[1 : len(s)-1]

		splitIndex := -1
		depth := 0
		for i, ch := range s {
			if depth == 0 && ch == ',' {
				splitIndex = i
				break
			}
			if ch == '[' {
				depth++
			}
			if ch == ']' {
				depth--
			}
		}

		if splitIndex == -1 {
			log.Fatal(fmt.Sprintf("Could not split string %v", s))
		}

		left := parseSnailNumber(s[:splitIndex])
		right := parseSnailNumber(s[splitIndex+1:])

		left.parent = res
		right.parent = res
		res.left = left
		res.right = right

		return res
	} else {
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}

		return &snailnumber{
			leaf:  true,
			value: n,
		}
	}
}

func (l *snailnumber) add(r *snailnumber) *snailnumber {
	res := &snailnumber{leaf: false, left: l, right: r}
	l.parent = res
	r.parent = res
	res.reduce()
	return res
}

func (num *snailnumber) explodable() *snailnumber {
	return findPairByDepth(num, 1)
}

func (num *snailnumber) splittable() *snailnumber {
	if num.leaf {
		if num.value > 9 {
			return num
		} else {
			return nil
		}
	} else {
		leftSplittable := num.left.splittable()
		if leftSplittable != nil {
			return leftSplittable
		}

		return num.right.splittable()
	}
}

func (num *snailnumber) explode() {
	root := num
	for root.parent != nil {
		root = root.parent
	}
	leaves := orderAllLeaves(root)

	leftIndex := -1
	rightIndex := -1
	for i, node := range leaves {
		if node == num.left {
			leftIndex = i
		}

		if node == num.right {
			rightIndex = i
		}
	}

	if leftIndex == -1 || rightIndex == -1 {
		log.Fatal("oops")
	}

	var leftLeaf *snailnumber
	if leftIndex > 0 {
		leftLeaf = leaves[leftIndex-1]
		leftLeaf.value += num.left.value
	}

	var rightLeaf *snailnumber
	if rightIndex < len(leaves)-1 {
		rightLeaf = leaves[rightIndex+1]
		rightLeaf.value += num.right.value
	}

	num.leaf = true
	num.value = 0
	num.left = nil
	num.right = nil
}

func (num *snailnumber) split() {
	val := num.value

	num.leaf = false
	num.value = 0

	left := &snailnumber{
		leaf:   true,
		parent: num,
		value:  val / 2,
	}
	num.left = left

	right := &snailnumber{
		leaf:   true,
		parent: num,
		value:  (val / 2) + (val % 2),
	}
	num.right = right
}

func (num *snailnumber) reduce() {
	for {
		expl := num.explodable()
		if expl != nil {
			expl.explode()
			continue
		}

		spl := num.splittable()
		if spl != nil {
			spl.split()
			continue
		}

		break
	}
}

func (num *snailnumber) String() string {
	if num.leaf {
		return fmt.Sprintf("%v", num.value)
	}
	return fmt.Sprintf("[%v, %v]", num.left, num.right)
}

func (num *snailnumber) magnitude() int {
	if num.leaf {
		return num.value
	}

	return num.left.magnitude()*3 + num.right.magnitude()*2
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

func maximumSum(l string, r string) int {
	m1 := parseSnailNumber(l).add(parseSnailNumber(r)).magnitude()
	m2 := parseSnailNumber(r).add(parseSnailNumber(l)).magnitude()
	if m1 > m2 {
		return m1
	}
	return m2
}

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	// var sum *snailnumber
	// for line := range reader {
	// 	number := parseSnailNumber(line)
	// 	if sum == nil {
	// 		sum = number
	// 	} else {
	// 		sum = sum.add(number)
	// 	}
	// }
	// fmt.Println(sum)
	// fmt.Println(sum.magnitude())

	numbers := make([]string, 0, 100)
	for line := range reader {
		numbers = append(numbers, line)
	}

	res := math.MinInt
	for i := 0; i < len(numbers); i++ {
		for j := i + 1; j < len(numbers); j++ {
			m := maximumSum(numbers[i], numbers[j])
			if m > res {
				res = m
			}
		}
	}

	fmt.Println(res)
}
