package main

import "testing"

// solvedGrid is a known-valid completed Sudoku grid (the classic
// Wikipedia example), used wherever a full valid board is needed.
var solvedGrid = [9][9]int{
	{5, 3, 4, 6, 7, 8, 9, 1, 2},
	{6, 7, 2, 1, 9, 5, 3, 4, 8},
	{1, 9, 8, 3, 4, 2, 5, 6, 7},
	{8, 5, 9, 7, 6, 1, 4, 2, 3},
	{4, 2, 6, 8, 5, 3, 7, 9, 1},
	{7, 1, 3, 9, 2, 4, 8, 5, 6},
	{9, 6, 1, 5, 3, 7, 2, 8, 4},
	{2, 8, 7, 4, 1, 9, 6, 3, 5},
	{3, 4, 5, 2, 8, 6, 1, 7, 9},
}

func TestNewBoardMarksGivens(t *testing.T) {
	var grid [9][9]int
	grid[2][3] = 7
	b := NewBoard(grid)
	if !b.Cell(2, 3).Given {
		t.Error("nonzero cell should be a given")
	}
	if b.Cell(0, 0).Given {
		t.Error("empty cell should not be a given")
	}
}

func TestSetCellOnEmptyCell(t *testing.T) {
	b := NewBoard([9][9]int{})
	if !b.SetCell(4, 4, 7) {
		t.Fatal("SetCell rejected a legal write")
	}
	if got := b.Cell(4, 4).Value; got != 7 {
		t.Fatalf("value = %d, want 7", got)
	}
	if !b.SetCell(4, 4, 0) {
		t.Fatal("SetCell rejected clearing a player cell")
	}
	if got := b.Cell(4, 4).Value; got != 0 {
		t.Fatalf("value after clear = %d, want 0", got)
	}
}

func TestGivenCellsCannotBeModified(t *testing.T) {
	var grid [9][9]int
	grid[0][0] = 5
	b := NewBoard(grid)
	if b.SetCell(0, 0, 3) {
		t.Error("SetCell reported success on a given cell")
	}
	if got := b.Cell(0, 0).Value; got != 5 {
		t.Errorf("given cell changed to %d", got)
	}
}

func TestSetCellRejectsOutOfRangeValues(t *testing.T) {
	b := NewBoard([9][9]int{})
	if b.SetCell(0, 0, 10) || b.SetCell(0, 0, -1) {
		t.Error("SetCell accepted an out-of-range value")
	}
}

func TestRowConflict(t *testing.T) {
	b := NewBoard([9][9]int{})
	b.SetCell(0, 0, 5)
	b.SetCell(0, 8, 5)
	assertConflicts(t, b, [][2]int{{0, 0}, {0, 8}})
}

func TestColumnConflict(t *testing.T) {
	b := NewBoard([9][9]int{})
	b.SetCell(0, 4, 9)
	b.SetCell(7, 4, 9)
	assertConflicts(t, b, [][2]int{{0, 4}, {7, 4}})
}

func TestBoxConflict(t *testing.T) {
	b := NewBoard([9][9]int{})
	// (0,0) and (2,2) share a box but not a row or column.
	b.SetCell(0, 0, 3)
	b.SetCell(2, 2, 3)
	assertConflicts(t, b, [][2]int{{0, 0}, {2, 2}})
}

func TestNoFalseConflictsOnValidBoard(t *testing.T) {
	b := NewBoard(solvedGrid)
	if got := b.Conflicts(); len(got) != 0 {
		t.Errorf("valid board reported conflicts: %v", got)
	}
}

func TestWonOnCompleteValidBoard(t *testing.T) {
	if !NewBoard(solvedGrid).Won() {
		t.Error("complete valid board should be won")
	}
}

func TestNotWonWithEmptyCell(t *testing.T) {
	grid := solvedGrid
	grid[8][8] = 0
	if NewBoard(grid).Won() {
		t.Error("board with an empty cell should not be won")
	}
}

func TestNotWonWithConflict(t *testing.T) {
	grid := solvedGrid
	grid[0][0] = 3 // duplicates the 3 at (0,1): full board, one conflict
	if NewBoard(grid).Won() {
		t.Error("full board with a conflict should not be won")
	}
}

// assertConflicts checks Conflicts() returns exactly the given coordinates.
func assertConflicts(t *testing.T, b *Board, want [][2]int) {
	t.Helper()
	got := b.Conflicts()
	if len(got) != len(want) {
		t.Fatalf("got %d conflict cells (%v), want %d", len(got), got, len(want))
	}
	for _, w := range want {
		if !got[w] {
			t.Errorf("missing conflict at %v", w)
		}
	}
}
