package main

import (
	"os"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/urfave/cli"
)

var evalCommand = cli.Command{
	Name:    "eval",
	Aliases: []string{"e"},
	Usage:   "evaluate an arrai expression",
	Action:  eval,
}

func eval(c *cli.Context) error {
	source := c.Args().Get(0)

	expr, err := syntax.Parse(syntax.NewStringLexer(source))
	if err != nil {
		return err
	}

	global := &rel.Scope{}
	value, err := expr.Eval(global, global)
	if err != nil {
		return err
	}

	s := value.String()
	os.Stdout.WriteString(s)
	if s[len(s)-1] != '\n' {
		os.Stdout.Write([]byte{'\n'})
	}

	return nil
}
