// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cli provides a bare bones CLI with commands.
//
// cli supports standard Go flags.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

// Command is a CLI command.
type Command struct {
	// Name or Aliases are used to match argv[1] and select this command.
	Name    string
	Aliases []string

	// Short is a one-line description for the command.
	Short string

	// Run is called if argv[1] matches the command. Flags are parsed before Run is called.
	Run func(args []string)

	flags *flag.FlagSet
}

// Flags returns a modifiable flag set for this command.
//
// If argv[1] matches this command, these flags will be parsed.
func (c *Command) Flags() *flag.FlagSet {
	if c.flags == nil {
		c.flags = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	}
	return c.flags
}

// An App is composed of many commands.
type App []Command

// Help returns the app's help string.
func (a App) Help() string {
	var s strings.Builder
	w := tabwriter.NewWriter(&s, 1, 2, 4, ' ', 0)
	fmt.Fprintf(w, "Commands:\n\n")
	for _, cmd := range a {
		fmt.Fprintf(w, "\t%s\t%s\n", cmd.Name, cmd.Short)
	}
	w.Flush()
	return s.String()
}

func (a App) commandFor(args []string) *Command {
	if len(args) == 0 {
		return nil
	}
	for _, cmd := range a {
		if args[0] == cmd.Name || slices.Contains(cmd.Aliases, args[0]) {
			return &cmd
		}
	}
	return nil
}

func (a App) run(errW io.Writer, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(errW, "No program name provided\n")
		return 1
	}
	cmd := a.commandFor(args[1:])
	if cmd == nil {
		fmt.Fprintf(errW, "%s", a.Help())
		return 1
	}

	cmd.Flags().SetOutput(errW)
	if err := cmd.Flags().Parse(args[2:]); err != nil {
		cmd.Flags().Output()
		return 1
	}

	cmd.Run(cmd.Flags().Args())
	return 0
}

// Run runs the app. Expects program name as the first arg, and an optional command name next.
func (a App) Run(args []string) {
	os.Exit(a.run(os.Stderr, args))
}
