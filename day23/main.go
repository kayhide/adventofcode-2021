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
	layout := parse(input)
	layout.display()

	for l, c := range layout.explore(0) {
		l.display()
		fmt.Println(c)
	}
}

func run2(input string) {
	layout := parse2(input)
	layout.display()

	for l, c := range layout.explore(0) {
		l.display()
		fmt.Println(c)
	}
}

func merge(xs, ys CostMap) CostMap {
	for l, c := range ys {
		if x, b := xs[l]; !b || c < x {
			xs[l] = c
		}
	}
	return xs
}

type Amph rune

func (a Amph) display() {
	if a.empty() {
		fmt.Print("_")
	}
	fmt.Print(string(a))
}

func (a Amph) empty() bool {
	return a == 0
}

func (a Amph) cost() int {
	switch rune(a) {
	case 'A':
		return 1
	case 'B':
		return 10
	case 'C':
		return 100
	case 'D':
		return 1000
	}
	return 0
}

func (a Amph) destination() int {
	switch rune(a) {
	case 'A':
		return 2
	case 'B':
		return 4
	case 'C':
		return 6
	case 'D':
		return 8
	}
	return 0
}

type CostMap map[Layout]int

// type Layout [15]Amph
// type Layout [15 + 8]Amph
type Layout struct {
	size  int
	amphs [15 + 8]Amph
}

func (l Layout) depth() int {
	return (l.size-3)/4 - 1
}

func (l Layout) explore(cost int) CostMap {
	res := CostMap{}
	if l.done() {
		res[l] = cost
		return res
	}
	for l2, c := range l.step() {
		merge(res, l2.explore(c+cost))
	}
	return res
}

func (l Layout) step() CostMap {
	res := CostMap{}
	if i, j, b := l.findFixable(); b {
		next := l
		c := next.move(i, j)
		res[next] = c
		return res
	}
	for i := 0; i < 11; i++ {
		for j := 0; j < 11; j++ {
			if l.reachable(i, j) {
				next := l
				c := next.move(i, j)
				res[next] = c
			}
		}
	}
	return res
}

func (l Layout) findFixable() (src, dst int, b bool) {
	for i := 0; i < 11; i++ {
		if l.any(i) {
			j := l.top(i).destination()
			if l.reachable(i, j) {
				src, dst, b = i, j, true
				return
			}
		}
	}
	src, dst, b = 0, 0, false
	return
}

func (l Layout) done() bool {
	for _, i := range []int{2, 4, 6, 8} {
		if !l.full(i) || !l.fixed(i) {
			return false
		}
	}
	return true
}

func (l Layout) cost(src, dst int) int {
	steps := dst - src
	if steps < 0 {
		steps = -steps
	}

	amph := l.top(src)
	if l.room(src) {
		steps++
		for _, a := range l.get(src) {
			if a.empty() {
				steps++
			}
		}
	}
	if l.room(dst) {
		for _, a := range l.get(dst) {
			if a.empty() {
				steps++
			}
		}
	}
	return steps * amph.cost()
}

func (l Layout) reachable(src, dst int) bool {
	if src == dst {
		return false
	}
	if l.empty(src) {
		return false
	}
	if l.fixed(src) {
		return false
	}
	if l.full(dst) {
		return false
	}
	if l.hall(src) && l.hall(dst) {
		return false
	}
	if l.room(dst) {
		if l.top(src).destination() != dst {
			return false
		}
		if !l.empty(dst) && !l.fixed(dst) {
			return false
		}
	}

	min, max := src, dst
	if dst < src {
		min, max = dst, src
	}

	for i := min + 1; i < max; i++ {
		if l.hall(i) && l.full(i) {
			return false
		}
	}
	return true
}

func (l *Layout) move(s, t int) int {
	c := l.cost(s, t)
	a := l.pop(s)
	l.push(t, a)
	return c
}

