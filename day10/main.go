package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

func main() {
	input, _ := ioutil.ReadFile("input.txt")
	// input, _ := ioutil.ReadFile("input-1.txt")
	run1(string(input))
	run2(string(input))
}

func run1(input string) {
	chunk := parse(input)
	sum := 0
	for _, line := range chunk.lines {
		_, rest := reduce(line)
		if 0 < len(rest) {
			switch rest[0] {
			case ')':
				sum += 3
			case ']':
				sum += 57
			case '}':
				sum += 1197
			case '>':
				sum += 25137
			}
		}
	}

	fmt.Printf("%d\n", sum)
}

func run2(input string) {
	chunk := parse(input)
	scores := []int{}
	for _, line := range chunk.lines {
		buf, rest := reduce(line)
		if 0 == len(rest) {
			score := 0
			for i := len(buf) - 1; 0 <= i; i-- {
				switch buf[i] {
				case '(':
					score = score*5 + 1
				case '[':
					score = score*5 + 2
				case '{':
					score = score*5 + 3
				case '<':
					score = score*5 + 4
				}
			}
			// fmt.Println(buf)
			// fmt.Println(score)
			scores = append(scores, score)
		}
	}

	sort.Ints(scores)
	fmt.Printf("%d\n", scores[len(scores)/2])
}

type Chunk struct {
	lines []string
}

func parse(input string) Chunk {
	lines := strings.Split(input, "\n")
	res := []string{}
	for _, line := range lines {
		if 0 < len(line) {
			res = append(res, line)
		}
	}
	return Chunk{res}
}

func reduce(line string) (string, string) {
	buf := ""
	for i, c := range line {
		buf_, b := consume(buf, c)
		if !b {
			return buf_, line[i:]
		}
		buf = buf_
	}
	return buf, ""
}

func consume(buf string, c rune) (string, bool) {
	switch c {
	case '(':
		return buf + string(c), true
	case '[':
		return buf + string(c), true
	case '{':
		return buf + string(c), true
	case '<':
		return buf + string(c), true
	}
	if len(buf) == 0 {
		return buf, false
	}
	last := buf[len(buf)-1]
	switch {
	case last == '(' && c == ')':
		return buf[:len(buf)-1], true
	case last == '[' && c == ']':
		return buf[:len(buf)-1], true
	case last == '{' && c == '}':
		return buf[:len(buf)-1], true
	case last == '<' && c == '>':
		return buf[:len(buf)-1], true
	}
	return buf, false
}
