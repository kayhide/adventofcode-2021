package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input, _ := ioutil.ReadFile("input.txt")
	// input, _ := ioutil.ReadFile("input-1.txt")
	run1(string(input))
	run2(string(input))
}

func run1(input string) {
	lines := strings.Split(input, "\n")
	depth := 0
	pos := 0
	for _, line := range lines {
		xs := strings.Split(line, " ")
		if len(xs) == 2 {
			v, e := strconv.Atoi(xs[1])
			if e == nil {
				// fmt.Printf("%s\n", line)
				switch xs[0] {
				case "up":
					depth -= v
				case "down":
					depth += v
				case "forward":
					pos += v
				}
			}

		}
	}
	fmt.Printf("%d\n", depth*pos)
}

func run2(input string) {
	lines := strings.Split(input, "\n")
	depth := 0
	pos := 0
	aim := 0
	for _, line := range lines {
		xs := strings.Split(line, " ")
		if len(xs) == 2 {
			v, e := strconv.Atoi(xs[1])
			if e == nil {
				// fmt.Printf("%s\n", line)
				switch xs[0] {
				case "up":
					aim -= v
				case "down":
					aim += v
				case "forward":
					pos += v
					depth += aim * v
				}
			}

		}
	}
	fmt.Printf("%d\n", depth*pos)
}
