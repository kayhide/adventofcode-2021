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
	inst := parse(input)
	for i := 0; i < 10; i++ {
		inst.step()
		// fmt.Println(inst.polymer)
	}

	counts := inst.count()
	min, max := len(inst.polymer), 0
	for _, v := range counts {
		if v < min {
			min = v
		}
		if max < v {
			max = v
		}
	}
	fmt.Println(max - min)
}

func run2(input string) {
	inst := parse(input)
	counter := inst.toCounter()
	for i := 0; i < 40; i++ {
		inst.stepCounter(&counter)
	}

	counts := counter.count()
	counts[inst.polymer[0]]++
	min, max := counts[inst.polymer[0]], counts[inst.polymer[0]]
	for _, v := range counts {
		if v < min {
			min = v
		}
		if max < v {
			max = v
		}
	}
	// fmt.Println(counts)
	// fmt.Println(min, max)
	fmt.Println(max - min)
}

type Pair struct {
	fst, snd byte
}

type Inst struct {
	polymer    []byte
	insertions map[Pair]byte
}

type Counter struct {
	counts map[Pair]int
}

func parse(input string) Inst {
	lines := strings.Split(input, "\n")
	polymer := []byte{}
	insertions := map[Pair]byte{}
	for _, c := range lines[0] {
		polymer = append(polymer, byte(c))
	}

	for _, line := range lines[1:] {
		if 0 < len(line) {
			xs := strings.Split(line, " -> ")
			insertions[Pair{xs[0][0], xs[0][1]}] = xs[1][0]
		}
	}
	return Inst{polymer, insertions}
}

func (inst *Inst) step() {
	next := make([]byte, len(inst.polymer)*2)
	next[0] = inst.polymer[0]
	p := 1
	for i, x := range inst.polymer[1:] {
		last := inst.polymer[i]
		if dst, b := inst.insertions[Pair{last, x}]; b {
			next[p] = dst
			next[p+1] = x
			p += 2
		} else {
			next[p] = x
			p++
		}
	}
	inst.polymer = next[:p]
}

func (inst *Inst) count() map[byte]int {
	res := map[byte]int{}
	for _, c := range inst.polymer {
		if _, b := res[c]; !b {
			res[c] = 1
		} else {
			res[c]++

		}
	}
	return res
}

func (inst *Inst) toCounter() Counter {
	counts := map[Pair]int{}
	for i, c := range inst.polymer[1:] {
		pair := Pair{inst.polymer[i], c}
		if _, b := counts[pair]; !b {
			counts[pair] = 1
		} else {
			counts[pair]++
		}
	}
	return Counter{counts}
}

func (inst *Inst) stepCounter(counter *Counter) {
	diff := map[Pair]int{}
	for k, n := range counter.counts {
		if x, b := inst.insertions[k]; b {
			p1, p2 := Pair{k.fst, x}, Pair{x, k.snd}
			ensure(&diff, k, 0)
			ensure(&diff, p1, 0)
			ensure(&diff, p2, 0)
			diff[k] -= n
			diff[p1] += n
			diff[p2] += n
		}
	}
	for k, d := range diff {
		ensure(&counter.counts, k, 0)
		counter.counts[k] += d
	}
}

func ensure(xs *map[Pair]int, k Pair, v int) {
	if _, b := (*xs)[k]; !b {
		(*xs)[k] = v
	}
}

func (counter *Counter) count() map[byte]int {
	res := map[byte]int{}
	for k, n := range counter.counts {
		c := k.snd
		if _, b := res[c]; !b {
			res[c] = 0
		}
		res[c] += n
	}
	return res
}
