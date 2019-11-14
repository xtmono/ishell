package ishell

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/flynn-archive/go-shlex"
)

type iHelper struct {
	cmd      *Cmd
	disabled func() bool
}

func (ic iHelper) Do(line []rune, pos int) (help []rune, lines int) {
	if ic.disabled != nil && ic.disabled() {
		return nil, 0
	}
	var words []string
	if w, err := shlex.Split(string(line)); err == nil {
		words = w
	} else {
		// fall back
		words = strings.Fields(string(line))
	}

	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)
	lines = 0
	prefix := ""
	if len(words) > 0 && pos > 0 && line[pos-1] != ' ' {
		prefix = words[len(words)-1]
		lines += ic.getHelpString(prefix, words[:len(words)-1], tw)
	} else {
		lines += ic.getHelpString(prefix, words, tw)
	}
	tw.Flush()

	if buf.Len() == 0 {
		return []rune("\nNo context help"), 1
	}
	return []rune(buf.String()), lines
}

func (ic iHelper) getHelpString(prefix string, w []string, tw *tabwriter.Writer) int {
	cmd, _ := ic.cmd.FindCmd(w)
	if cmd == nil {
		cmd = ic.cmd
	}
	cnt := 0
	for _, c := range cmd.children {
		if strings.HasPrefix(c.Name, prefix) {
			if c.LongHelp != "" {
				fmt.Fprintf(tw, "\n\t%s\t\t\t%s", c.Name, c.LongHelp)
			} else if c.Help != "" {
				fmt.Fprintf(tw, "\n\t%s\t\t\t%s", c.Name, c.Help)
			} else {
				fmt.Fprintf(tw, "\n\t%s", c.Name)
			}
			cnt++
		}
	}
	return cnt
}
