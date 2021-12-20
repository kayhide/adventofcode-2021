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
	scanners := parse(input)
	s0 := scanners[0]
	done := map[int]bool{0: true}
	for len(done) < len(scanners) {
		last := len(done)
		for i, t := range scanners {
			if _, b := done[i]; !b {
				c := s0.match(&t)
				if 12 <= c.n {
					t.transform = c.f
					// t.display()
					// s0.display()
					s0.merge(&t)
					// s0.display()
					done[i] = true
				}
				// fmt.Println(i, n, len(s0.poses))
				fmt.Print(".")
			}
		}
		if last == len(done) {
			panic("Something is wrong")
		}
		fmt.Print("*")
	}

	fmt.Println("")
	fmt.Println(len(s0.poses))
}

func run2(input string) {
	scanners := parse(input)
	s0 := scanners[0]
	done := map[int]bool{0: true}
	poses := make([]Pos, len(scanners))
	for len(done) < len(scanners) {
		last := len(done)
		for i, t := range scanners {
			if _, b := done[i]; !b {
				c := s0.match(&t)
				if 12 <= c.n {
					t.transform = c.f
					s0.merge(&t)
					done[i] = true
					poses[i] = c.f(Pos{0, 0, 0})
				}
				// fmt.Println(i, n, len(s0.poses))
				fmt.Print(".")
			}
		}
		if last == len(done) {
			panic("Something is wrong")
		}
		fmt.Print("*")
	}
	fmt.Println("")

	max := 0
	for i, p := range poses {
		for _, p1 := range poses[i:] {
			d := abs(p.x-p1.x) + abs(p.y-p1.y) + abs(p.z-p1.z)
			if max < d {
				max = d
			}

		}
	}
	// fmt.Println(poses)
	fmt.Println(max)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

type Pos struct {
	x, y, z int
}

type Scanner struct {
	poses     map[Pos]int
	transform func(Pos) Pos
}

func parse(input string) []Scanner {
	lines := strings.Split(input, "\n")
	res := []Scanner{}
	poses := map[Pos]int{}
	for _, line := range lines {
		if 0 < len(line) {
			if line[0:3] == "---" {
				if 0 < len(poses) {
					scanner := Scanner{poses, func(p Pos) Pos { return p }}
					res = append(res, scanner)
					poses = map[Pos]int{}
				}
			} else {
				poses[parsePos(line)] = 1
			}
		}
	}
	if 0 < len(poses) {
		scanner := Scanner{poses, func(p Pos) Pos { return p }}
		res = append(res, scanner)
	}
	return res
}

func parsePos(s string) Pos {
	vs := strings.Split(s, ",")
	x, _ := strconv.Atoi(vs[0])
	y, _ := strconv.Atoi(vs[1])
	z, _ := strconv.Atoi(vs[2])
	return Pos{x, y, z}
}

func (s *Scanner) display() {
	area := 4000
	w, h := 20, 20
	xys := make([]int, w*h)
	for p, _ := range s.poses {
		p = s.transform(p)
		if -area < p.x && p.x < area && -area < p.y && p.y < area {
			xys[w*((p.y+area)/(2*area/h))+((p.x+area)/(2*area/w))]++
		}
	}
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if 0 < xys[w*j+i] {
				fmt.Print(xys[w*j+i])
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func (s *Scanner) merge(t *Scanner) {
	for p := range t.poses {
		s.poses[t.transform(p)] = 1
	}
}

type Transform func(Pos) Pos

type Count struct {
	n int
	f Transform
}

func (s *Scanner) match(t *Scanner) Count {
	ch := make(chan Count)
	for _, f := range variations() {
		go s.countLappings(t, f, ch)
	}

	max := Count{}
	max.n = 1
	for range variations() {
		for range s.poses {
			for range t.poses {
				r := <-ch
				if max.n < r.n {
					max = r
				}
			}
		}
	}

	return max
}

func (s *Scanner) countLappings(t *Scanner, f Transform, ch chan Count) {
	for sp := range s.poses {
		for tp := range t.poses {
			tp = f(tp)
			v := Pos{sp.x - tp.x, sp.y - tp.y, sp.z - tp.z}
			g := func(p Pos) Pos {
				p = f(p)
				return Pos{p.x + v.x, p.y + v.y, p.z + v.z}
			}
			go countDups(s.poses, t.poses, g, ch)
		}
	}
}

func countDups(ss map[Pos]int, ts map[Pos]int, f Transform, ch chan Count) {
	res := 0
	for t, _ := range ts {
		if _, b := ss[f(t)]; b {
			res++
		}
	}
	ch <- Count{res, f}
}

func variations() []Transform {
	res := []Transform{
		func(p Pos) Pos { return p },
		func(p Pos) Pos { return Pos{p.y, p.z, p.x} },
		func(p Pos) Pos { return Pos{p.z, p.x, p.y} },

		func(p Pos) Pos { return Pos{p.y, -p.x, p.z} },
		func(p Pos) Pos { return Pos{-p.x, p.z, p.y} },
		func(p Pos) Pos { return Pos{p.z, p.y, -p.x} },

		func(p Pos) Pos { return Pos{-p.x, -p.y, p.z} },
		func(p Pos) Pos { return Pos{-p.y, p.z, -p.x} },
		func(p Pos) Pos { return Pos{p.z, -p.x, -p.y} },

		func(p Pos) Pos { return Pos{-p.y, p.x, p.z} },
		func(p Pos) Pos { return Pos{p.x, p.z, -p.y} },
		func(p Pos) Pos { return Pos{p.z, -p.y, p.x} },

		func(p Pos) Pos { return Pos{p.y, p.x, -p.z} },
		func(p Pos) Pos { return Pos{p.x, -p.z, p.y} },
		func(p Pos) Pos { return Pos{-p.z, p.y, p.x} },

		func(p Pos) Pos { return Pos{-p.x, p.y, -p.z} },
		func(p Pos) Pos { return Pos{p.y, -p.z, -p.x} },
		func(p Pos) Pos { return Pos{-p.z, -p.x, p.y} },

		func(p Pos) Pos { return Pos{-p.y, -p.x, -p.z} },
		func(p Pos) Pos { return Pos{-p.x, -p.z, -p.y} },
		func(p Pos) Pos { return Pos{-p.z, -p.y, -p.x} },

		func(p Pos) Pos { return Pos{p.x, -p.y, -p.z} },
		func(p Pos) Pos { return Pos{-p.y, -p.z, p.x} },
		func(p Pos) Pos { return Pos{-p.z, p.x, -p.y} }}
	return res
}
