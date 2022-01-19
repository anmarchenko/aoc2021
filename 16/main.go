package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type packet struct {
	version    int
	packetType int
	value      int
	subpackets []*packet
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

const PACKET_TYPE_SUM = 0
const PACKET_TYPE_PRODUCT = 1
const PACKET_TYPE_MINIMUM = 2
const PACKET_TYPE_MAXIMUM = 3
const PACKET_TYPE_LITERAL = 4
const PACKET_TYPE_GREATER_THAN = 5
const PACKET_TYPE_LESS_THAN = 6
const PACKET_TYPE_EQUAL_TO = 7

const BITS_PER_VALUE = 5

const MODE_VERSION = 0
const MODE_TYPE = 1
const MODE_VALUE = 2
const MODE_LENGTH_BIT = 3
const MODE_READ_SUBPACKETS_COUNT = 4
const MODE_READ_SUBPACKETS_LENGTH = 5
const MODE_SUBPACKETS_BY_COUNT = 6
const MODE_SUBPACKETS_BY_LENGTH = 7

func parseMessage(message []byte) (*packet, int) {
	p := &packet{version: 0, value: 0, packetType: 0, subpackets: make([]*packet, 0)}

	read := 0
	readingMode := MODE_VERSION

	valueLast := false
	valueRead := 0

	subCountLeft := 11
	subCount := 0

	subLengthLeft := 15
	subLength := 0

	for read < len(message) {
		bit := int(message[read])
		switch readingMode {
		case MODE_VERSION:
			p.version = p.version*2 + bit
			read++
			if read == 3 {
				readingMode = MODE_TYPE
			}
		case MODE_TYPE:
			p.packetType = p.packetType*2 + bit
			read++
			if read == 6 {
				if p.packetType == PACKET_TYPE_LITERAL {
					readingMode = MODE_VALUE
				} else {
					readingMode = MODE_LENGTH_BIT
				}
			}
		case MODE_VALUE:
			if valueRead%BITS_PER_VALUE == 0 {
				if valueLast {
					return p, read
				}
				// check if next is last
				if bit == 0 {
					valueLast = true
				}
			} else {
				p.value = p.value*2 + bit
			}
			read++
			valueRead++
		case MODE_LENGTH_BIT:
			if bit == 0 {
				readingMode = MODE_READ_SUBPACKETS_LENGTH
			} else {
				readingMode = MODE_READ_SUBPACKETS_COUNT
			}
			read++
		case MODE_READ_SUBPACKETS_LENGTH:
			subLength = subLength*2 + bit
			read++
			subLengthLeft--
			if subLengthLeft == 0 {
				readingMode = MODE_SUBPACKETS_BY_LENGTH
			}
		case MODE_READ_SUBPACKETS_COUNT:
			subCount = subCount*2 + bit
			read++
			subCountLeft--
			if subCountLeft == 0 {
				readingMode = MODE_SUBPACKETS_BY_COUNT
			}
		case MODE_SUBPACKETS_BY_LENGTH:
			subPacket, readSubpacket := parseMessage(message[read:])
			p.subpackets = append(p.subpackets, subPacket)

			read += readSubpacket
			subLength -= readSubpacket

			if subLength == 0 {
				return p, read
			}
		case MODE_SUBPACKETS_BY_COUNT:
			subPacket, readSubpacket := parseMessage(message[read:])
			p.subpackets = append(p.subpackets, subPacket)

			read += readSubpacket
			subCount--

			if subCount == 0 {
				return p, read
			}
		}
	}
	return p, read
}

func printPackets(p *packet) {
	fmt.Printf("value = [%v]; version = [%v]; type = [%v]\n", p.value, p.version, p.packetType)
	for _, subPacket := range p.subpackets {
		printPackets(subPacket)
	}
}

func sumVersions(p *packet) int {
	sum := p.version
	for _, subPacket := range p.subpackets {
		sum += sumVersions(subPacket)
	}
	return sum
}

func calculate(p *packet) int {
	switch p.packetType {
	case PACKET_TYPE_LITERAL:
		return p.value
	case PACKET_TYPE_SUM:
		res := 0
		for _, sub := range p.subpackets {
			res += calculate(sub)
		}
		return res
	case PACKET_TYPE_PRODUCT:
		res := 1
		for _, sub := range p.subpackets {
			res *= calculate(sub)
		}
		return res
	case PACKET_TYPE_MAXIMUM:
		res := -1
		for _, sub := range p.subpackets {
			val := calculate(sub)
			if val > res {
				res = val
			}
		}
		return res
	case PACKET_TYPE_MINIMUM:
		res := math.MaxInt
		for _, sub := range p.subpackets {
			val := calculate(sub)
			if val < res {
				res = val
			}
		}
		return res
	case PACKET_TYPE_EQUAL_TO:
		left := calculate(p.subpackets[0])
		right := calculate(p.subpackets[1])
		if left == right {
			return 1
		} else {
			return 0
		}
	case PACKET_TYPE_GREATER_THAN:
		left := calculate(p.subpackets[0])
		right := calculate(p.subpackets[1])
		if left > right {
			return 1
		} else {
			return 0
		}
	case PACKET_TYPE_LESS_THAN:
		left := calculate(p.subpackets[0])
		right := calculate(p.subpackets[1])
		if left < right {
			return 1
		} else {
			return 0
		}
	}
	return -1
}

const BITS_PER_LETTER = 4

func main() {
	reader, err := Readlines("input")
	if err != nil {
		log.Fatal(err)
	}

	var hexMessage string
	for line := range reader {
		hexMessage = line
	}
	message := make([]byte, len(hexMessage)*BITS_PER_LETTER)
	for i, ch := range strings.Split(hexMessage, "") {
		n, err := strconv.ParseUint(ch, 16, 32)
		if err != nil {
			log.Fatal(err)
		}
		for j := BITS_PER_LETTER - 1; j >= 0; j-- {
			message[i*BITS_PER_LETTER+j] = byte(n & 0x1)
			n = n >> 1
		}
	}
	res, _ := parseMessage(message)
	printPackets(res)
	fmt.Println(sumVersions(res))
	fmt.Println(calculate(res))
}
