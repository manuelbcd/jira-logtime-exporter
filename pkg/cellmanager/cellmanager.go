package cellmanager

import (
	"fmt"
	"strings"
)

type Cell struct {
	row      int
	col      string
	colIndex int
}

const LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

/**
Cell initialization (A1)
*/
func (c *Cell) Init() {
	c.colIndex = 0
	c.col = string(LETTERS[c.colIndex])
	c.row = 1
}

/**
Column initialization
*/
func (c *Cell) InitCol() {
	c.colIndex = 0
	c.col = string(LETTERS[c.colIndex])
}

/**
Increment row
*/
func (c *Cell) IncRow() {
	c.row++
}

/**
Increment column
*/
func (c *Cell) IncCol() *Cell {
	c.colIndex++
	if c.colIndex >= len(LETTERS) {
		c.colIndex = 0
	}
	c.col = string(LETTERS[c.colIndex])
	return c
}

/**
Set a specific column
*/
func (c *Cell) SetCol(column byte) *Cell {
	c.colIndex = strings.Index(LETTERS, string(column))
	c.col = string(column)
	return c
}

/**
Return cell position (col+row)
*/
func (c *Cell) GetStr() string {
	return fmt.Sprintf("%s%d", c.col, c.row)
}
