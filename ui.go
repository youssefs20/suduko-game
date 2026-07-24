package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// UI renders the game onto a tcell screen. It holds no game state.
type UI struct {
	screen tcell.Screen
}

const (
	gridX = 1 // left margin of the grid
	gridY = 1 // top margin of the grid
)

var (
	lineStyle     = tcell.StyleDefault
	givenStyle    = tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	playerStyle   = tcell.StyleDefault.Foreground(tcell.ColorAqua)
	conflictStyle = tcell.StyleDefault.Foreground(tcell.ColorRed)
	statusStyle   = tcell.StyleDefault.Foreground(tcell.ColorGray)
	winStyle      = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
)

// DrawMenu renders the difficulty-selection screen.
func (u *UI) DrawMenu() {
	u.screen.Clear()
	lines := []string{
		"S U D O K U",
		"",
		"Choose difficulty:",
		"",
		"  1  Easy",
		"  2  Medium",
		"  3  Hard",
		"",
		"q  quit",
	}
	for i, l := range lines {
		u.drawText(gridX+1, gridY+i, l, tcell.StyleDefault)
	}
	u.screen.Show()
}

// DrawGame renders the board, cursor, status line, and win banner.
func (u *UI) DrawGame(b *Board, curR, curC int, d Difficulty, won bool) {
	u.screen.Clear()

	// Separator lines: heavy every third line, light otherwise.
	for i := 0; i <= 9; i++ {
		u.drawText(gridX, gridY+i*2, sepLine(i), lineStyle)
	}

	conflicts := b.Conflicts()
	for r := 0; r < 9; r++ {
		y := gridY + r*2 + 1
		// Vertical separators: heavy at box boundaries, light between cells.
		for c := 0; c <= 9; c++ {
			ch := '│'
			if c%3 == 0 {
				ch = '┃'
			}
			u.screen.SetContent(gridX+c*4, y, ch, nil, lineStyle)
		}
		// Cell contents: " d " with the style covering the whole cell so
		// the cursor's reverse-video block is visible on empty cells too.
		for c := 0; c < 9; c++ {
			cell := b.Cell(r, c)
			style := playerStyle
			if cell.Given {
				style = givenStyle
			}
			if conflicts[[2]int{r, c}] {
				style = conflictStyle
			}
			if r == curR && c == curC {
				style = style.Reverse(true)
			}
			ch := ' '
			if cell.Value != 0 {
				ch = rune('0' + cell.Value)
			}
			x := gridX + c*4 + 1
			u.screen.SetContent(x, y, ' ', nil, style)
			u.screen.SetContent(x+1, y, ch, nil, style)
			u.screen.SetContent(x+2, y, ' ', nil, style)
		}
	}

	status := fmt.Sprintf("%s  |  arrows/hjkl move   1-9 place   0/del clear   n new   q quit", d)
	u.drawText(gridX, gridY+20, status, statusStyle)
	if won {
		u.drawText(gridX, gridY+22, "You win!  Press n for a new game or q to quit.", winStyle)
	}
	u.screen.Show()
}

// sepLine returns horizontal separator i (0 = top border, 9 = bottom).
func sepLine(i int) string {
	switch {
	case i == 0:
		return "┏━━━┯━━━┯━━━┳━━━┯━━━┯━━━┳━━━┯━━━┯━━━┓"
	case i == 9:
		return "┗━━━┷━━━┷━━━┻━━━┷━━━┷━━━┻━━━┷━━━┷━━━┛"
	case i%3 == 0:
		return "┣━━━┿━━━┿━━━╋━━━┿━━━┿━━━╋━━━┿━━━┿━━━┫"
	default:
		return "┠───┼───┼───╂───┼───┼───╂───┼───┼───┨"
	}
}

func (u *UI) drawText(x, y int, s string, style tcell.Style) {
	for i, r := range []rune(s) {
		u.screen.SetContent(x+i, y, r, nil, style)
	}
}
