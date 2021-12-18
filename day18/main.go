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
	// input = "[[[[[9,8],1],2],3],4]"
	// input = "[7,[6,[5,[4,[3,2]]]]]"
	// input = "[[6,[5,[4,[3,2]]]],1]"
	// input = "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]"
	// input = "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]"
	// input = "[[[[4,3],4],4],[7,[[8,4],9]]]\n[1,1]"
	// input = "[1,1]\n[2,2]\n[3,3]\n[4,4]"
	// input = "[1,1]\n[2,2]\n[3,3]\n[4,4]\n[5,5]"
	// input = "[1,1]\n[2,2]\n[3,3]\n[4,4]\n[5,5]\n[6,6]"
	// input = "[[[4,0],[5,4]],[[7,7],[6,0]]]\n[[1,[5,5]],[[0,[[5,5],[5,5]]],[0,6]]]"

	ns := parse(input)
	res := ns[0]
	for _, n := range ns[1:] {
		// fmt.Println(res)
		// fmt.Println(" + ", n)
		res = add(res, n)
		// fmt.Println(" = ", res)
	}
	// fmt.Printf("%d\n", res)
	fmt.Println(res.magnitude())
}

func run2(input string) {
	ns := parse(input)
	max := 0
	for _, n := range ns {
		for _, m := range ns {
			if n != m {
				v := reduce(add(n, m)).magnitude()
				if max < v {
					max = v

				}
			}
		}
		fmt.Print(".")
	}
	// fmt.Printf("%d\n", res)
	fmt.Println("")
	fmt.Println(max)
}

type Num interface {
	paths() []Path
	at(Path) Num
	over(Path, func(Num) Num) Num
	magnitude() int
}

type Atom int

type Pair struct {
	x, y Num
}

func parse(input string) []Num {
	lines := strings.Split(input, "\n")
	res := make([]Num, 0)
	for _, line := range lines {
		if 0 < len(line) {
			n, _ := parseNum(line)
			res = append(res, n)
		}
	}
	return res
}

func parseNum(s string) (Num, int) {
	i := 0
	switch s[i] {
	case '[':
		n, j := parseNum(s[i+1:])
		i += j + 1
		m, k := parseNum(s[i+1:])
		i += k + 1
		return Pair{n, m}, i + 1
	default:
		n, _ := strconv.Atoi(string(s[i]))
		return Atom(n), i + 1
	}
}

func add(n, m Num) Num {
	return reduce(Pair{n, m})
}

func reduce(n Num) Num {
	for {
		m := split(explode(n))
		if n == m {
			return m
		}
		n = m
	}
}

func explode(n Num) Num {
	updated := true
	for updated {
		updated = false
		paths := n.paths()
		for i, path := range paths {
			if 5 <= len(path) && isExplodable(n.at(path.up())) {
				l, r := n.at(path).(Atom), n.at(paths[i+1]).(Atom)
				if 0 < i {
					n = n.over(paths[i-1], func(m Num) Num {
						v, _ := m.(Atom)
						return Atom(int(l) + int(v))
					})
				}

				i, path = i+1, paths[i+1]
				if i < len(paths)-1 {
					n = n.over(paths[i+1], func(m Num) Num {
						v, _ := m.(Atom)
						return Atom(int(r) + int(v))
					})
				}
				n = n.over(path.up(), func(m Num) Num { return Atom(0) })
				updated = true
				// fmt.Println(" => ", n)
				break
			}
		}
	}
	return n
}

func split(n Num) Num {
	paths := n.paths()
	for _, path := range paths {
		v := n.at(path).(Atom)
		if 9 < v {
			n = n.over(path, func(Num) Num { return Pair{Atom(v / 2), Atom((v + 1) / 2)} })
			return n
		}
	}
	return n
}

func isExplodable(n Num) bool {
	if p, b := n.(Pair); b {
		_, b1 := p.x.(Atom)
		_, b2 := p.y.(Atom)
		return b1 && b2
	}
	return false
}

type Path []bool

func (path Path) up() Path {
	res := make(Path, len(path)-1)
	copy(res, path)
	return res
}

func (path Path) left() (Path, bool) {
	i := -1
	for j, b := range path {
		if b {
			i = j
		}
	}
	if i == -1 {
		return path, false
	}
	res := make(Path, len(path))
	copy(res, path)
	res[i] = false
	return res, true
}

func (path Path) right() (Path, bool) {
	i := -1
	for j, b := range path {
		if !b {
			i = j
		}
	}
	if i == -1 {
		return path, false
	}
	res := make(Path, len(path))
	copy(res, path)
	res[i] = !res[i]
	return res, true
}

func (n Atom) paths() []Path {
	return []Path{Path([]bool{})}
}
func (n Pair) paths() (res []Path) {
	for _, path := range n.x.paths() {
		res = append(res, append(Path([]bool{false}), path...))
	}
	for _, path := range n.y.paths() {
		res = append(res, append(Path([]bool{true}), path...))
	}
	return
}

func (n Atom) at(path Path) Num {
	if 0 < len(path) {
		panic("Wrong path")
	}
	return n
}

func (n Pair) at(path Path) Num {
	if 0 == len(path) {
		return n
	}
	if path[0] {
		return n.y.at(path[1:])
	}
	return n.x.at(path[1:])
}

func (n Atom) over(path Path, f func(Num) Num) Num {
	if 0 < len(path) {
		panic("Wrong path")
	}
	return f(n)
}

func (n Pair) over(path Path, f func(Num) Num) Num {
	if 0 == len(path) {
		panic("Wrong path")
	}
	if 1 == len(path) {
		if path[0] {
			return Pair{n.x, f(n.y)}
		}
		return Pair{f(n.x), n.y}
	}
	if path[0] {
		return Pair{n.x, n.y.over(path[1:], f)}
	}
	return Pair{n.x.over(path[1:], f), n.y}
}

func (n Atom) magnitude() int {
	return int(n)
}

func (n Pair) magnitude() int {
	return 3*n.x.magnitude() + 2*n.y.magnitude()
}
