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
	MAX_AGE       = 20
	AGE_THRESHOLD = 10
)

func (g Grid) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for x, col := range g {
		for y, age := range col {
			if age < MAX_AGE {
				color := termbox.Attribute(MAX_AGE - age)
				termbox.SetCell(x*2, y, ' ', termbox.ColorDefault, color)
				termbox.SetCell(x*2+1, y, ' ', termbox.ColorDefault, color)
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

func (g Grid) Get(x, y int) int {
	if x < 0 || x >= g.Width() ||
		y < 0 || y >= g.Height() {
		return -1
	}
	return g[x][y]
}

func (g *Grid) Set(x, y, age int) {
	(*g)[x][y] = age
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
			col = append(col, MAX_AGE+1)
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
	c := time.Tick(50 * time.Millisecond)
	for _ = range c {
		random(grid)
		grid.Age()
		grid.Draw()
	}
}

func square(grid *Grid) {
	for x := 0; x < grid.Width(); x++ {
		for y := 0; y < grid.Height(); y++ {
			if grid.Get(x, y) != 1 {
				continue
			}

			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					if grid.Get(x+dx, y+dy) > MAX_AGE {
						grid.Set(x+dx, y+dy, 0)
					}
				}
			}
		}
	}
}

func diamond(grid *Grid) {
	for x := 0; x < grid.Width(); x++ {
		for y := 0; y < grid.Height(); y++ {
			if grid.Get(x, y) != 1 {
				continue
			}

			if grid.Get(x-1, y) > MAX_AGE {
				grid.Set(x-1, y, 0)
			}
			if grid.Get(x, y-1) > MAX_AGE {
				grid.Set(x, y-1, 0)
			}
			if grid.Get(x+1, y) > MAX_AGE {
				grid.Set(x+1, y, 0)
			}
			if grid.Get(x, y+1) > MAX_AGE {
				grid.Set(x, y+1, 0)
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
			if grid.Get(x, y) != 1 {
				continue
			}

			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					if grid.Get(x+dx, y+dy) > MAX_AGE-AGE_THRESHOLD && should() {
						grid.Set(x+dx, y+dy, 0)
					}
				}
			}
		}
	}
}