func (l *Layout) pop(i int) Amph {
	as := l.get(i)
	j := len(as) - 1
	for ; 0 <= j; j-- {
		if !as[j].empty() {
			break
		}
	}
	a := as[j]
	as[j] = Amph(0)
	l.put(i, as)
	return a
}

func (l *Layout) push(i int, a Amph) {
	as := l.get(i)
	j := 0
	for ; j < len(as); j++ {
		if as[j].empty() {
			break
		}
	}
	as[j] = a
	l.put(i, as)
}

func (l *Layout) top(i int) Amph {
	amph := Amph(0)
	for _, a := range l.get(i) {
		if !a.empty() {
			amph = a
		}
	}
	return amph
}

func (l Layout) fixed(i int) bool {
	for _, a := range l.get(i) {
		if !a.empty() && a.destination() != i {
			return false
		}
	}
	return true
}

func (l Layout) get(i int) []Amph {
	res := []Amph{}
	d := l.depth()
	if i < 2 {
		res = append(res, l.amphs[i])
	} else if i < 9 {
		j := 2 + (d+1)*((i-2)/2)
		if l.hall(i) {
			res = append(res, l.amphs[j+d])
		} else {
			res = append(res, l.amphs[j:j+d]...)
		}
	} else {
		j := i + (d-1)*4
		res = append(res, l.amphs[j])
	}
	return res
}

func (l *Layout) put(i int, as []Amph) {
	d := l.depth()
	if i < 2 {
		l.amphs[i] = as[0]
	} else if i < 9 {
		j := 2 + (d+1)*((i-2)/2)
		if l.hall(i) {
			l.amphs[j+d] = as[0]
		} else {
			for k := 0; k < d; k++ {
				l.amphs[j+k] = as[k]
			}
		}
	} else {
		j := i + (d-1)*4
		l.amphs[j] = as[0]
	}

}

func (l Layout) hall(i int) bool {
	return !l.room(i)
}

func (l Layout) room(i int) bool {
	return i == 2 || i == 4 || i == 6 || i == 8
}

func (l Layout) empty(i int) bool {
	for _, a := range l.get(i) {
		if !a.empty() {
			return false
		}
	}
	return true
}

func (l Layout) full(i int) bool {
	for _, a := range l.get(i) {
		if a.empty() {
			return false
		}
	}
	return true
}

func (l Layout) any(i int) bool {
	for _, a := range l.get(i) {
		if !a.empty() {
			return true
		}
	}
	return false
}

func (l Layout) display() {
	amphss := map[int][]Amph{}
	for i := 0; i < 11; i++ {
		amphss[i] = l.get(i)
	}
	for i := 0; i < 11; i++ {
		if l.hall(i) {
			amphss[i][0].display()
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println("")
	for j := l.depth() - 1; 0 <= j; j-- {
		for i := 0; i < 11; i++ {
			if l.hall(i) {
				fmt.Print(" ")
			} else {
				amphss[i][j].display()
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func parse(input string) Layout {
	lines := strings.Split(input, "\n")
	lines = lines[2:4]
	layout := Layout{}
	layout.size = 7 + 4*len(lines)
	for j := 0; j < len(lines); j++ {
		line := lines[len(lines)-j-1]
		for i := 0; i < 4; i++ {
			layout.push(i*2+2, Amph(rune(line[i*2+3])))
		}
	}
	return layout
}

func parse2(input string) Layout {
	lines := strings.Split(input, "\n")
	extra := []string{
		"  #D#C#B#A#",
		"  #D#B#A#C#"}
	lines = append(lines[2:3], extra[0], extra[1], lines[3])
	layout := Layout{}
	layout.size = 7 + 4*len(lines)
	for j := 0; j < len(lines); j++ {
		line := lines[len(lines)-j-1]
		for i := 0; i < 4; i++ {
			layout.push(i*2+2, Amph(rune(line[i*2+3])))
		}
	}
	return layout
}
