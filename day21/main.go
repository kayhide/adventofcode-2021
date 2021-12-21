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
	players := parse(input)

	tick := 0
	dice := Dice100{1, 0}
	for !finished(1000, players) {
		n := tick % 2
		step := dice.roll() + dice.roll() + dice.roll()
		players[n].proceed(step)
		// fmt.Println(n, players[n])
		tick++
	}
	// fmt.Println(players)

	score := players[tick%2].score
	fmt.Println(score * dice.count())
}

func run2(input string) {
	players := parse(input)

	tick := 0
	us := Universes{Game{players[0], players[1]}: 1}
	for !us.done() {
		// fmt.Println(tick, us)
		us = us.step(tick)
		tick++
	}

	wins := map[int]int{0: 0, 1: 0}
	for g, n := range us {
		// fmt.Println(g, n, g.winner())
		wins[g.winner()] += n
	}

	max := 0
	for _, n := range wins {
		if max < n {
			max = n
		}
	}
	fmt.Println(max)
}

type Game struct {
	player1, player2 Player
}

func (g Game) winner() int {
	if 21 <= g.player1.score {
		return 0
	} else if 21 <= g.player2.score {
		return 1
	}
	return -1
}

func (g Game) step(tick int) (res []Game) {
	player := g.player1
	if tick%2 == 1 {
		player = g.player2
	}

	for k := 1; k <= 3; k++ {
		for j := 1; j <= 3; j++ {
			for i := 1; i <= 3; i++ {
				next := Player{player.at, player.score}
				next.proceed(i + j + k)
				if tick%2 == 0 {
					res = append(res, Game{next, g.player2})
				} else {
					res = append(res, Game{g.player1, next})
				}

			}
		}
	}
	return
}

type Universes map[Game]int

func (us *Universes) done() bool {
	for g := range *us {
		if -1 == g.winner() {
			return false
		}
	}
	return true
}

func (us *Universes) step(tick int) Universes {
	res := Universes{}
	for g, n := range *us {
		if -1 == g.winner() {
			for _, x := range g.step(tick) {
				if _, b := res[x]; !b {
					res[x] = 0
				}
				res[x] += n
			}
		} else {
			if _, b := res[g]; !b {
				res[g] = 0
			}
			res[g] += n

		}
	}
	return res
}

type Player struct {
	at    int
	score int
}

func (p *Player) proceed(n int) {
	p.at += n
	for 10 < p.at {
		p.at -= 10
	}
	p.score += p.at
}

func finished(n int, players []Player) bool {
	for _, p := range players {
		if n <= p.score {
			return true
		}
	}
	return false
}

type Dice interface {
	roll() int
	count() int
}

type Dice100 struct {
	next int
	n    int
}

func (d *Dice100) roll() int {
	res := d.next
	d.next++
	if 100 < d.next {
		d.next = 1
	}
	d.n++
	return res
}

func (d *Dice100) count() int {
	return d.n
}

func parse(input string) []Player {
	lines := strings.Split(input, "\n")
	players := []Player{}

	for _, line := range lines {
		if 0 < len(line) {
			x := strings.Split(line, ": ")[1]
			n, _ := strconv.Atoi(x)
			players = append(players, Player{n, 0})
		}
	}
	return players
}
