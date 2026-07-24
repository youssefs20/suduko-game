package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

type mode int

const (
	modeMenu mode = iota
	modePlay
	modeWon
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, "sudoku:", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintln(os.Stderr, "sudoku:", err)
		os.Exit(1)
	}
	defer screen.Fini()
	run(screen)
}

// run drives the game until the player quits. Modes: menu (difficulty
// selection) -> play -> won; 'n' returns to the menu from either.
func run(screen tcell.Screen) {
	ui := &UI{screen: screen}
	m := modeMenu
	var board *Board
	var diff Difficulty
	curR, curC := 0, 0

	draw := func() {
		if m == modeMenu {
			ui.DrawMenu()
		} else {
			ui.DrawGame(board, curR, curC, diff, m == modeWon)
		}
	}
	draw()

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			draw()
		case *tcell.EventKey:
			key := ev.Key()
			r := rune(0)
			if key == tcell.KeyRune {
				r = ev.Rune()
			}

			if key == tcell.KeyEscape || key == tcell.KeyCtrlC || r == 'q' {
				return
			}

			switch m {
			case modeMenu:
				switch r {
				case '1':
					diff = Easy
				case '2':
					diff = Medium
				case '3':
					diff = Hard
				default:
					continue
				}
				board, _ = Generate(diff)
				curR, curC = 0, 0
				m = modePlay
				draw()

			case modePlay, modeWon:
				if r == 'n' {
					m = modeMenu
					draw()
					continue
				}
				if m == modeWon {
					continue // board locked except n/q
				}
				switch {
				case key == tcell.KeyUp || r == 'k':
					if curR > 0 {
						curR--
					}
				case key == tcell.KeyDown || r == 'j':
					if curR < 8 {
						curR++
					}
				case key == tcell.KeyLeft || r == 'h':
					if curC > 0 {
						curC--
					}
				case key == tcell.KeyRight || r == 'l':
					if curC < 8 {
						curC++
					}
				case r >= '1' && r <= '9':
					board.SetCell(curR, curC, int(r-'0'))
					if board.Won() {
						m = modeWon
					}
				case key == tcell.KeyBackspace || key == tcell.KeyBackspace2 ||
					key == tcell.KeyDelete || r == '0':
					board.SetCell(curR, curC, 0)
				}
				draw()
			}
		}
	}
}
