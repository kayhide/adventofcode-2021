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
	vents := parse(input)
	// fmt.Println(vents)
	board := initBoard(vents)
	for _, vent := range vents {
		board = drawVent(board, vent)
	}
	count := 0
	for _, row := range board {
		for _, p := range row {
			if 2 <= p {
				count++
			}
		}
	}
	fmt.Println(count)
}

func run2(input string) {
	vents := parse(input)
	// fmt.Println(vents)
	board := initBoard(vents)
	for _, vent := range vents {
		board = drawVent(board, vent)
		board = drawVentDiagonal(board, vent)
	}
	count := 0
	for _, row := range board {
		for _, p := range row {
			if 2 <= p {
				count++
			}
		}
	}
	// for _, row := range board {
	//    	fmt.Println(row)
	// }
	fmt.Println(count)
}

type Point struct {
	x, y int
}
type Vent struct {
	pt0, pt1 Point
}

type Board [][]int

func parse(input string) []Vent {
	lines := strings.Split(input, "\n")
	res := []Vent{}
	for _, line := range lines {
		if 0 < len(line) {
			poses := strings.Split(line, " -> ")
			res = append(res, Vent{parsePoint(poses[0]), parsePoint(poses[1])})
		}
	}
	return res
}

func parsePoint(s string) Point {
	ns := strings.Split(s, ",")
	x, _ := strconv.Atoi(ns[0])
	y, _ := strconv.Atoi(ns[1])
	return Point{x, y}
}

func initBoard(vents []Vent) Board {
	m := Point{}
	for _, vent := range vents {
		x := max(vent.pt0.x, vent.pt1.x)
		if m.x < x {
			m.x = x
		}
		y := max(vent.pt0.y, vent.pt1.y)
		if m.y < y {
			m.y = y
		}
	}
	res := make([][]int, m.y+1)
	for j := range res {
		res[j] = make([]int, m.x+1)
	}
	return res
}

func min(x int, y int) int {
	if x > y {
		return y
	}
	return x
}

func max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

func drawVent(board Board, vent Vent) Board {
	pt0 := vent.pt0
	pt1 := vent.pt1
	if pt0.y == pt1.y {
		for i := min(pt0.x, pt1.x); i <= max(pt0.x, pt1.x); i++ {
			board[pt0.y][i]++
		}
	}
	if pt0.x == pt1.x {
		for i := min(pt0.y, pt1.y); i <= max(pt0.y, pt1.y); i++ {
			board[i][pt0.x]++
		}
	}
	return board
}

func drawVentDiagonal(board Board, vent Vent) Board {
	var pt0, pt1 Point
	if vent.pt0.x < vent.pt1.x {
		pt0 = vent.pt0
		pt1 = vent.pt1
	} else {
		pt0 = vent.pt1
		pt1 = vent.pt0
	}
	if (pt1.x - pt0.x) == (pt1.y - pt0.y) {
		for i := 0; i <= (pt1.x - pt0.x); i++ {
			board[pt0.y+i][pt0.x+i]++
		}
	}
	if (pt1.x - pt0.x) == -(pt1.y - pt0.y) {
		for i := 0; i <= (pt1.x - pt0.x); i++ {
			board[pt0.y-i][pt0.x+i]++
		}
	}
	return board
}
