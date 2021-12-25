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
	// run2(string(input))
}

func run1(input string) {
	layout := parse(input)
	// layout.display()

	n := 0
	for {
		n++
		next := layout.step()
		if next.eql(&layout) {
			break
		}
		layout = next
	}
	fmt.Println(n)
}

func run2(input string) {
}

type Cucumber rune

func (c Cucumber) eastward() bool {
	return c == '>'
}

func (c Cucumber) southward() bool {
	return c == 'v'
}

func (c Cucumber) empty() bool {
	return c == 0
}

func (c Cucumber) display() {
	if c.empty() {
		fmt.Print(".")
	} else {
		fmt.Print(string(c))
	}
}

type Layout struct {
	w, h      int
	cucumbers []Cucumber
}

func (l *Layout) eql(l2 *Layout) bool {
	for i, c := range l.cucumbers {
		if c != l2.cucumbers[i] {
			return false
		}
	}
	return true
}

func (l *Layout) empty(x, y int) bool {
	return l.get(x%l.w, y%l.h).empty()
}

func (l *Layout) ground() Layout {
	res := Layout{}
	res.w, res.h = l.w, l.h
	res.cucumbers = make([]Cucumber, len(l.cucumbers))
	return res
}

func (l *Layout) step() Layout {
	l1 := l.stepEast()
	l2 := l1.stepSouth()
	return l2
}
func (l *Layout) stepEast() Layout {
	res := l.ground()
	for j := 0; j < l.h; j++ {
		for i := 0; i < l.w; i++ {
			c := l.get(i, j)
			if c.eastward() {
				if l.empty(i+1, j) {
					res.put(i+1, j, c)
				} else {
					res.put(i, j, c)
				}
			} else if c.southward() {
				res.put(i, j, c)
			}
		}
	}
	return res
}

func (l *Layout) stepSouth() Layout {
	res := l.ground()
	for j := 0; j < l.h; j++ {
		for i := 0; i < l.w; i++ {
			c := l.get(i, j)
			if c.eastward() {
				res.put(i, j, c)
			} else if c.southward() {
				if l.empty(i, j+1) {
					res.put(i, j+1, c)
				} else {
					res.put(i, j, c)
				}
			}
		}
	}
	return res
}

func (l Layout) display() {
	for j := 0; j < l.h; j++ {
		for i := 0; i < l.w; i++ {
			l.get(i, j).display()
		}
		fmt.Println("")
	}
	fmt.Println("")

}

func (l *Layout) put(x, y int, c Cucumber) {
	l.cucumbers[(x%l.w)+(y%l.h)*l.w] = c
}

func (l *Layout) get(x, y int) Cucumber {
	return l.cucumbers[x+y*l.w]
}

func parse(input string) Layout {
	lines := strings.Split(input, "\n")
	layout := Layout{}
	layout.w = len(lines[0])
	for _, line := range lines {
		if 0 < len(line) {
			layout.h++
		}
	}
	layout.cucumbers = make([]Cucumber, layout.w*layout.h)
	for j := 0; j < layout.h; j++ {
		for i := 0; i < layout.w; i++ {
			c := lines[j][i]
			if c == '>' || c == 'v' {
				layout.put(i, j, Cucumber(c))
			}
		}
	}
	return layout
}
