package main

import (
	"befunge/stack"
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
)

type BefungeMatrix struct {
	code   [][]int
	s      *stack.Stack
	pos    [2]int
	vector [2]int
}

func getLines(s string) []string {
	main := []string{}
	aggr := ""
	max := 0
	for _, c := range s {
		if c == '\n' {
			if max < len(aggr) {
				max = len(aggr)
			}
			main = append(main, aggr)
			aggr = ""
			continue
		}
		aggr += string(c)
	}
	if max < len(aggr) {
		max = len(aggr)
	}
	main = append(main, aggr)
	for i, v := range main {
		if len(v) < max {
			for j := 0; j <= max-len(v); j++ {
				main[i] += " "
			}
		}
	}
	return main
}

func createMatrix(s []string) [][]int {
	main := [][]int{}
	for _, line := range s {
		temp := []int{}
		for _, c := range line {
			temp = append(temp, int(c))
		}
		main = append(main, temp)
	}
	return main
}

func (Bm *BefungeMatrix) setRight() {
	Bm.vector = [2]int{0, 1}
}
func (Bm *BefungeMatrix) setLeft() {
	Bm.vector = [2]int{0, -1}
}
func (Bm *BefungeMatrix) setDown() {
	Bm.vector = [2]int{1, 0}
}
func (Bm *BefungeMatrix) setUp() {
	Bm.vector = [2]int{-1, 0}
}

func (Bm *BefungeMatrix) Next() {
	Bm.pos[0] += Bm.vector[0]
	Bm.pos[1] += Bm.vector[1]
	tmp := [2]int{Bm.pos[0], Bm.pos[1]}

	if Bm.pos[0] == -1 {
		// going up
		tmp[0] = len(Bm.code) - 1
	} else if Bm.pos[1] == -1 {
		// going left
		tmp[1] = len(Bm.code[0]) - 1
	} else if Bm.pos[0] == len(Bm.code) {
		// going down
		tmp[0] = 0
	} else if Bm.pos[1] == len(Bm.code[Bm.pos[0]]) {
		// bottom right going right
		tmp[1] = 0
	}
	Bm.pos = tmp
	// fmt.Println(Bm.s)
}

func runFromFile(fpath string) error {
	if path.Ext(fpath) == ".bf" {
		f, err := ioutil.ReadFile(fpath)
		if err != nil {
			return errors.New("an error occured while reading the file: " + err.Error())
		}
		s := stack.Stack{Items: []int{}}
		bm := BefungeMatrix{code: createMatrix(getLines(string(f))), s: &s, vector: [2]int{0, 1}, pos: [2]int{0, 0}}
		bm.RunCode()
	}
	return nil
}

func (Bm *BefungeMatrix) RunCode() {
	stringMode := false
	for Bm.code[Bm.pos[0]][Bm.pos[1]] != '@' {
		curr := Bm.code[Bm.pos[0]][Bm.pos[1]]
		if stringMode && curr == '"' {
			stringMode = !stringMode
			Bm.Next()
			continue
		}
		if stringMode {
			Bm.s.Push(curr)
			Bm.Next()
			continue
		}
		switch curr {
		case '+':
			// Addition +
			a := Bm.s.Pop()
			b := Bm.s.Pop()
			Bm.s.Push(a + b)
		case '-':
			// Subtraction -
			a := Bm.s.Pop()
			b := Bm.s.Pop()
			Bm.s.Push(b - a)
		case '*':
			// Multiplication *
			a := Bm.s.Pop()
			b := Bm.s.Pop()
			Bm.s.Push(a * b)
		case '/':
			// Division /
			a := Bm.s.Pop()
			b := Bm.s.Pop()
			Bm.s.Push(int(b / a))
		case '%':
			// Modulo %
			a := Bm.s.Pop()
			b := Bm.s.Pop()
			Bm.s.Push(b % a)
		case '!':
			// Logical Not !
			a := Bm.s.Pop()
			if a == 0 {
				Bm.s.Push(1)
			} else {
				Bm.s.Push(0)
			}
		case '`':
			// Greate Than `
			a := Bm.s.Pop()
			b := Bm.s.Pop()
			if b > a {
				Bm.s.Push(1)
			} else {
				Bm.s.Push(0)
			}
		case '>':
			// right >
			Bm.setRight()
		case '<':
			// left <
			Bm.setLeft()
		case '^':
			// up ^
			Bm.setUp()
		case 'v':
			// down v
			Bm.setDown()
		case '?':
			// random direction ?
			randNum := rand.Intn(4)
			switch randNum {
			case 0:
				// right >
				Bm.setRight()
			case 1:
				// left <
				Bm.setLeft()
			case 2:
				// up ^
				Bm.setUp()
			case 3:
				// down v
				Bm.setDown()
			}
		case '_':
			// Horizontal IF _
			n := Bm.s.Pop()
			if n == 0 {
				Bm.setRight()
			} else {
				Bm.setLeft()
			}
		case '|':
			// Vertical IF |
			n := Bm.s.Pop()
			if n == 0 {
				Bm.setDown()
			} else {
				Bm.setUp()
			}
		case '"':
			// StringMode "
			stringMode = !stringMode
		case ':':
			// Duplicate Push :
			Bm.s.Push(Bm.s.Peek())
		case '\\':
			// Swap Stack \
			m := Bm.s.Pop()
			n := Bm.s.Pop()
			Bm.s.Push(m)
			Bm.s.Push(n)
		case '$':
			// Discard Stack Top $
			Bm.s.Pop()
		case '.':
			// Output as Int .
			fmt.Printf("%d", Bm.s.Pop())
		case ',':
			// Output as ASCII Char ,
			fmt.Printf("%s", string(rune(Bm.s.Pop())))
		case '#':
			// Bridge Jump Cmd #
			Bm.Next()
			Bm.Next()
			continue
		case 'g':
			// Get call g
			y := Bm.s.Pop()
			x := Bm.s.Pop()
			if y >= 0 && x >= 0 && y < len(Bm.code) && x < len(Bm.code[y]) {
				Bm.s.Push(Bm.code[y][x])
			} else {
				Bm.s.Push(0)
			}
		case 'p':
			// Put call p
			y := Bm.s.Pop()
			x := Bm.s.Pop()
			v := Bm.s.Pop()
			if y >= 0 && x >= 0 && y < len(Bm.code) && x < len(Bm.code[y]) {
				Bm.code[y][x] = v
			}
		case '&':
			// Get Int Input &
			n := 0
			fmt.Scanf("%d", &n)
			Bm.s.Push(n)
		case '~':
			// Get Char Input ~
			reader := bufio.NewReader(os.Stdin)
			r, _, _ := reader.ReadRune()
			Bm.s.Push(int(r))
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// Push to Stack
			num, _ := strconv.Atoi(string(rune(curr)))
			Bm.s.Push(num)
		}
		Bm.Next()
		// if Bm.pos[0] == -1 || Bm.pos[1] == -1 {
		// 	fmt.Println("Stack", Bm.s)
		// 	fmt.Println("Pos", Bm.pos)
		// 	fmt.Println("Vector", Bm.vector)
		// 	fmt.Println("Cmd", Bm.code)
		// }
	}
}

func main() {
	// Befunge-97/98 have issues
	if os.Args[1] != "" {
		err := runFromFile(os.Args[1])
		if err != nil {
			panic(errors.New("an error occured while opening the given file path: " + err.Error()))
		}
	} else {
		panic(errors.New("an error occured while opening the file path: " + os.Args[1]))
	}
}
