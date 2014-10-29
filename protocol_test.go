// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// test protocol for ots2
package goots

import (
	// . "github.com/GiterLab/goots/otstype"
	"fmt"
	"testing"
)

func Test_make_headers_string(t *testing.T) {
	headers := map[string]string{
		"x-ots-date":         "Tue, 12 Aug 2014 10:23:03 GMT",
		"x-ots-apiversion":   "2014-08-08",
		"x-ots-accesskeyid":  "29j2NtzlUr8hjP8b",
		"x-ots-instancename": "naketest",
		"x-ots-contentmd5":   "1B2M2Y8AsgTpgAmY7PhCfg==",
		"x-ots-signature":    "testforx-ots-signature",
	}
	str := protocol._make_headers_string(headers)
	fmt.Println(str)
}
