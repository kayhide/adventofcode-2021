package main

import (
	"fmt"
	"io/ioutil"
	"math"
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
	swarm := parse(input)
	min, max := minmax(swarm.crabs)

	cheapest := []int{}
	for i := min; i <= max; i++ {
		cost := sum(swarm.alignAt(i))
		if len(cheapest) == 0 {
			cheapest = []int{cost}
		} else {
			if cost < cheapest[0] {
				cheapest[0] = cost
			}
		}
	}

	fmt.Printf("%d\n", cheapest[0])
}

func run2(input string) {
	swarm := parse(input)
	min, max := minmax(swarm.crabs)

	cheapest := []int{}
	for i := min; i <= max; i++ {
		cost := sum(swarm.alignAt2(i))
		if len(cheapest) == 0 {
			cheapest = []int{cost}
		} else {
			if cost < cheapest[0] {
				cheapest[0] = cost
			}
		}
	}

	fmt.Printf("%d\n", cheapest[0])
}

type Swarm struct {
	crabs []int
}

func parse(input string) Swarm {
	lines := strings.Split(input, "\n")
	res := []int{}
	for _, line := range strings.Split(lines[0], ",") {
		x, e := strconv.Atoi(line)
		if e == nil {
			res = append(res, x)
		}
	}
	return Swarm{res}
}

func (swarm *Swarm) alignAt(pos int) []int {
	costs := make([]int, len([]int(swarm.crabs)))
	for i, p := range swarm.crabs {
		costs[i] = int(math.Abs(float64(p - pos)))
	}
	return costs
}

func (swarm *Swarm) alignAt2(pos int) []int {
	costs := make([]int, len([]int(swarm.crabs)))
	for i, p := range swarm.crabs {
		n := int(math.Abs(float64(p - pos)))
		costs[i] = n * (n + 1) / 2
	}
	return costs
}

func sum(xs []int) int {
	sum := 0
	for _, x := range xs {
		sum += x
	}
	return sum
}

func minmax(xs []int) (int, int) {
	if 0 == len(xs) {
		return 0, 0
	}
	min := xs[0]
	max := xs[0]
	for _, x := range xs[1:] {
		if x < min {
			min = x
		}
		if max < x {
			max = x
		}
	}
	return min, max
}

func ave(xs []int) float64 {
	s := sum(xs)
	return float64(s) / float64(len(xs))
}
