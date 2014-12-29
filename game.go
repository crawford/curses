package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/nsf/termbox-go"
)

type Point struct {
	x, y int
	age int
}

func (p Point) Draw() {
	color := termbox.Attribute(23 - p.age)
	termbox.SetCell(p.x * 2, p.y, ' ', termbox.ColorDefault, color)
	termbox.SetCell(p.x * 2 + 1, p.y, ' ', termbox.ColorDefault, color)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	running := true
	go func() {
		<-c
		running = false
	}()

	points := []*Point{}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	//width, height := termbox.Size()
	//width /= 2
	termbox.SetOutputMode(termbox.OutputGrayscale)
	termbox.SetInputMode(termbox.InputMouse | termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	go gameloop(&points)

	for running == true {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventMouse:
			points = append(points, &Point{
				x: ev.MouseX/2,
				y: ev.MouseY,
			})
			drawPoints(points)
		default:
			running = false
		}
	}
	termbox.Close()
}

func gameloop(points *[]*Point) {
	c := time.Tick(100 * time.Millisecond)
	for _ = range c {
		agePoints(points)
		np := []*Point{}
		for _, p := range *points {
			if p.y > 0 && p.age == 1 {
				np = append(np, &Point{
					x: p.x,
					y: p.y-1,
				})
			}
		}
		*points = append(*points, np...)

		drawPoints(*points)
	}
}

func agePoints(points *[]*Point) {
	np := *points
	*points = []*Point{}
	for _, p := range np {
		if p.age < 23 {
			p.age++
			*points = append(*points, p)
		}
	}
}

func drawPoints(points []*Point) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, p := range points {
		p.Draw()
	}
	termbox.Flush()
}
