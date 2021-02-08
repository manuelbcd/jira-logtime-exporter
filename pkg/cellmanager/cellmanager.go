package cellmanager

import "fmt"

type Cell struct {
	row 		int
	col 		string
	colIndex	int
}

const LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (c * Cell) Init() {
	c.colIndex = 0
	c.col = string(LETTERS[c.colIndex])
	c.row = 1
}

func (c * Cell) InitCol() {
	c.colIndex = 0
	c.col = string(LETTERS[c.colIndex])
}

func (c * Cell) IncRow()  {
	c.row ++
}

func (c * Cell) IncCol() * Cell {
	c.colIndex ++
	if c.colIndex >= len(LETTERS){
		c.colIndex = 0
	}
	c.col = string(LETTERS[c.colIndex])
	return c
}

func (c * Cell) GetStr() string {
	return fmt.Sprintf("%s%d", c.col, c.row)
}