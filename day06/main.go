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
	school := parse(input)
	// fmt.Printf("%d\n", school)
	for i := 0; i < 80; i++ {
		// fmt.Printf("%d %d\n", i, school)
		school.Step()
	}
	fmt.Printf("%d\n", school.Count())
}

func run2(input string) {
	school := parse(input)
	// fmt.Printf("%d\n", school)
	for i := 0; i < 256; i++ {
		// fmt.Printf("%d\n", i)
		school.Step()
	}
	fmt.Printf("%d\n", school.Count())
}

func parse(input string) School {
	lines := strings.Split(input, "\n")
	res := make([]int, 9)
	for _, x := range strings.Split(lines[0], ",") {
		n, _ := strconv.Atoi(x)
		res[n]++
	}
	return School{res}
}

type School struct {
	timers []int
}

func (school *School) Step() {
	born := school.timers[0]
	next := append(school.timers[1:], born)
	next[6] += born
	school.timers = next
}

func (school *School) Count() int {
	sum := 0
	for _, x := range school.timers {
		sum += x
	}
	return sum
}
