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
	algo, field := parse(input)
	// algo.display()
	// field.display()
	for i := 0; i < 2; i++ {
		outer := i%2 == 1 && algo[0]
		field = field.enhance(outer, algo)

	}
	fmt.Println(field.count())
}

func run2(input string) {
	algo, field := parse(input)
	// algo.display()
	// field.display()
	for i := 0; i < 50; i++ {
		outer := i%2 == 1 && algo[0]
		field = field.enhance(outer, algo)

	}
	fmt.Println(field.count())
}

type Algo []bool

type Pos struct {
	x, y int
}

type Field [][]bool

func (field Field) count() (sum int) {
	for _, row := range field {
		for _, b := range row {
			if b {
				sum++
			}
		}
	}
	return sum
}

func parse(input string) (Algo, Field) {
	lines := strings.Split(input, "\n")
	algo := parseBools(lines[0])
	bss := [][]bool{}
	for _, line := range lines[1:] {
		if 0 < len(line) {
			bss = append(bss, parseBools(line))
		}
	}
	return algo, Field(bss)
}

func parseBools(s string) []bool {
	bs := make([]bool, len(s))
	for i, c := range s {
		bs[i] = c == '#'
	}
	return bs
}

func printRow(bs []bool) {
	for _, b := range bs {
		if b {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println("")
}

func (algo Algo) display() {
	printRow([]bool(algo))
}

func (field Field) display() {
	for _, row := range field {
		printRow(row)
	}
	fmt.Println("")
}

func (field Field) at(outer bool, p Pos) bool {
	w, h := len(field[0]), len(field)
	if 0 <= p.x && p.x < w && 0 <= p.y && p.y < h {
		return field[p.y][p.x]
	}
	return outer
}

func (field Field) enhance(outer bool, algo Algo) Field {
	w, h := len(field[0]), len(field)
	bss := make([][]bool, h+2)
	for j := 0; j < h+2; j++ {
		bss[j] = make([]bool, w+2)
		for i := 0; i < w+2; i++ {
			v := field.windowValue(outer, Pos{i - 1, j - 1})
			bss[j][i] = algo[v]
		}
	}
	return Field(bss)
}

func (field Field) windowValue(outer bool, p Pos) int {
	bs := make([]bool, 9)
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			bs[j*3+i+4] = field.at(outer, Pos{p.x + i, p.y + j})
		}
	}
	return toInt(bs)
}

func toInt(bs []bool) int {
	res := 0
	for _, b := range bs {
		res = res << 1
		if b {
			res++
		}
	}
	return res
}
