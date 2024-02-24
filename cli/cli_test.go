// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"reflect"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	type flags struct {
		output string
		input  string
	}
	var f flags
	var cmd string
	var cmdArgs []string

	makeCmd := Command{
		Name:  "make",
		Short: "create uimage",
		Run: func(args []string) {
			cmdArgs = args
			cmd = "make"
		},
	}
	makeCmd.Flags().StringVar(&f.output, "o", "", "Output")
	makeCmd.Flags().StringVar(&f.input, "i", "", "Input")

	listCmd := Command{
		Name:    "list",
		Short:   "list uimage",
		Aliases: []string{"ls", "l"},
		Run: func(args []string) {
			cmdArgs = args
			cmd = "list"
		},
	}
	app := App{makeCmd, listCmd}

	for _, tt := range []struct {
		name        string
		args        []string
		wantCmd     string
		wantCmdArgs []string
		wantFlags   flags
		wantExit    int
		wantPrint   string
	}{
		{
			name:        "cmd with flag",
			args:        []string{"uimage", "make", "-o", "high", "foobar", "bla"},
			wantCmdArgs: []string{"foobar", "bla"},
			wantCmd:     "make",
			wantFlags: flags{
				output: "high",
			},
			wantExit:  0,
			wantPrint: "",
		},
		{
			name:      "not exist",
			args:      []string{"uimage", "notmake", "-o", "low"},
			wantExit:  1,
			wantPrint: "Commands:\n\n    make    create uimage\n    list    list uimage\n",
		},
		{
			name:      "cmd exist but flag doesn't",
			args:      []string{"uimage", "list", "-o", "low"},
			wantExit:  1,
			wantPrint: "flag provided but not defined: -o\nUsage of list:\n",
		},
		{
			name:        "cmd with no flags",
			args:        []string{"uimage", "list", "anything"},
			wantExit:    0,
			wantCmd:     "list",
			wantCmdArgs: []string{"anything"},
		},
		{
			name:        "alias",
			args:        []string{"uimage", "ls", "anything"},
			wantExit:    0,
			wantCmd:     "list",
			wantCmdArgs: []string{"anything"},
		},
		{
			name:      "no program name",
			wantExit:  1,
			wantPrint: "No program name provided\n",
		},
		{
			name:      "no command name",
			args:      []string{"uimage"},
			wantExit:  1,
			wantPrint: "Commands:\n\n    make    create uimage\n    list    list uimage\n",
		},
		{
			name:      "cmd help",
			args:      []string{"uimage", "make", "-h"},
			wantExit:  1,
			wantPrint: "Usage of make:\n  -i string\n    \tInput\n  -o string\n    \tOutput\n",
		},
		{
			name:      "app help",
			args:      []string{"uimage", "-h"},
			wantExit:  1,
			wantPrint: "Commands:\n\n    make    create uimage\n    list    list uimage\n",
		},
		{
			name:      "app help 2",
			args:      []string{"uimage", "help"},
			wantExit:  1,
			wantPrint: "Commands:\n\n    make    create uimage\n    list    list uimage\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f = flags{}
			cmd = ""
			cmdArgs = nil

			var s strings.Builder
			exitCode := app.run(&s, tt.args)
			t.Logf("App:\n%s", s.String())
			if exitCode != tt.wantExit {
				t.Errorf("run = %d, want %d", exitCode, tt.wantExit)
			}
			if cmd != tt.wantCmd {
				t.Errorf("run = cmd %s, want cmd %s", cmd, tt.wantCmd)
			}
			if !reflect.DeepEqual(cmdArgs, tt.wantCmdArgs) {
				t.Errorf("run = args %+v, want %+v", cmdArgs, tt.wantCmdArgs)
			}
			if !reflect.DeepEqual(f, tt.wantFlags) {
				t.Errorf("run = flags %+v, want %+v", f, tt.wantFlags)
			}
			if got := s.String(); got != tt.wantPrint {
				t.Errorf("run = %#v, want %#v", got, tt.wantPrint)
			}
		})
	}
}
