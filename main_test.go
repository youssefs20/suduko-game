package main

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

// runWithKeys starts run() on a simulation screen, injects the given key
// runes, and fails the test if run() does not return within the timeout.
func runWithKeys(t *testing.T, runes []rune) {
	t.Helper()
	s := tcell.NewSimulationScreen("")
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Fini()

	done := make(chan struct{})
	go func() {
		run(s)
		close(done)
	}()
	for _, r := range runes {
		s.InjectKey(tcell.KeyRune, r, tcell.ModNone)
	}
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("run did not exit after injected keys")
	}
}

func TestQuitFromMenu(t *testing.T) {
	runWithKeys(t, []rune{'q'})
}

func TestStartGameThenQuit(t *testing.T) {
	// '1' selects Easy (triggers generation), then a move and a quit.
	runWithKeys(t, []rune{'1', 'j', 'q'})
}
