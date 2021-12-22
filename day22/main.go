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
	// input, _ := ioutil.ReadFile("input-2.txt")
	run1(string(input))
	run2(string(input))
}

func run1(input string) {
	steps := parse(input)
	f := Field{}
	f.init(Range{-50, 50})
	for _, step := range steps {
		// fmt.Println(step)
		f.apply(step)
	}

	sum := 0
	for _, b := range f.vs {
		if b {
			// fmt.Println(decode(i, 101))
			sum++
		}
	}
	fmt.Println(sum)
}

func run2(input string) {
	steps := parse(input)
	state := []Step{steps[0]}
	for _, step := range steps[1:] {
		// fmt.Println(i+1, len(steps), len(state), step)
		next := []Step{}
		if step.on {
			next = append(next, step)
		}
		for _, t := range state {
			c := step.cuboid.intersect(t.cuboid)
			// fmt.Println(" ", c.volume())
			if 0 < c.volume() {
				if t.on {
					next = append(next, Step{false, c})
				} else {
					next = append(next, Step{true, c})
				}
			}
		}
		// state = append(state, next...)
		state = merge(state, next)
	}

	sum := 0
	for _, s := range state {
		// fmt.Println(s)
		if s.on {
			sum += s.cuboid.volume()
		} else {
			sum -= s.cuboid.volume()

		}
	}
	fmt.Println(sum)
}

func merge(a, b []Step) []Step {
	xs := map[Cuboid]int{}
	for _, s := range append(a, b...) {
		if _, b := xs[s.cuboid]; !b {
			xs[s.cuboid] = 0
		}
		if s.on {
			xs[s.cuboid]++
		} else {
			xs[s.cuboid]--
		}
	}
	// fmt.Println(xs)
	res := []Step{}
	for c, n := range xs {
		if 0 < n {
			for i := 0; i < n; i++ {
				res = append(res, Step{true, c})
			}
		} else if n < 0 {
			for i := 0; i < -n; i++ {
				res = append(res, Step{false, c})
			}

		}
	}
	return res
}

type Field struct {
	minmax Range
	vs     []bool
}

func (f *Field) init(r Range) {
	dim := r.max - r.min
	f.vs = make([]bool, dim*dim*dim)
	f.minmax = r
}

func (f *Field) apply(step Step) {
	dim, offset := f.minmax.max-f.minmax.min, -f.minmax.min
	c := step.cuboid.intersect(Cuboid{f.minmax, f.minmax, f.minmax})
	for k := c.z.min; k < c.z.max; k++ {
		for j := c.y.min; j < c.y.max; j++ {
			for i := c.x.min; i < c.x.max; i++ {
				idx := encode(i+offset, j+offset, k+offset, dim)
				f.vs[idx] = step.on
			}
		}
	}
}

func minimum(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func maximum(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type Range struct {
	min, max int
}

func (r1 Range) intersect(r2 Range) (res Range) {
	res.min = maximum(r1.min, r2.min)
	res.max = minimum(r1.max, r2.max)
	if res.max <= res.min {
		res = Range{}
	}
	return
}

type Cuboid struct {
	x, y, z Range
}

func (c Cuboid) volume() int {
	zero := Range{}
	if c.x == zero || c.y == zero || c.z == zero {
		return 0
	}
	return (c.x.max - c.x.min) * (c.y.max - c.y.min) * (c.z.max - c.z.min)
}

func (c1 Cuboid) intersect(c2 Cuboid) (res Cuboid) {
	res.x = c1.x.intersect(c2.x)
	res.y = c1.y.intersect(c2.y)
	res.z = c1.z.intersect(c2.z)
	return
}

type Step struct {
	on     bool
	cuboid Cuboid
}

func parse(input string) []Step {
	lines := strings.Split(input, "\n")
	res := []Step{}
	for _, line := range lines {
		if 0 < len(line) {
			res = append(res, parseStep(line))
		}
	}
	return res
}

func parseStep(s string) Step {
	step := Step{}
	ss := strings.Split(s, " ")
	step.on = ss[0] == "on"
	eqs := strings.Split(ss[1], ",")
	for _, eq := range eqs {
		kv := strings.Split(eq, "=")
		minmax := strings.Split(kv[1], "..")
		min, _ := strconv.Atoi(minmax[0])
		max, _ := strconv.Atoi(minmax[1])
		r := Range{min, max + 1}
		switch kv[0] {
		case "x":
			step.cuboid.x = r
		case "y":
			step.cuboid.y = r
		case "z":
			step.cuboid.z = r
		}
	}
	return step
}

func encode(x, y, z, n int) int {
	return x + (y * n) + z*n*n
}

func decode(v, n int) (x, y, z int) {
	x = v % n
	y = (v / n) % n
	z = v / n / n
	return
}
