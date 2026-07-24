package main

import "math/rand"

// Difficulty selects how many clues a generated puzzle keeps.
type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	default:
		return "Hard"
	}
}

// Clues returns the target number of given clues for the difficulty.
func (d Difficulty) Clues() int {
	switch d {
	case Easy:
		return 40
	case Medium:
		return 32
	default:
		return 26
	}
}

// Generate produces a puzzle with exactly one solution at the given
// difficulty, plus the solved grid it was derived from.
func Generate(d Difficulty) (*Board, [9][9]int) {
	for {
		var full [9][9]int
		if !fill(&full, 0) {
			continue // practically unreachable; retry from scratch
		}
		puzzle := removeClues(full, d.Clues())
		return NewBoard(puzzle), full
	}
}

// fill completes grid from cell index pos (row-major, 0-80) using
// randomized backtracking.
func fill(grid *[9][9]int, pos int) bool {
	if pos == 81 {
		return true
	}
	r, c := pos/9, pos%9
	for _, i := range rand.Perm(9) {
		v := i + 1
		if canPlace(grid, r, c, v) {
			grid[r][c] = v
			if fill(grid, pos+1) {
				return true
			}
			grid[r][c] = 0
		}
	}
	return false
}

// canPlace reports whether v can go at (r, c) without colliding with an
// existing value in the same row, column, or box.
func canPlace(grid *[9][9]int, r, c, v int) bool {
	for i := 0; i < 9; i++ {
		if grid[r][i] == v || grid[i][c] == v {
			return false
		}
	}
	br, bc := r/3*3, c/3*3
	for i := br; i < br+3; i++ {
		for j := bc; j < bc+3; j++ {
			if grid[i][j] == v {
				return false
			}
		}
	}
	return true
}

// removeClues blanks cells in random order, keeping the solution unique,
// until target clues remain or no further removal is possible. If removal
// stalls above the target the puzzle is used anyway (slightly easier than
// requested).
func removeClues(grid [9][9]int, target int) [9][9]int {
	clues := 81
	for _, pos := range rand.Perm(81) {
		if clues == target {
			break
		}
		r, c := pos/9, pos%9
		saved := grid[r][c]
		grid[r][c] = 0
		if countSolutions(grid, 2) != 1 {
			grid[r][c] = saved
		} else {
			clues--
		}
	}
	return grid
}

// countSolutions counts the solutions of grid by backtracking, stopping
// early once limit is reached.
func countSolutions(grid [9][9]int, limit int) int {
	var solve func(pos int) int
	solve = func(pos int) int {
		for pos < 81 && grid[pos/9][pos%9] != 0 {
			pos++
		}
		if pos == 81 {
			return 1
		}
		r, c := pos/9, pos%9
		count := 0
		for v := 1; v <= 9; v++ {
			if canPlace(&grid, r, c, v) {
				grid[r][c] = v
				count += solve(pos + 1)
				grid[r][c] = 0
				if count >= limit {
					break
				}
			}
		}
		return count
	}
	return solve(0)
}
