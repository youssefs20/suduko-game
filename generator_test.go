package main

import "testing"

func TestFillProducesValidCompleteGrid(t *testing.T) {
	var grid [9][9]int
	if !fill(&grid, 0) {
		t.Fatal("fill failed on an empty grid")
	}
	if !NewBoard(grid).Won() {
		t.Fatal("filled grid is not a valid complete board")
	}
}

func TestCountSolutionsOnSolvedGrid(t *testing.T) {
	if got := countSolutions(solvedGrid, 2); got != 1 {
		t.Errorf("solved grid has %d solutions, want 1", got)
	}
}

func TestCountSolutionsDetectsAmbiguity(t *testing.T) {
	if got := countSolutions([9][9]int{}, 2); got != 2 {
		t.Errorf("empty grid reported %d solutions, want limit 2", got)
	}
}

func TestGenerate(t *testing.T) {
	for _, d := range []Difficulty{Easy, Medium, Hard} {
		t.Run(d.String(), func(t *testing.T) {
			board, full := Generate(d)

			if !NewBoard(full).Won() {
				t.Fatal("returned solved grid is not valid")
			}

			grid := gridOf(board)
			clues := 0
			for r := 0; r < 9; r++ {
				for c := 0; c < 9; c++ {
					cell := board.Cell(r, c)
					if (cell.Value != 0) != cell.Given {
						t.Fatalf("cell (%d,%d): Given flag inconsistent with value %d", r, c, cell.Value)
					}
					if cell.Value != 0 {
						clues++
						if cell.Value != full[r][c] {
							t.Fatalf("clue (%d,%d)=%d disagrees with solution %d", r, c, cell.Value, full[r][c])
						}
					}
				}
			}

			// Removal stops exactly at the target, or stalls slightly above
			// it when uniqueness blocks further removal.
			if clues < d.Clues() || clues > d.Clues()+10 {
				t.Errorf("clue count %d not near target %d", clues, d.Clues())
			}

			if got := countSolutions(grid, 2); got != 1 {
				t.Errorf("puzzle has %d solutions, want exactly 1", got)
			}

			solved, ok := solveGrid(grid)
			if !ok {
				t.Fatal("puzzle is unsolvable")
			}
			if solved != full {
				t.Error("solving the puzzle does not reproduce the original grid")
			}
		})
	}
}

// gridOf extracts the values of a board into a plain grid.
func gridOf(b *Board) [9][9]int {
	var g [9][9]int
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			g[r][c] = b.Cell(r, c).Value
		}
	}
	return g
}

// solveGrid solves a puzzle by deterministic backtracking, returning the
// completed grid and whether a solution exists.
func solveGrid(g [9][9]int) ([9][9]int, bool) {
	var solve func(pos int) bool
	solve = func(pos int) bool {
		for pos < 81 && g[pos/9][pos%9] != 0 {
			pos++
		}
		if pos == 81 {
			return true
		}
		r, c := pos/9, pos%9
		for v := 1; v <= 9; v++ {
			if canPlace(&g, r, c, v) {
				g[r][c] = v
				if solve(pos + 1) {
					return true
				}
				g[r][c] = 0
			}
		}
		return false
	}
	ok := solve(0)
	return g, ok
}
