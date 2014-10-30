// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// test protocol for ots2
package goots

import (
	"fmt"
	"testing"

	"code.google.com/p/goprotobuf/proto"
	. "github.com/GiterLab/goots/otstype"
	"github.com/GiterLab/goots/protobuf"
)

func Test_make_headers_string(t *testing.T) {
	fmt.Println("Test protocol start...")
	headers := DictString{
		"x-ots-date":         "Tue, 12 Aug 2014 10:23:03 GMT",
		"x-ots-apiversion":   "2014-08-08",
		"x-ots-accesskeyid":  "29j2NtzlUr8hjP8b",
		"x-ots-instancename": "naketest",
		"x-ots-contentmd5":   "1B2M2Y8AsgTpgAmY7PhCfg==",
		"x-ots-signature":    "testforx-ots-signature",
	}
	str := newProtocol(nil)._make_headers_string(headers)
	// fmt.Println(str)

	//
	// x-ots-accesskeyid:29j2NtzlUr8hjP8b
	// x-ots-apiversion:2014-08-08
	// x-ots-contentmd5:1B2M2Y8AsgTpgAmY7PhCfg==
	// x-ots-date:Tue, 12 Aug 2014 10:23:03 GMT
	// x-ots-instancename:naketest
	//
	if str != "x-ots-accesskeyid:29j2NtzlUr8hjP8b\nx-ots-apiversion:2014-08-08\nx-ots-contentmd5:1B2M2Y8AsgTpgAmY7PhCfg==\nx-ots-date:Tue, 12 Aug 2014 10:23:03 GMT\nx-ots-instancename:naketest" {
		t.Fail()
	}
}

func Test_call_signature_method(t *testing.T) {
	stringToSign := "/ListTable\nPOST\n\nx-ots-accesskeyid:29j2NtzlUr8hjP8b\nx-ots-apiversion:2014-08-08\nx-ots-contentmd5:1B2M2Y8AsgTpgAmY7PhCfg==\nx-ots-date:Tue, 12 Aug 2014 10:23:03 GMT\nx-ots-instancename:naketest\n"
	key := "8AKqXmNBkl85QK70cAOuH4bBd3gS0J"

	protocol := newProtocol(nil)
	protocol.Set("", key, "", "", "")
	str := protocol._call_signature_method(stringToSign)
	// fmt.Println(str)

	if str != "4xap392B7EBpN+RmlHgNowjoG1w=" {
		t.Fail()
	}
}

func Test_make_request_signature(t *testing.T) {
	headers := DictString{
		"x-ots-date":         "Tue, 12 Aug 2014 10:23:03 GMT",
		"x-ots-apiversion":   "2014-08-08",
		"x-ots-accesskeyid":  "29j2NtzlUr8hjP8b",
		"x-ots-instancename": "naketest",
		"x-ots-contentmd5":   "1B2M2Y8AsgTpgAmY7PhCfg==",
		"x-ots-signature":    "testforx-ots-signature",
	}
	key := "8AKqXmNBkl85QK70cAOuH4bBd3gS0J"
	protocol := newProtocol(nil)
	protocol.Set("", key, "", "", "")
	query := "/ListTable"
	signature, err := protocol._make_request_signature(query, headers)
	if err != nil {
		t.Fail()
	}
	// fmt.Println("signature:", signature)

	if signature != "4xap392B7EBpN+RmlHgNowjoG1w=" {
		t.Fail()
	}
}

func Test_make_headers(t *testing.T) {
	key := "8AKqXmNBkl85QK70cAOuH4bBd3gS0J"
	protocol := newProtocol(nil)
	protocol.Set("", key, "", "", "")
	query := "/ListTable"

	proto_list_table := new(protobuf.ListTableRequest)

	body := proto.MarshalTextString(proto_list_table)
	header, err := protocol._make_headers([]byte(body), query)
	if err != nil {
		t.Fail()
	}

	if header["x-ots-contentmd5"].(string) != "1B2M2Y8AsgTpgAmY7PhCfg==" {
		t.Fail()
	}
}
