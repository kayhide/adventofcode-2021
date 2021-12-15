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
	cavern := parse(input)
	// fmt.Println(cavern)
	ex := Explorer{map[Pos]int{}}
	ex.sums[Pos{0, 0}] = 0

	cavern.explore(&ex, cavern.end)
	cavern.review(&ex)
	fmt.Println(ex.sums[cavern.end])
}

func run2(input string) {
	cavern := parse(input)
	// cavern.display()
	cavern.enlarge(5)
	// cavern.display()

	ex := Explorer{map[Pos]int{}}
	ex.sums[Pos{0, 0}] = 0

	cavern.explore(&ex, cavern.end)
	cavern.review(&ex)
	fmt.Println(ex.sums[cavern.end])
}

type Pos struct {
	x, y int
}

type Cavern struct {
	risks [][]int
	end   Pos
}

type Explorer struct {
	sums map[Pos]int
}

func parse(input string) Cavern {
	lines := strings.Split(input, "\n")
	risks := [][]int{}
	for _, line := range lines {
		if 0 < len(line) {
			row := make([]int, len(line))
			for i, c := range line {
				v, _ := strconv.Atoi(string(c))
				row[i] = v
			}
			risks = append(risks, row)
		}
	}
	end := Pos{len(risks[0]) - 1, len(risks) - 1}
	return Cavern{risks, end}
}

func (cavern *Cavern) isInside(p Pos) bool {
	return 0 <= p.x && p.x <= cavern.end.x && 0 <= p.y && p.y <= cavern.end.y
}

func (cavern *Cavern) at(p Pos) int {
	return cavern.risks[p.y][p.x]
}

func (cavern *Cavern) explore(ex *Explorer, p Pos) int {
	if x, b := ex.sums[p]; b {
		return x
	}
	p1 := Pos{p.x - 1, p.y}
	p2 := Pos{p.x, p.y - 1}
	xs := []int{}
	if cavern.isInside(p1) {
		xs = append(xs, cavern.explore(ex, p1))
	}
	if cavern.isInside(p2) {
		xs = append(xs, cavern.explore(ex, p2))
	}

	min := xs[0]
	for _, x := range xs[1:] {
		if x < min {
			min = x
		}
	}
	ex.sums[p] = min + cavern.at(p)
	return ex.sums[p]
}

func (cavern *Cavern) poses() []Pos {
	poses := []Pos{}
	for j := 0; j <= cavern.end.y; j++ {
		for i := 0; i <= cavern.end.x; i++ {
			poses = append(poses, Pos{i, j})
		}
	}
	return poses

}
func (cavern *Cavern) review(ex *Explorer) {
	for {
		fmt.Print(".")
		poses := cavern.poses()
		next := []Pos{}
		for _, p := range poses {
			if b := cavern.reviewAt(ex, p); b {
				next = append(next, p)
			}
		}
		if len(next) == 0 {
			break
		}
		poses = next
	}
	fmt.Println("")
}

func (cavern *Cavern) reviewAt(ex *Explorer, p Pos) bool {
	if p.x == 0 && p.y == 0 {
		return false
	}
	p1 := Pos{p.x - 1, p.y}
	p2 := Pos{p.x + 1, p.y}
	p3 := Pos{p.x, p.y - 1}
	p4 := Pos{p.x, p.y + 1}
	poses := []Pos{}
	min := ex.sums[cavern.end]
	for _, p_ := range []Pos{p1, p2, p3, p4} {
		if cavern.isInside(p_) {
			if v := ex.sums[p_]; v < min {
				min = v
			} else {
				poses = append(poses, p_)
			}
		}
	}

	v := min + cavern.at(p)
	if x, b := ex.sums[p]; !b || x != v {
		ex.sums[p] = v
		for _, q := range poses {
			cavern.reviewAt(ex, q)
		}
		return true
	}
	return false
}

func (cavern *Cavern) display() {
	for j := 0; j <= cavern.end.y; j++ {
		for i := 0; i <= cavern.end.x; i++ {
			fmt.Printf("%2d", cavern.at(Pos{i, j}))
		}
		fmt.Println("")
	}
	fmt.Println("")
}
func (cavern *Cavern) show(ex *Explorer) {
	for j := 0; j <= cavern.end.y; j++ {
		for i := 0; i <= cavern.end.x; i++ {
			fmt.Printf("%2d", cavern.at(Pos{i, j}))
		}
		fmt.Println("")
	}
	fmt.Println("")
	// for j :=0; j <= cavern.end.y; j++ {
	//     for i := 0; i <= cavern.end.x; i++ {
	//         fmt.Printf("%3d", ex.sums[Pos{i,j}] / 10)
	//     }
	//     fmt.Println("")
	// }
	// fmt.Println("")
}

func (cavern *Cavern) enlarge(s int) {
	size := Pos{cavern.end.x + 1, cavern.end.y + 1}
	risks := make([][]int, size.y*s)
	for j := 0; j < size.y*s; j++ {
		risks[j] = make([]int, size.x*s)
		for i := 0; i < size.x*s; i++ {
			v := cavern.at(Pos{i % size.x, j % size.y})
			d := j/size.y + i/size.x
			v = (v+d-1)%9 + 1
			risks[j][i] = v
		}
	}
	cavern.risks = risks
	cavern.end = Pos{size.x*s - 1, size.y*s - 1}
}
