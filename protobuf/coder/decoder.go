// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// encoder for ots2
package coder

import (
	// "fmt"
	"reflect"

	// "code.google.com/p/goprotobuf/proto"
	. "github.com/GiterLab/goots/log"
	// . "github.com/GiterLab/goots/otstype"
	// . "github.com/GiterLab/goots/protobuf"
)

var api_decode_map = NewFuncmap()

// request encode for ots2
func DecodeRequest(api_name string, args ...interface{}) (req []reflect.Value, err error) {
	if _, ok := api_decode_map[api_name]; !ok {
		return nil, (OTSClientError{}.Set("No PB encode method for API %s", api_name))
	}

	req, err = api_decode_map.Call(api_name, args...)
	if err != nil {
		return nil, err
	}

	return req, nil
}
