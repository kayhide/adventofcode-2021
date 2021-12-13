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
	cavern := parse(input)
	// cavern.display()
	for cavern.steps < 100 {
		cavern.step()
	}
	// cavern.display()
	sum := cavern.totalFlashes()
	fmt.Println(sum)
}

func run2(input string) {
	cavern := parse(input)
	// cavern.display()
	count := len(cavern.poses())
	lastFlashes := 0
	for i := 0; ; i++ {
		cavern.step()
		totalFlashes := cavern.totalFlashes()
		if lastFlashes+count == totalFlashes {
			break
		}
		lastFlashes = totalFlashes
	}
	// cavern.display()
	fmt.Println(cavern.steps)
}

type Pos struct {
	x, y int
}

type Octopus struct {
	level   int
	flashes int
}
type Cavern struct {
	octopuses [][]Octopus
	steps     int
}

func parse(input string) Cavern {
	lines := strings.Split(input, "\n")
	res := [][]Octopus{}
	for _, line := range lines {
		if 0 < len(line) {
			xs := []Octopus{}
			for _, c := range line {
				n, _ := strconv.Atoi(string(c))
				xs = append(xs, Octopus{n, 0})
			}
			res = append(res, xs)
		}
	}
	return Cavern{res, 0}
}

func (cavern *Cavern) display() {
	for _, row := range cavern.octopuses {
		for _, n := range row {
			if n.level < 10 {
				fmt.Printf("%d", n.level)
			} else {
				fmt.Printf("%c", n.level-10+int('a'))
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (cavern *Cavern) poses() (res []Pos) {
	for j, row := range cavern.octopuses {
		for i := range row {
			res = append(res, Pos{i, j})
		}
	}
	return
}

func (cavern *Cavern) isInside(p Pos) bool {
	return 0 <= p.x &&
		p.x < len(cavern.octopuses[0]) &&
		0 <= p.y &&
		p.y < len(cavern.octopuses)

}

func (cavern *Cavern) at(p Pos) *Octopus {
	return &cavern.octopuses[p.y][p.x]
}

func (cavern *Cavern) neighbors(p Pos) (res []Pos) {
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i != 0 || j != 0 {
				p := Pos{p.x + i, p.y + j}
				if cavern.isInside(p) {
					res = append(res, p)
				}
			}
		}
	}
	return
}

func (cavern *Cavern) step() {
	cavern.inc()
	cavern.flash()
	cavern.steps++
}

func (cavern *Cavern) inc() {
	for _, p := range cavern.poses() {
		octopus := cavern.at(p)
		octopus.level++
	}
}

func (cavern *Cavern) flash() {
	for {
		flashings := []Pos{}
		for _, p := range cavern.poses() {
			if 9 < cavern.at(p).level {
				flashings = append(flashings, p)
			}
		}
		if len(flashings) == 0 {
			return
		}
		for _, p := range flashings {
			octopus := cavern.at(p)
			octopus.level = 0
			octopus.flashes++
			cavern.chargeNeighbors(p)
			// cavern.display()
		}
	}
}

func (cavern *Cavern) chargeNeighbors(p Pos) {
	for _, q := range cavern.neighbors(p) {
		octopus := cavern.at(q)
		if 0 < octopus.level {
			octopus.level++
		}
	}
}

func (cavern *Cavern) totalFlashes() (sum int) {
	for _, p := range cavern.poses() {
		sum += cavern.at(p).flashes
	}
	return
}
