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
	target := parse(input)

	max := Toss{}
	for vy := 0; vy < -target.p0.y+1; vy++ {
		toss := Toss{}
		toss.v0 = Pos{0, vy}
		toss.runY(target)
		if toss.success && max.hi < toss.hi {
			// fmt.Println(toss)
			max = toss
		}
	}
	fmt.Printf("%d\n", max.hi)
}

func run2(input string) {
	target := parse(input)

	count := 0
	for vx := 1; vx < target.p1.x+1; vx++ {
		for vy := target.p0.y; vy < -target.p0.y+1; vy++ {
			toss := Toss{}
			toss.v0 = Pos{vx, vy}
			toss.runXY(target)
			if toss.success {
				// fmt.Println(toss)
				count++
			}
		}
	}
	fmt.Printf("%d\n", count)
}

type Pos struct {
	x, y int
}

type Target struct {
	p0, p1 Pos
}

func parse(input string) Target {
	lines := strings.Split(input, "\n")
	line := lines[0]
	line = strings.TrimSpace(strings.Split(line, ":")[1])
	xys := strings.Split(line, ",")
	xs := strings.Split(strings.Split(xys[0], "=")[1], "..")
	ys := strings.Split(strings.Split(xys[1], "=")[1], "..")

	// fmt.Println(xs, ys)
	x0, _ := strconv.Atoi(xs[0])
	x1, _ := strconv.Atoi(xs[1])
	y0, _ := strconv.Atoi(ys[0])
	y1, _ := strconv.Atoi(ys[1])
	return Target{Pos{x0, y0}, Pos{x1, y1}}
}

type Toss struct {
	v0      Pos
	hi      int
	success bool
}

func (toss *Toss) runY(target Target) {
	p := Pos{}
	v := toss.v0
	for target.p1.y <= p.y {
		// fmt.Println(p.y, v.y)
		p.y += v.y
		if v.y == 0 {
			toss.hi = p.y
		}
		v.y--
	}
	toss.success = target.p0.y <= p.y
}

func (toss *Toss) runXY(target Target) {
	p := Pos{}
	v := toss.v0
	for target.p0.y < p.y {
		// fmt.Println(p, v)
		p.x += v.x
		p.y += v.y
		if isInside(target, p) {
			toss.success = true
			return
		}
		if 0 < v.x {
			v.x--
		}
		v.y--
	}
	toss.success = false
}

func isInside(target Target, p Pos) bool {
	return target.p0.x <= p.x && p.x <= target.p1.x && target.p0.y <= p.y && p.y <= target.p1.y
}
