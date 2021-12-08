package main

import (
	"fmt"
	"io/ioutil"
	"sort"
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
	entries := parse(input)
	sum := 0
	for _, entry := range entries {
		for _, sig := range entry.output {
			n := len(sig)
			if n == 2 || n == 4 || n == 3 || n == 7 {
				sum++
			}
		}
	}
	fmt.Printf("%d\n", sum)
}

func run2(input string) {
	entries := parse(input)
	sum := 0
	for _, entry := range entries {
		digits := entry.analyze()
		v := 0
		for _, x := range entry.output {
			// fmt.Println(x)
			v = v*10 + digits[normalize(x)]
		}
		sum += v
	}
	fmt.Println(sum)
}

type Entry struct {
	patterns []string
	output   []string
}

type Digits map[string]int

func parse(input string) []Entry {
	lines := strings.Split(input, "\n")
	res := []Entry{}
	for _, line := range lines {
		if 0 < len(line) {
			patterns := []string{}
			output := []string{}
			xs := strings.Split(line, "|")
			for _, x := range strings.Split(strings.TrimSpace(xs[0]), " ") {
				patterns = append(patterns, x)
			}
			for _, x := range strings.Split(strings.TrimSpace(xs[1]), " ") {
				output = append(output, x)
			}
			res = append(res, Entry{patterns, output})
		}
	}
	return res

}

func (entry *Entry) analyze() Digits {
	pats := entry.patterns
	pat := make([]string, 10)
	for _, k := range entry.patterns {
		switch len(k) {
		case 2:
			pat[1] = k
			pats = rem(k, pats)
		case 4:
			pat[4] = k
			pats = rem(k, pats)
		case 3:
			pat[7] = k
			pats = rem(k, pats)
		case 7:
			pat[8] = k
			pats = rem(k, pats)
		}
	}

	{
		p9 := pat[4] + pat[7]
		for _, k := range pats {
			if contains(k, p9) {
				pat[9] = k
				pats = rem(k, pats)
				break
			}
		}

	}
	{
		p0 := pat[1]
		for _, k := range pats {
			if len(k) == 6 && contains(k, p0) {
				pat[0] = k
				pats = rem(k, pats)
				break
			}
		}
	}
	{
		for _, k := range pats {
			if len(k) == 6 {
				pat[6] = k
				pats = rem(k, pats)
				break
			}
		}
	}
	{
		p2 := del(pat[0], pat[4])
		for _, k := range pats {
			if len(k) == 5 && contains(k, p2) {
				pat[2] = k
				pats = rem(k, pats)
				break
			}
		}
	}
	{
		p5 := del(pat[0], pat[2])
		for _, k := range pats {
			if len(k) == 5 && contains(k, p5) {
				pat[5] = k
				pats = rem(k, pats)
				break
			}
		}
	}
	{
		for _, k := range pats {
			pat[3] = k
		}
	}
	digits := Digits(map[string]int{})
	for i, v := range pat {
		digits[normalize(v)] = i
	}
	return digits
}

func del(xs, ys string) string {
	res := xs
	for _, y := range ys {
		res = strings.Replace(res, string(y), "", -1)
	}
	return res
}

func contains(xs, ys string) bool {
	for _, y := range ys {
		if !strings.Contains(xs, string(y)) {
			return false
		}
	}
	return true
}

func rem(s string, xs []string) []string {
	res := []string{}
	for _, x := range xs {
		if x != s {
			res = append(res, x)
		}
	}
	return res
}

func normalize(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}
