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
	// input, _ := ioutil.ReadFile("input-2.txt")
	// input, _ := ioutil.ReadFile("input-3.txt")
	run1(string(input))
	run2(string(input))
}

func run1(input string) {
	// input = "8A004A801A8002F478"
	// input = "620080001611562C8802118E34"
	// input = "C0015000016115A2E0802F182340"
	// input = "A0016C880162017C3686B18A3D4780"
	bits := parse(input)
	// bits.display()
	packets, _ := bits.decode(1)
	sum := packets[0].versionSum()
	fmt.Println(sum)
}

func run2(input string) {
	// input = "C200B40A82"
	// input = "04005AC33890"
	// input = "880086C3E88112"
	// input = "CE00C43D881120"
	// input = "D8005AC2A8F0"
	// input = "F600BC2D8F"
	// input = "9C005AC2F8F0"
	// input = "9C0141080250320F1802104A08"
	bits := parse(input)
	// bits.display()
	packets, _ := bits.decode(1)
	res := packets[0].calc()
	fmt.Println(res)
}

type Bits []bool

func (bits Bits) display() {
	for _, b := range bits {
		if b {
			fmt.Print(1)
		} else {
			fmt.Print(0)
		}
	}
	fmt.Println("")
}

func parse(input string) Bits {
	lines := strings.Split(input, "\n")
	bits := []bool{}
	for _, line := range lines {
		if 0 < len(line) {
			for _, c := range line {
				u, _ := strconv.ParseUint(string(c), 16, 16)
				// fmt.Println(u)
				for i := 3; 0 <= i; i-- {
					// fmt.Print(u >> i & 1)
					bits = append(bits, 1 == u>>i&1)
				}
				// fmt.Println("")
			}

		}
	}
	return Bits(bits)
}

type Packet struct {
	version int
	ty      int
	lit     int
	subs    []Packet
}

func (bits Bits) decode(limit int) (res []Packet, i int) {
	for (limit == 0 || len(res) < limit) && i < len(bits) {
		// bits[i:].display()
		cur := Packet{}
		cur.version = toInt(bits[i : i+3])
		cur.ty = toInt(bits[i+3 : i+6])
		if cur.ty == 4 {
			lit, n := Bits(bits[i+6:]).decodeLit()
			cur.lit = lit
			i += 6 + n
		} else if cur.ty != 4 {
			if bits[i+6] {
				count := toInt(Bits(bits[i+7 : i+7+11]))
				subs, n := Bits(bits[i+18:]).decode(count)
				cur.subs = subs
				i += 18 + n
			} else {
				length := toInt(Bits(bits[i+7 : i+7+15]))
				cur.subs, _ = Bits(bits[i+22 : i+22+length]).decode(0)
				i += 22 + length
			}
		} else {
			break
		}

		res = append(res, cur)
	}
	// bits.display()
	// fmt.Println(res)
	return
}

func (bits Bits) decodeLit() (v, i int) {
	bs := []bool{}
	for {
		bs = append(bs, bits[i+1:i+5]...)
		if !bits[i] {
			break
		}
		i += 5
	}
	return toInt(bs), i + 5
}

func toInt(bits Bits) (res int) {
	for _, b := range bits {
		res = res << 1
		if b {
			res++
		}
	}
	return
}

func (packet *Packet) versionSum() (sum int) {
	sum += packet.version
	if packet.subs != nil {
		for _, p := range packet.subs {
			sum += p.versionSum()
		}
	}
	return
}

func (packet *Packet) calc() (ret int) {
	switch packet.ty {
	case 0:
		for _, p := range packet.subs {
			ret += p.calc()
		}
	case 1:
		ret = 1
		for _, p := range packet.subs {
			ret *= p.calc()
		}
	case 2:
		ret = packet.subs[0].calc()
		for _, p := range packet.subs[1:] {
			if v := p.calc(); v < ret {
				ret = v
			}
		}
	case 3:
		ret = packet.subs[0].calc()
		for _, p := range packet.subs[1:] {
			if v := p.calc(); ret < v {
				ret = v
			}
		}
	case 4:
		ret = packet.lit
	case 5:
		if packet.subs[0].calc() > packet.subs[1].calc() {
			ret = 1
		}
	case 6:
		if packet.subs[0].calc() < packet.subs[1].calc() {
			ret = 1
		}
	case 7:
		if packet.subs[0].calc() == packet.subs[1].calc() {
			ret = 1
		}
	}
	return
}
