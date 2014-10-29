// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// test client for ots2
package goots

import (
	. "github.com/GiterLab/goots/otstype"
	"testing"
)

func Test_New(t *testing.T) {
	o, err := New("http://127.0.0.1:8800", "OTSMultiUser177_accessid", "OTSMultiUser177_accesskey", "TestInstance177")
	if err != nil {
		t.Fail()
	}
	t.Log(o)

	o, err = New("http://127.0.0.1:8800", "OTSMultiUser177_accessid", "OTSMultiUser177_accesskey", "TestInstance177",
		60, 60, "ots2-client-test", "utf-8")
	if err != nil {
		t.Fail()
	}
	t.Log(o)
	// t.Fail()
}

func Test_Set(t *testing.T) {
	o, err := New("http://127.0.0.1:8800", "OTSMultiUser177_accessid", "OTSMultiUser177_accesskey", "TestInstance177")
	if err != nil {
		t.Fail()
	}
	t.Log(o)

	o = o.Set(DictString{
		"EndPoint": "http://127.0.0.1:8888",
		// "NotExist": 123,
	})

	t.Log(o)
}
