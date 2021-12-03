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
	lines := strings.Split(input, "\n")
	n := len(lines[0])
	xs, ys := countBits(lines)
	zs := make([]int, n)
	for i := range zs {
		if ys[i] < xs[i] {
			zs[i] = 1
		}
	}
	gamma := fromBin(zs)
	epsilon := fromBin(negate(zs))
	fmt.Printf("%d\n", gamma*epsilon)
}

func run2(input string) {
	lines := strings.Split(input, "\n")
	x := oxgen_generator(lines)
	y := co2_scrubber(lines)
	fmt.Printf("%d\n", x*y)
}

func oxgen_generator(lines []string) int {
	for i := range lines[0] {
		xs, ys := countBits(lines)
		// fmt.Println(xs)
		// fmt.Println(ys)
		var n int
		var b byte
		if ys[i] <= xs[i] {
			n = xs[i]
			b = '1'
		} else {
			n = ys[i]
			b = '0'
		}

		lines = filterLines(lines, i, b, n)
		// fmt.Println(i, b, n)
		// fmt.Println(lines)
		if len(lines) == 1 {
			break
		}
	}
	xs, _ := countBits(lines)
	return fromBin(xs)
}

func co2_scrubber(lines []string) int {
	for i := range lines[0] {
		xs, ys := countBits(lines)
		// fmt.Println(xs)
		// fmt.Println(ys)
		var n int
		var b byte
		if ys[i] <= xs[i] {
			n = ys[i]
			b = '0'
		} else {
			n = xs[i]
			b = '1'
		}

		lines = filterLines(lines, i, b, n)
		// fmt.Println(i, b, n)
		// fmt.Println(lines)
		if len(lines) == 1 {
			break
		}
	}
	xs, _ := countBits(lines)
	return fromBin(xs)
}

func filterLines(lines []string, pos int, bit byte, count int) []string {
	res := []string{}
	for _, line := range lines {
		if line[pos] == bit {
			res = append(res, line)
		}
		if len(res) == count {
			return res
		}
	}
	return res
}

func countBits(lines []string) ([]int, []int) {
	n := len(lines[0])
	xs := make([]int, n)
	ys := make([]int, n)
	for _, line := range lines {
		for i, x := range line {
			if x == 49 {
				xs[i]++
			} else {
				ys[i]++
			}
		}
	}
	return xs, ys
}

func fromBin(xs []int) int {
	res := 0
	for _, x := range xs {
		res = res*2 + x
	}
	return res
}

func negate(xs []int) []int {
	res := make([]int, len(xs))
	for i, x := range xs {
		res[i] = -1 * (x - 1)
	}
	return res
}
