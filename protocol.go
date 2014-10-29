// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// protocol for ots2
package goots

import (
	"fmt"
	"strings"
	// "net/url"
	// "reflect"
	// "sync"
	// "time"

	// . "github.com/GiterLab/goots/log"
	// . "github.com/GiterLab/goots/otstype"
	// "github.com/GiterLab/goots/urllib"
	// "code.google.com/p/goprotobuf/proto"
)

var API_VERSION = "2014-08-08"
var protocol = ots_protocol{
	api_version: API_VERSION,
}

type ots_protocol struct {
	api_version   string
	user_id       string
	user_key      string
	instance_name string
	encoding      string
	// encoder func()
	// decoder func()
	logger string
}

func (o *ots_protocol) Set(user_id, user_key, instance_name, encoding, logger string) *ots_protocol {
	o.user_id = user_id
	o.user_key = user_key
	o.instance_name = instance_name
	o.encoding = encoding
	o.logger = logger

	return o
}

func (o *ots_protocol) _make_headers_string(headers map[string]string) string {
	if len(headers) == 0 {
		return "\n"
	}

	strslice := make([]string, len(headers))
	i := 0
	for k, v := range headers {
		if strings.HasPrefix(k, "x-ots-") && k != "x-ots-signature" {
			strslice[i] = fmt.Sprintf("%s:%s", strings.ToLower(k), strings.TrimSpace(v))
			i++
		}
	}

	fmt.Println(strslice)
	return ""
}
