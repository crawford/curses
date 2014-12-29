package main

import (
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/nsf/termbox-go"
)

type Grid [][]int

const (
	MAX_AGE=10
)

func (g Grid) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for x, col := range g {
		for y, age := range col {
			if age < MAX_AGE {
				color := termbox.Attribute(MAX_AGE - age)
				termbox.SetCell(x * 2, y, ' ', termbox.ColorDefault, color)
				termbox.SetCell(x * 2 + 1, y, ' ', termbox.ColorDefault, color)
			}
		}
	}
	termbox.Flush()
}

func (g *Grid) Age() {
	for x, col := range *g {
		for y, age := range col {
			if age <= MAX_AGE {
				(*g)[x][y]++
			}
		}
	}
}

func (g Grid) Width() int {
	return len(g)
}

func (g Grid) Height() int {
	if g.Width() == 0 {
		return 0
	}
	return len(g[0])
}

func NewGrid(width, height int) *Grid {
	grid := Grid{}
	for x := 0; x < width/2; x++ {
		col := []int{}
		for y := 0; y < height; y++ {
			col = append(col, MAX_AGE + 1)
		}
		grid = append(grid, col)
	}
	return &grid
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	running := true
	go func() {
		<-c
		running = false
	}()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	grid := NewGrid(termbox.Size())
	termbox.SetOutputMode(termbox.OutputGrayscale)
	termbox.SetInputMode(termbox.InputMouse | termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	go gameloop(grid)

	for running == true {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventMouse:
			(*grid)[ev.MouseX/2][ev.MouseY] = 0
			grid.Draw()
		default:
			running = false
		}
	}
	termbox.Close()
}

func gameloop(grid *Grid) {
	c := time.Tick(100 * time.Millisecond)
	for _ = range c {
		random(grid)
		grid.Age()
		grid.Draw()
	}
}

func square(grid *Grid) {
	for x := 0; x < grid.Width(); x++ {
		for y := 0; y < grid.Height(); y++ {
			if (*grid)[x][y] != 1 {
				continue
			}

			if x > 0 && (*grid)[x - 1][y] > MAX_AGE {
				(*grid)[x - 1][y] = 0
			}
			if x > 0 && y > 0 && (*grid)[x - 1][y - 1] > MAX_AGE {
				(*grid)[x - 1][y - 1] = 0
			}
			if y > 0 && (*grid)[x][y - 1] > MAX_AGE {
				(*grid)[x][y - 1] = 0
			}
			if x < grid.Width() - 1 && y > 0 && (*grid)[x + 1][y - 1] > MAX_AGE {
				(*grid)[x + 1][y - 1] = 0
			}
			if x < grid.Width() - 1 && (*grid)[x + 1][y] > MAX_AGE {
				(*grid)[x + 1][y] = 0
			}
			if x < grid.Width() - 1 && y < grid.Height() - 1 && (*grid)[x + 1][y + 1] > MAX_AGE {
				(*grid)[x + 1][y + 1] = 0
			}
			if y < grid.Height() - 1 && (*grid)[x][y + 1] > MAX_AGE {
				(*grid)[x][y + 1] = 0
			}
			if x > 0 && y < grid.Height() - 1 && (*grid)[x - 1][y + 1] > MAX_AGE {
				(*grid)[x - 1][y + 1] = 0
			}
		}
	}
}

func diamond(grid *Grid) {
	for x := 0; x < grid.Width(); x++ {
		for y := 0; y < grid.Height(); y++ {
			if (*grid)[x][y] != 1 {
				continue
			}

			if x > 0 && (*grid)[x - 1][y] > MAX_AGE {
				(*grid)[x - 1][y] = 0
			}
			if y > 0 && (*grid)[x][y - 1] > MAX_AGE {
				(*grid)[x][y - 1] = 0
			}
			if x < grid.Width() - 1 && (*grid)[x + 1][y] > MAX_AGE {
				(*grid)[x + 1][y] = 0
			}
			if y < grid.Height() - 1 && (*grid)[x][y + 1] > MAX_AGE {
				(*grid)[x][y + 1] = 0
			}
		}
	}
}

func random(grid *Grid) {
	should := func() bool {
		return rand.Float32() > 0.6
	}
	for x := 0; x < grid.Width(); x++ {
		for y := 0; y < grid.Height(); y++ {
			if (*grid)[x][y] != 1 {
				continue
			}

			if x > 0 && (*grid)[x - 1][y] > MAX_AGE && should() {
				(*grid)[x - 1][y] = 0
			}
			if x > 0 && y > 0 && (*grid)[x - 1][y - 1] > MAX_AGE && should() {
				(*grid)[x - 1][y - 1] = 0
			}
			if y > 0 && (*grid)[x][y - 1] > MAX_AGE && should() {
				(*grid)[x][y - 1] = 0
			}
			if x < grid.Width() - 1 && y > 0 && (*grid)[x + 1][y - 1] > MAX_AGE && should() {
				(*grid)[x + 1][y - 1] = 0
			}
			if x < grid.Width() - 1 && (*grid)[x + 1][y] > MAX_AGE && should() {
				(*grid)[x + 1][y] = 0
			}
			if x < grid.Width() - 1 && y < grid.Height() - 1 && (*grid)[x + 1][y + 1] > MAX_AGE && should() {
				(*grid)[x + 1][y + 1] = 0
			}
			if y < grid.Height() - 1 && (*grid)[x][y + 1] > MAX_AGE && should() {
				(*grid)[x][y + 1] = 0
			}
			if x > 0 && y < grid.Height() - 1 && (*grid)[x - 1][y + 1] > MAX_AGE && should() {
				(*grid)[x - 1][y + 1] = 0
			}
		}
	}
}
