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
	insts := parse(input)
	blocks := toBlocks(insts)

	runner := Runner{}
	runner.init(blocks)

	res, _ := runner.descend([]int{}, Alu{})
	fmt.Println("")
	for _, n := range res {
		fmt.Print(n)
	}
	fmt.Println("")
}

func run2(input string) {
	insts := parse(input)
	blocks := toBlocks(insts)

	runner := Runner{}
	runner.init(blocks)

	res, _ := runner.ascend([]int{}, Alu{})
	fmt.Println("")
	for _, n := range res {
		fmt.Print(n)
	}
	fmt.Println("")
}

type Runner struct {
	blocks []Block
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
	if len(r.cache)%100000 == 0 {
		fmt.Print(".")
	}
	if 100000000 < len(r.cache) {
		r.cache = map[CacheEntry]bool{}
	}
	r.cache[toCacheEntry(ix, input, a)] = true
	return r.blocks[ix].exec(input, a)
}

func (r *Runner) visited(ix, input int, a Alu) bool {
	_, b := r.cache[toCacheEntry(ix, input, a)]
	return b
}

func (r *Runner) descend(inputs []int, a Alu) ([]int, bool) {
	ix := len(inputs)
	if ix == 14 {
		// fmt.Println(inputs, len(r.cache), a)
		return inputs, a[2] == 0
	}

	for i := 9; 0 < i; i-- {
		if !r.visited(ix, i, a) {
			a1 := r.exec(ix, i, a)
			ns, b := r.descend(append([]int{}, append(inputs, i)...), a1)
			if b {
				return ns, true
			}
		}
	}
	return nil, false
}

func (r *Runner) ascend(inputs []int, a Alu) ([]int, bool) {
	ix := len(inputs)
	if ix == 14 {
		// fmt.Println(inputs, len(r.cache), a)
		return inputs, a[2] == 0
	}

	for i := 1; i <= 9; i++ {
		if !r.visited(ix, i, a) {
			a1 := r.exec(ix, i, a)
			ns, b := r.ascend(append([]int{}, append(inputs, i)...), a1)
			if b {
				return ns, true
			}
		}
	}
	return nil, false
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

func toBlocks(insts []Inst) []Block {
	res := []Block{}
	block := Block{}
	for ix, inst := range insts {
		switch i := inst.(type) {
		case Inp:
			if 0 < len(block.ops) {
				res = append(res, block)
				block = Block{}
				block.ix = ix
				block.ops = []Op{}
			}
			block.inp = i
		case Op:
			block.ops = append(block.ops, i)

		}
	}
	res = append(res, block)
	return res
}

func toInput(n int) ([]int, bool) {
	res := make([]int, 14)
	for i := 0; i < 14; i++ {
		if n%10 == 0 {
			return res, false
		}
		res[13-i] = n % 10
		n = n / 10
	}
	return res, true
}

type Alu [4]int

type Block struct {
	ix  int
	inp Inp
	ops []Op
}

type Inst interface {
	show()
}

type Inp struct {
	line string
	f    func(int, Alu) Alu
}

func (i Inp) show() {
	fmt.Println(i.line)
}

type Op struct {
	line string
	f    func(Alu) Alu
}

func (i Op) show() {
	fmt.Println(i.line)
}

func parse(input string) []Inst {
	lines := strings.Split(input, "\n")
	res := []Inst{}
	for _, line := range lines {
		if 0 < len(line) {
			res = append(res, parseInst(line))
		}
	}
	return res
}

func parseInst(s string) Inst {
	ss := strings.Split(s, " ")
	switch ss[0] {
	case "inp":
		i := toIdx(ss[1])
		return Inp{
			s,
			func(n int, a Alu) Alu {
				res := a
				res[i] = n
				return res
			}}
	case "add":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				s,
				func(alu Alu) Alu {
					alu[a] = alu[a] + b
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				s,
				func(alu Alu) Alu {
					alu[a] = alu[a] + alu[b]
					return alu
				}}

		}
	case "mul":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				s,
				func(alu Alu) Alu {
					alu[a] = alu[a] * b
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				s,
				func(alu Alu) Alu {
					alu[a] = alu[a] * alu[b]
					return alu
				}}

		}
	case "div":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				s,
				func(alu Alu) Alu {
					alu[a] = alu[a] / b
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				s,
				func(alu Alu) Alu {
					alu[a] = alu[a] / alu[b]
					return alu
				}}

		}
	case "mod":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				s,
				func(alu Alu) Alu {
					alu[a] = alu[a] % b
					return alu
				}}
		} else {
			b := toIdx(ss[2])
			return Op{
				s,
				func(alu Alu) Alu {
					alu[a] = alu[a] % alu[b]
					return alu
				}}

		}
	case "eql":
		a := toIdx(ss[1])
		if b, err := strconv.Atoi(ss[2]); err == nil {
			return Op{
				s,
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
				s,
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
