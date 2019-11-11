// +build windows

package ishell

import (
	"github.com/xtmono/readline"
)

func clearScreen(s *Shell) error {
	return readline.ClearScreen(s.writer)
}
