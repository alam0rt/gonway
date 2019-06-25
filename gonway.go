package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

/*
   The number of live neighbours is calculated for each square
   All the cells having less than two live neighbours die of solitude
   All the cells having more than three live neighbours die of overpopulation
   A new cell is born on each empty square that has three live neighbours.
*/

type point struct {
	x, y int
}

type Cell struct {
	living   bool
	location point
}

// to create a cell, we need to give it a location and
// we let it be alive
func (u Universe) spawnCell(p point, l bool) {
	u[p].location = p
	u[p].living = l
}

func (p point) GetLocation() point { return p }

func (p *point) PrintLocation() {
	fmt.Printf("x: %v, y: %v\n", p.x, p.y)
}

func (p *point) GetNeighbours() [8]point {
	var n [8]point
	n[0] = point{x: p.x - 1, y: p.y - 1} //  5 6 7  Neighbours "look" like this
	n[1] = point{x: p.x - 1, y: p.y}     //  3 p 4
	n[2] = point{x: p.x - 1, y: p.y + 1} //  0 1 2
	n[3] = point{x: p.x, y: p.y - 1}
	n[4] = point{x: p.x, y: p.y + 1}
	n[5] = point{x: p.x + 1, y: p.y - 1}
	n[6] = point{x: p.x + 1, y: p.y}
	n[7] = point{x: p.x + 1, y: p.y + 1}
	return n
}

const Width int = 80
const Height int = 24
const Tick int = 1

func printGrid(u Universe) {
	h := 0
	w := 0
	chars := Width * Height
	for h < Height {
		for c := 0; c < chars; c++ {
			if c%Width == 0 {
				h += 1 // move down a row
				w = 0  // and back to the start
				fmt.Printf("\n")
			}
			s := "."
			p := point{x: w, y: h}
			// if the cell actually exists
			if u[p] != nil {
				// and if it's alive
				if u[p].living == true {
					// write an indication
					s = "o"
				}
			}

			fmt.Printf("%s", s)
			w += 1
		}
	}
}

type Universe map[point]*Cell

func (u Universe) livingNeighbours(c *Cell) int {
	x := 0
	for _, v := range c.location.GetNeighbours() {
		if u[v] != nil {
			if u[v].living == true {
				x += 1
			}
		}
	}
	return x
}

// Every tick, the reaper cometh
func (u Universe) reaper() {
	for _, v := range u {
		n := u.livingNeighbours(v)
		if n < 2 && v.living == true {
			// die from starvation
			u.spawnCell(v.location, false)
		} else if (n == 2 || n == 3) && v.living == true {
			// live on
			u.spawnCell(v.location, true)
		} else if n > 3 && v.living == true {
			// die from overpopulation
			u.spawnCell(v.location, false)
		} else if n == 3 && v.living == false {
			// replicate
			u.spawnCell(v.location, true)
		}
	}

}

func main() {
	u := make(Universe)

	// initiate the grid with "dead" cells
	chars := Width * Height
	w := 0
	for h := 0; h < Height; {
		for c := 0; c < chars; c++ {
			if c%Width == 0 {
				h += 1 // move down a row
				w = 0  // and back to the start
			}
			p := point{x: w, y: h}
			u[p] = &Cell{}
			u.spawnCell(p, false)
			w += 1
		}
	}

	//   o o
	// o o
	//   o
	// create a "living" cell

	// center
	p := point{x: Width / 2, y: Height / 2}
	u[p] = &Cell{}
	u.spawnCell(p, true)

	// right
	p = point{x: Width/2 + 1, y: Height / 2}
	u[p] = &Cell{}
	u.spawnCell(p, true)

	p = point{x: Width / 2, y: Height/2 + 1}
	u[p] = &Cell{}
	u.spawnCell(p, true)

	p = point{x: Width/2 - 1, y: Height/2 + 1}
	u[p] = &Cell{}
	u.spawnCell(p, true)

	p = point{x: Width / 2, y: Height/2 + 2}
	u[p] = &Cell{}
	u.spawnCell(p, true)

	for {
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
		printGrid(u)
		fmt.Printf("Point %v has %v neigh\n", p, u.livingNeighbours(u[p]))
		time.Sleep(time.Second / 16)
		u.reaper()
	}
}
