// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llog_test

import (
	"log/slog"

	"github.com/u-root/uio/llog"
)

func ExampleDefault_withtime() {
	l := llog.Default()
	l.Infof("An INFO level string")
	l.Debugf("A DEBUG level that does not appear")

	l.Level = slog.LevelDebug
	l.Debugf("A DEBUG level that appears")
}
