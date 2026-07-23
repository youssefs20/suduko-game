package main

// Cell is one square of the board. Value 0 means empty. Given cells are
// puzzle clues and cannot be modified by the player.
type Cell struct {
	Value int
	Given bool
}

// Board is a 9x9 Sudoku board.
type Board struct {
	cells [9][9]Cell
}

// NewBoard builds a board from a grid of values; nonzero cells become givens.
func NewBoard(grid [9][9]int) *Board {
	b := &Board{}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			b.cells[r][c] = Cell{Value: grid[r][c], Given: grid[r][c] != 0}
		}
	}
	return b
}

// Cell returns the cell at (row, col).
func (b *Board) Cell(row, col int) Cell {
	return b.cells[row][col]
}

// SetCell writes v (1-9, or 0 to clear) into (row, col). It returns false
// and leaves the board unchanged if the cell is a given or v is out of range.
func (b *Board) SetCell(row, col, v int) bool {
	if v < 0 || v > 9 || b.cells[row][col].Given {
		return false
	}
	b.cells[row][col].Value = v
	return true
}

// Conflicts returns the coordinates of every filled cell whose value
// collides with another filled cell in the same row, column, or 3x3 box.
// Both cells of each collision are included.
func (b *Board) Conflicts() map[[2]int]bool {
	conflicts := map[[2]int]bool{}
	mark := func(r1, c1, r2, c2 int) {
		v := b.cells[r1][c1].Value
		if v != 0 && v == b.cells[r2][c2].Value {
			conflicts[[2]int{r1, c1}] = true
			conflicts[[2]int{r2, c2}] = true
		}
	}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			for c2 := c + 1; c2 < 9; c2++ {
				mark(r, c, r, c2) // row pairs
			}
			for r2 := r + 1; r2 < 9; r2++ {
				mark(r, c, r2, c) // column pairs
			}
		}
	}
	for br := 0; br < 9; br += 3 {
		for bc := 0; bc < 9; bc += 3 {
			var coords [][2]int
			for r := br; r < br+3; r++ {
				for c := bc; c < bc+3; c++ {
					coords = append(coords, [2]int{r, c})
				}
			}
			for i := 0; i < len(coords); i++ {
				for j := i + 1; j < len(coords); j++ {
					mark(coords[i][0], coords[i][1], coords[j][0], coords[j][1])
				}
			}
		}
	}
	return conflicts
}

// Won reports whether the board is completely filled with no conflicts.
func (b *Board) Won() bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b.cells[r][c].Value == 0 {
				return false
			}
		}
	}
	return len(b.Conflicts()) == 0
}
