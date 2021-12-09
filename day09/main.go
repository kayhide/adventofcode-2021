package main

import (
	"fmt"
	"io/ioutil"
	"sort"
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
	heightmap := parse(input)
	sum := 0
	for j, row := range heightmap.rows {
		for i, n := range row {
			if heightmap.isLowest(i, j) {
				sum += n + 1
			}
		}
	}
	fmt.Printf("%d\n", sum)
}

func run2(input string) {
	heightmap := parse(input)
	lows := []Pos{}
	for j, row := range heightmap.rows {
		for i := range row {
			if heightmap.isLowest(i, j) {
				lows = append(lows, Pos{i, j})
			}
		}
	}
	basins := []Basin{}
	for _, pos := range lows {
		if !includes(pos, &basins) {
			basins = append(basins, heightmap.findBasin(pos))
		}
	}

	sizes := []int{}
	for _, basin := range basins {
		sizes = append(sizes, len(basin.poses))
	}
	sort.Ints(sizes)

	prod := 1
	for _, l := range sizes[len(sizes)-3:] {
		prod *= l
	}

	fmt.Println(prod)
}

type Heightmap struct {
	rows []Row
}

type Row []int

type Pos struct {
	x, y int
}
type Basin struct {
	poses []Pos
}

func parse(input string) Heightmap {
	lines := strings.Split(input, "\n")
	rows := []Row{}
	for _, line := range lines {
		if 0 < len(line) {
			row := []int{}
			for _, x := range line {
				i, _ := strconv.Atoi(string(x))
				row = append(row, i)
			}
			rows = append(rows, Row(row))
		}
	}
	return Heightmap{rows}
}

func (heightmap *Heightmap) at(p Pos) int {
	return heightmap.rows[p.y][p.x]
}

func (heightmap *Heightmap) isLowest(i, j int) bool {
	p0 := Pos{i, j}
	n := heightmap.at(Pos{i, j})
	poses := []Pos{
		Pos{p0.x - 1, p0.y},
		Pos{p0.x + 1, p0.y},
		Pos{p0.x, p0.y - 1},
		Pos{p0.x, p0.y + 1}}
	for _, p := range poses {
		if heightmap.isInside(p) && heightmap.at(p) <= n {
			return false
		}
	}
	return true
}

func (heightmap *Heightmap) findBasin(p Pos) Basin {
	basin := Basin{[]Pos{}}
	heightmap.explor(p, &basin)
	return basin
}

func (heightmap *Heightmap) explor(p Pos, basin *Basin) {
	if heightmap.isInside(p) && !basin.contains(p) && heightmap.rows[p.y][p.x] < 9 {
		basin.poses = append(basin.poses, p)
		heightmap.explor(Pos{p.x - 1, p.y}, basin)
		heightmap.explor(Pos{p.x + 1, p.y}, basin)
		heightmap.explor(Pos{p.x, p.y - 1}, basin)
		heightmap.explor(Pos{p.x, p.y + 1}, basin)
	}
}

func (heightmap *Heightmap) isInside(p Pos) bool {
	return 0 <= p.x && p.x < len(heightmap.rows[0]) &&
		0 <= p.y && p.y < len(heightmap.rows)
}

func (basin *Basin) contains(p Pos) bool {
	for _, q := range basin.poses {
		if p == q {
			return true
		}
	}
	return false
}

func includes(p Pos, basins *[]Basin) bool {
	for _, basin := range *basins {
		if basin.contains(p) {
			return true
		}
	}
	return false
}
