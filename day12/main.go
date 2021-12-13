package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	input, _ := ioutil.ReadFile("input.txt")
	// input, _ := ioutil.ReadFile("input-1.txt")
	// input, _ := ioutil.ReadFile("input-2.txt")
	run1(string(input))
	run2(string(input))
}

func run1(input string) {
	system := parse(input)
	// system.display()
	routes := system.explore()
	fmt.Println(len(routes))
}

func run2(input string) {
	system := parse(input)
	// system.display()
	routes := system.explore2()
	fmt.Println(len(routes))
}

type Cave struct {
	label string
	small bool
	dsts  []string
}

type System struct {
	caves map[string]Cave
}

type Route struct {
	paths []string
}

func parse(input string) System {
	lines := strings.Split(input, "\n")
	caves := map[string]Cave{}
	for _, line := range lines {
		if 0 < len(line) {
			xs := strings.Split(line, "-")
			src := xs[0]
			dst := xs[1]
			for i := 0; i < 2; i++ {
				if _, b := caves[src]; !b && src != "end" {
					caves[src] = Cave{src, src == strings.ToLower(src), []string{}}
				}
				if cave, b := caves[src]; b && dst != "start" {
					cave.dsts = append(cave.dsts, dst)
					caves[src] = cave
				}
				src, dst = dst, src
			}
		}
	}
	return System{caves}
}

func (system *System) display() {
	for src, cave := range system.caves {
		if cave.small {
			fmt.Printf("%s (small)\n", src)
		} else {
			fmt.Printf("%s (big)\n", src)
		}
		for _, dst := range cave.dsts {
			fmt.Printf("  -> %s\n", dst)
		}
	}
}

func (system *System) explore() (routes []Route) {
	route := Route{[]string{"start"}}
	return system.step(route)
}

func (system *System) step(route Route) (res []Route) {
	src := route.paths[len(route.paths)-1]
	if src == "end" {
		res = append(res, route)
		return
	}
	cave := system.caves[src]
	for _, dst := range cave.dsts {
		if !system.caves[dst].small || !route.includes(dst) {
			next := Route{append(route.paths, dst)}
			for _, next_next := range system.step(next) {
				res = append(res, next_next)
			}
		}
	}
	return
}

func (system *System) explore2() (routes []Route) {
	route := Route{[]string{"start"}}
	return system.step2(route)
}

func (system *System) step2(route Route) (res []Route) {
	src := route.paths[len(route.paths)-1]
	if src == "end" {
		res = []Route{route}
		return
	}
	cave, b := system.caves[src]
	if !b {
		return
	}
	counts := map[string]int{}
	for _, path := range route.paths[1:] {
		if system.caves[path].small {
			counts[path] = route.count(path)
		}
	}
	n := 0
	for _, m := range counts {
		if 2 < m {
			return
		}
		if 2 == m {
			n++
		}
		if 1 < n {
			return
		}
	}
	for _, dst := range cave.dsts {
		paths := make([]string, len(route.paths))
		copy(paths, route.paths)
		nexts := system.step2(Route{append(paths, dst)})
		res = concat(res, nexts)
	}
	return
}

func (route *Route) includes(x string) bool {
	for _, path := range route.paths {
		if path == x {
			return true
		}
	}
	return false
}

func (route *Route) count(x string) (res int) {
	for _, path := range route.paths {
		if path == x {
			res++
		}
	}
	return
}

func concat(xs, ys []Route) []Route {
	res := make([]Route, len(xs))
	copy(res, xs)
	for _, y := range ys {
		res = append(res, y)
	}
	return res
}
