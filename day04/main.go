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
	calls, boards := parse(input)
	// fmt.Println(calls)
	// fmt.Println(boards)
	game := Game{calls, boards}
	var winning Board
	var call int
	for {
		call = step(&game)
		done, i := isDone(&game)
		if done {
			winning = game.boards[i]
			break
		}
	}
	// fmt.Println(game)
	// fmt.Println(winning)
	// fmt.Println(call)
	xs := unmarkedNumbers(&winning)
	fmt.Println(sum(xs) * call)
}

func run2(input string) {
	calls, boards := parse(input)
	// fmt.Println(calls)
	// fmt.Println(boards)
	game := Game{calls, boards}
	var winning Board
	var call int
	for {
		call = step(&game)
		done, i := isDone(&game)
		if done {
			winning = game.boards[i]
			if 1 == len(game.boards) {
				break
			}
			boards := []Board{}
			for _, board := range game.boards {
				if !isWinning(&board) {
					boards = append(boards, board)
				}
			}
			game.boards = boards
		}
	}
	// fmt.Println(len(game.boards))
	// fmt.Println(winning)
	// fmt.Println(call)
	xs := unmarkedNumbers(&winning)
	fmt.Println(sum(xs) * call)
}

func parse(input string) ([]int, []Board) {
	lines := strings.Split(input, "\n")
	first := lines[0]
	calls := []int{}
	boards := []Board{}
	for _, x := range strings.Split(first, ",") {
		i, _ := strconv.Atoi(x)
		calls = append(calls, i)
	}

	board := Board{}
	for i, line := range lines {
		if 0 == len(line) {
			if 0 < len(board.cells) {
				boards = append(boards, board)
				board = Board{}
			}
		} else if 0 < i {
			cells := []Cell{}
			for _, x := range strings.Split(line, " ") {
				if 0 < len(x) {
					i, _ := strconv.Atoi(x)
					cells = append(cells, Cell{i, false})
				}
			}
			board.cells = append(board.cells, cells)
		}
	}
	return calls, boards
}

type Cell struct {
	n   int
	hit bool
}

type Board struct {
	cells [][]Cell
}

type Game struct {
	calls  []int
	boards []Board
}

func step(game *Game) int {
	call := game.calls[0]
	game.calls = game.calls[1:]
	for _, board := range game.boards {
		mark(&board, call)
	}
	return call
}

func mark(board *Board, call int) {
	for j, row := range board.cells {
		for i, cell := range row {
			if cell.n == call {
				board.cells[j][i] = Cell{call, true}
			}
		}
	}
}

func isDone(game *Game) (bool, int) {
	for i, board := range game.boards {
		if isWinning(&board) {
			return true, i
		}
	}
	return false, 0
}

type Pos struct {
	x int
	y int
}

func isWinning(board *Board) bool {
	poses := []Pos{}
	for j, row := range board.cells {
		poses = []Pos{}
		for i, _ := range row {
			poses = append(poses, Pos{i, j})
		}
		if isAllTrue(poses, board) {
			return true
		}
		poses = []Pos{}
		for i, _ := range row {
			poses = append(poses, Pos{j, i})
		}
		if isAllTrue(poses, board) {
			return true
		}
	}
	poses = []Pos{}
	for j := range board.cells {
		poses = append(poses, Pos{j, j})
	}
	if isAllTrue(poses, board) {
		return true
	}
	return false
}

func isAllTrue(poses []Pos, board *Board) bool {
	for _, pos := range poses {
		if !board.cells[pos.y][pos.x].hit {
			return false
		}
	}
	return true
}

func unmarkedNumbers(board *Board) []int {
	res := []int{}
	for _, row := range board.cells {
		for _, cell := range row {
			if !cell.hit {
				res = append(res, cell.n)
			}
		}
	}
	return res
}

func sum(xs []int) int {
	res := 0
	for _, x := range xs {
		res += x
	}
	return res
}
