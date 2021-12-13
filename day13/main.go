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
	paper := parse(input)
	// paper.display()
	paper.fold(0)
	fmt.Printf("%d\n", len(paper.poses))
}

func run2(input string) {
	paper := parse(input)
	// paper.display()
	for i := 0; i < len(paper.folds); i++ {
    	paper.fold(i)
	}
	paper.display()
}

type Pos struct {
	x, y int
}

type Fold struct {
	horizontal bool
	n          int
}

type Paper struct {
	poses []Pos
	folds []Fold
	end   Pos
}

func parse(input string) Paper {
	lines := strings.Split(input, "\n")
	poses := []Pos{}
	folds := []Fold{}
	for _, line := range lines {
		if 0 < len(line) {
			if strings.HasPrefix(line, "fold") {
				words := strings.Split(line, " ")
				lr := strings.Split(words[2], "=")
				n, _ := strconv.Atoi(lr[1])
				folds = append(folds, Fold{lr[0] == "y", n})
			} else {
				xy := strings.Split(line, ",")
				x, _ := strconv.Atoi(xy[0])
				y, _ := strconv.Atoi(xy[1])
				poses = append(poses, Pos{x, y})
			}
		}
	}
	end := Pos{0, 0}
	for _, p := range poses {
		if end.x < p.x {
			end.x = p.x
		}
		if end.y < p.y {
			end.y = p.y
		}
	}
	return Paper{poses, folds, end}
}

func (paper *Paper) display() {
	for j := 0; j <= paper.end.y; j++ {
		for i := 0; i <= paper.end.x; i++ {
			if includes(paper.poses, Pos{i, j}) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func (paper *Paper) fold(i int) {
	fold := paper.folds[i]
	poses := []Pos{}
	if fold.horizontal {
		for _, p := range paper.poses {
			if fold.n < p.y {
				p.y = 2*fold.n - p.y
			}
			poses = append(poses, p)
		}
		paper.end.y = fold.n
	} else {
		for _, p := range paper.poses {
			if fold.n < p.x {
				p.x = 2*fold.n - p.x
			}
			poses = append(poses, p)
		}
		paper.end.x = fold.n
	}

	paper.poses = distinct(poses)
}

func distinct(poses []Pos) (res []Pos) {
	for _, p := range poses {
		if !includes(res, p) {
			res = append(res, p)
		}
	}
	return
}

func includes(poses []Pos, p Pos) bool {
	for _, x := range poses {
		if x == p {
			return true
		}
	}
	return false
}
