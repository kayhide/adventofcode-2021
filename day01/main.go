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
	count := 0
	last := 0
	for i, line := range lines {
		x, e := strconv.Atoi(line)
		if e == nil {
			if 0 < i {
				if last < x {
					count++
				}
			}
			last = x
		}
	}
	fmt.Printf("%d\n", count)
}

func run2(input string) {
	lines := strings.Split(input, "\n")
	count := 0
	last := 0
	lasts := []int{}
	for _, line := range lines {
		x, e := strconv.Atoi(line)
		if e == nil {
			if len(lasts) == 3 {
				lasts = append(lasts[1:], x)
				y := sum(lasts)
				if last < y {
					count++
				}
				last = y
			} else {
				lasts = append(lasts, x)
				last = sum(lasts)
			}
		}
	}
	fmt.Printf("%d\n", count)
}

func sum(xs []int) int {
	res := 0
	for _, x := range xs {
		res += x
	}
	return res
}
