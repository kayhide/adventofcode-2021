package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	input, _ := ioutil.ReadFile("input.txt")
	// input, _ := ioutil.ReadFile("input-1.txt")
	run1(string(input))
	run2(string(input))
}

func run1(input string) {
	us, d := parseInput(input)
	n := len(us)
	// fmt.Println(us, n, d)

	xs := countBits(us, d)
	gamma := 0
	epsilon := 0
	for i, x := range xs {
		if n-x <= x {
			gamma += 1 << i
		}
		if x <= n-x {
			epsilon += 1 << i
		}
	}
	fmt.Printf("%d\n", gamma*epsilon)
}

func run2(input string) {
	us, d := parseInput(input)
	x := oxgen_generator(us, d)
	y := co2_scrubber(us, d)
	fmt.Printf("%d\n", x*y)
}

func parseInput(input string) ([]uint, int) {
	res := []uint{}
	n := 0
	for _, line := range strings.Split(input, "\n") {
		if 0 < len(line) {
			n = len(line)
			res = append(res, parseBits(line))
		}
	}
	return res, n
}

func oxgen_generator(us []uint, digits int) int {
	for i := digits - 1; 1 < len(us); i-- {
		count := len(us)
		xs := countBits(us, digits)
		// fmt.Println(us, xs)
		if count-xs[i] <= xs[i] {
			us = filterBits(us, i, 1, xs[i])
		} else {
			us = filterBits(us, i, 0, count-xs[i])
		}
		// fmt.Println(i, b, n, us)
	}
	// fmt.Println(int(us[0]))
	return int(us[0])
}

func co2_scrubber(us []uint, digits int) int {
	for i := digits - 1; 1 < len(us); i-- {
		count := len(us)
		xs := countBits(us, digits)
		// fmt.Println(us, xs)
		if count-xs[i] <= xs[i] {
			us = filterBits(us, i, 0, count-xs[i])
		} else {
			us = filterBits(us, i, 1, xs[i])
		}
		// fmt.Println(i, b, n, us)
	}
	// fmt.Println(int(us[0]))
	return int(us[0])
}

func filterBits(us []uint, pos int, bit uint, count int) []uint {
	res := []uint{}
	for _, u := range us {
		if (u>>pos)&1 == bit {
			res = append(res, u)
		}
		if len(res) == count {
			return res
		}
	}
	return res
}

func parseBits(s string) uint {
	var res uint = 0
	for _, x := range s {
		res = res << 1
		if x == '1' {
			res++
		}
	}
	return res
}

func countBits(us []uint, digits int) []int {
	xs := make([]int, digits)
	for _, u := range us {
		for i := range xs {
			if u&1 == 1 {
				xs[i]++
			}
			u = u >> 1
		}
	}
	return xs
}
