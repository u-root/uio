// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llog

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"testing"
)

func TestLevelFlag(t *testing.T) {
	for _, tt := range []struct {
		args []string
		want slog.Level
	}{
		{
			args: []string{"-level=4"},
			want: slog.LevelWarn,
		},
		{
			args: []string{},
			want: slog.LevelInfo,
		},
	} {
		f := flag.NewFlagSet("", flag.ContinueOnError)

		v := &Logger{}
		v.RegisterLevelFlag(f, "level")
		_ = f.Parse(tt.args)

		if v.Level != tt.want {
			t.Errorf("Parse(%#v) = %v, want %v", tt.args, v, tt.want)
		}
	}

	for _, tt := range []struct {
		args []string
		want slog.Level
		err  error
	}{
		{
			args: []string{"-v"},
			want: slog.LevelWarn,
		},
		{
			args: []string{},
			want: slog.LevelInfo,
		},
		{
			args: []string{"-v=true"},
			want: slog.LevelWarn,
		},
		{
			args: []string{"-v=true", "-v=false"},
			want: slog.LevelWarn,
		},
		{
			args: []string{"-v=foobar"},
			want: slog.LevelInfo,
			err:  strconv.ErrSyntax,
		},
	} {
		f := flag.NewFlagSet("", flag.ContinueOnError)

		v := &Logger{}
		v.RegisterVerboseFlag(f, "v", slog.LevelWarn)
		// Parse doesn't use %w.
		if err := f.Parse(tt.args); err != tt.err && err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Errorf("Parse(%#v) = %v, want %v", tt.args, err, tt.err)
		}
		if v.Level != tt.want {
			t.Errorf("Parse(%#v) = %v, want %v", tt.args, v, tt.want)
		}
	}

	for _, tt := range []struct {
		args []string
		want slog.Level
		err  error
	}{
		{
			args: []string{"-v"},
			want: slog.LevelDebug,
		},
		{
			args: []string{},
			want: slog.LevelInfo,
		},
		{
			args: []string{"-v=true"},
			want: slog.LevelDebug,
		},
		{
			args: []string{"-v=true", "-v=false"},
			want: slog.LevelDebug,
		},
		{
			args: []string{"-v=foobar"},
			want: slog.LevelInfo,
			err:  strconv.ErrSyntax,
		},
	} {
		f := flag.NewFlagSet("", flag.ContinueOnError)

		v := &Logger{}
		v.RegisterDebugFlag(f, "v")
		// Parse doesn't use %w.
		if err := f.Parse(tt.args); err != tt.err && err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Errorf("Parse(%#v) = %v, want %v", tt.args, err, tt.err)
		}
		if v.Level != tt.want {
			t.Errorf("Parse(%#v) = %v, want %v", tt.args, v, tt.want)
		}
	}
}

func TestNilLogger(t *testing.T) {
	for _, l := range []*Logger{nil, {}} {
		// Test that none of this panics.
		l.AtLevelFunc(slog.LevelDebug)("nothing")
		l.AtLevel(slog.LevelDebug).Printf("nothing")
		l.Debugf("nothing")
		l.Infof("nothing")
		l.Warnf("nothing")
		l.Errorf("nothing")
		l.Logf(slog.LevelDebug, "nothing")
	}
}

func TestLog(t *testing.T) {
	var s strings.Builder
	l := New(slog.LevelDebug, func(format string, args ...any) {
		fmt.Fprintf(&s, format+"\n", args...)
	})

	l.AtLevelFunc(slog.LevelDebug)("nothing")
	l.AtLevel(slog.LevelInfo).Printf("nothing")
	l.Debugf("nothing")
	l.Infof("nothing")
	l.Warnf("nothing")
	l.Errorf("nothing")
	l.Logf(slog.LevelDebug, "nothing")

	want := `DEBUG nothing
INFO nothing
DEBUG nothing
INFO nothing
WARN nothing
ERROR nothing
DEBUG nothing
`
	if got := s.String(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDefaults(t *testing.T) {
	var s strings.Builder
	log.SetOutput(&s)
	log.SetFlags(0)

	l := Debug()
	l.Debugf("foobar")
	want := "DEBUG foobar\n"
	if got := s.String(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	l = Default()
	l.Debugf("bazzed")
	l.Infof("stuff")
	want = "DEBUG foobar\nINFO stuff\n"
	if got := s.String(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	l = Test(t)
	l.Debugf("more foobar")
}
