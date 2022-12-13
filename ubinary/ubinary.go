// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ubinary provides a native endian binary.ByteOrder.
//
// Deprecated: use github.com/josharian/native instead.
package ubinary

import "github.com/josharian/native"

// NativeEndian is $GOARCH's implementation of byte order.
//
// Deprecated: use github.com/josharian/native.Endian. This package
// now just forwards to that one.
var NativeEndian = native.Endian
