// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llog_test

import (
	"flag"
	"log/slog"

	"github.com/u-root/uio/llog"
)

func someFunc(v llog.Printf) {
	v("logs at the given level %s", "foo")
}

func Example() {
	l := llog.Default()
	// If -v is set, l.Level becomes slog.LevelDebug.
	l.RegisterDebugFlag(flag.CommandLine, "v")
	flag.Parse()

	someFunc(l.Debugf)
	someFunc(l.AtLevelFunc(slog.LevelWarn))
	someFunc(l.Infof)
}
