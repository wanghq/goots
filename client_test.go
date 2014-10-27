// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// test client for ots2
package goots

import (
	"testing"
)

func Test_New(t *testing.T) {
	o := New("http://127.0.0.1:8800", "OTSMultiUser177_accessid", "OTSMultiUser177_accesskey", "TestInstance177")
	t.Log(o)
	t.Fail()
}
