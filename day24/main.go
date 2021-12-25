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
	blocks := parse(input)
	runner := Runner{}
	runner.init(blocks)
	is := make([]int, 9)
	for i := 0; i < 9; i++ {
		is[i] = 9 - i
	}
	runner.nexts = is

	ch := make(chan string)
	go runner.explore("", Alu{}, ch)

	res := <-ch
	fmt.Println("")
	fmt.Println(res)
}

func run2(input string) {
	blocks := parse(input)
	runner := Runner{}
	runner.init(blocks)
	is := make([]int, 9)
	for i := 0; i < 9; i++ {
		is[i] = i + 1
	}
	runner.nexts = is

	ch := make(chan string)
	go runner.explore("", Alu{}, ch)

	res := <-ch
	fmt.Println("")
	fmt.Println(res)
}

type Runner struct {
	blocks []Block
	nexts  []int
	cache  map[CacheEntry]bool
}

type CacheEntry [3]int

func toCacheEntry(ix, input int, alu Alu) CacheEntry {
	return CacheEntry{ix, input, alu[2]}
}

func (r *Runner) init(blocks []Block) {
	r.blocks = blocks
	r.cache = map[CacheEntry]bool{}
}

func (r *Runner) exec(ix, input int, a Alu) Alu {
	if 10000000 < len(r.cache) {
		fmt.Println("*")
		r.cache = map[CacheEntry]bool{}
	}
	if len(r.cache)%100000 == 0 {
		fmt.Print(".")
	}
	r.cache[toCacheEntry(ix, input, a)] = true
	return r.blocks[ix].exec(input, a)
}

func (r *Runner) visited(ix, input int, a Alu) bool {
	_, b := r.cache[toCacheEntry(ix, input, a)]
	return b
}

func (r *Runner) explore(inputs string, a Alu, ch chan string) bool {
	ix := len(inputs)
	if ix == len(r.blocks) {
		if a[2] == 0 {
			ch <- inputs
			return true
		}
		return false
	}

	for _, i := range r.nexts {
		if !r.visited(ix, i, a) {
			a1 := r.exec(ix, i, a)
			if r.explore(fmt.Sprintf("%s%d", inputs, i), a1, ch) {
				return true
			}
		}
	}
	return false
}

func (block Block) exec(i int, a Alu) Alu {
	a = block.inp.f(i, a)
	for _, op := range block.ops {
		// old := a
		a = op.f(a)
		// op.show()
		// fmt.Println(old, "->", a)
	}
	return a
}

type Alu [4]int

type Block struct {
	ix  int
	inp Inp
	ops []Op
}

func (block Block) empty() bool {
	return block.inp.f == nil
}

type Inp struct {
	line []string
	f    func(int, Alu) Alu
}

type Op struct {
	line []string
	f    func(Alu) Alu
}

func parse(input string) []Block {
	lines := strings.Split(input, "\n")
	res := []Block{}
	var block *Block
	for _, line := range lines {
		if 0 < len(line) {
			ss := strings.Split(line, " ")
			if inp, ok := parseInp(ss); ok {
				if block != nil {
					res = append(res, *block)
				}
				block = &Block{len(res), inp, []Op{}}
			} else {
				op := parseOp(ss)
				block.ops = append(block.ops, op)
			}
		}
	}
	res = append(res, *block)
	return res
}

func parseInp(ss []string) (Inp, bool) {
	switch ss[0] {
	case "inp":
		i := toIdx(ss[1])
		return Inp{
			ss,
			func(n int, a Alu) Alu {
				res := a
				res[i] = n
				return res
			}}, true
	}
	return Inp{}, false
}

func parseOp(ss []string) Op {
	switch ss[0] {
	case "add":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				ss,
				func(alu Alu) Alu {
					alu[a] = alu[a] + b
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				ss,
				func(alu Alu) Alu {
					alu[a] = alu[a] + alu[b]
					return alu
				}}

		}
	case "mul":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				ss,
				func(alu Alu) Alu {
					alu[a] = alu[a] * b
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				ss,
				func(alu Alu) Alu {
					alu[a] = alu[a] * alu[b]
					return alu
				}}

		}
	case "div":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				ss,
				func(alu Alu) Alu {
					alu[a] = alu[a] / b
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				ss,
				func(alu Alu) Alu {
					alu[a] = alu[a] / alu[b]
					return alu
				}}

		}
	case "mod":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				ss,
				func(alu Alu) Alu {
					alu[a] = alu[a] % b
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				ss,
				func(alu Alu) Alu {
					alu[a] = alu[a] % alu[b]
					return alu
				}}

		}
	case "eql":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				ss,
				func(alu Alu) Alu {
					if alu[a] == b {
						alu[a] = 1
					} else {
						alu[a] = 0
					}
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				ss,
				func(alu Alu) Alu {
					if alu[a] == alu[b] {
						alu[a] = 1
					} else {
						alu[a] = 0
					}
					return alu
				}}

		}
	}
	panic("Unknown instruction")
}

func toIdx(c string) int {
	switch c {
	case "x":
		return 0
	case "y":
		return 1
	case "z":
		return 2
	case "w":
		return 3
	}
	panic("Bad variable: " + c)
}
