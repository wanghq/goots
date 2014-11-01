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
	signature := protocol._call_signature_method(stringToSign)

	if signature != "4xap392B7EBpN+RmlHgNowjoG1w=" {
		t.Logf("signature shoud be %s, not %s", "4xap392B7EBpN+RmlHgNowjoG1w=", signature)
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
		t.Log(err)
		t.Fail()
	}
	// fmt.Println("signature:", signature)

	if signature != "4xap392B7EBpN+RmlHgNowjoG1w=" {
		t.Logf("signature shoud be %s, not %s", "4xap392B7EBpN+RmlHgNowjoG1w=", signature)
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
		t.Log(err)
		t.Fail()
	}

	if header["x-ots-contentmd5"].(string) != "1B2M2Y8AsgTpgAmY7PhCfg==" {
		t.Logf("x-ots-contentmd5 shoud be %s, not %s", "1B2M2Y8AsgTpgAmY7PhCfg=", header["x-ots-contentmd5"].(string))
		t.Fail()
	}
}

func Test_make_response_signature(t *testing.T) {
	key := "8AKqXmNBkl85QK70cAOuH4bBd3gS0J"
	protocol := newProtocol(nil)
	protocol.Set("", key, "", "", "")
	query := "/ListTable"

	headers := DictString{
		"x-ots-date":        "Tue, 12 Aug 2014 10:23:03 GMT",
		"x-ots-requestid":   "0005006c-0e81-db74-4a34-ce0a5df229a1",
		"x-ots-contenttype": "protocol buffer",
		"x-ots-contentmd5":  "1B2M2Y8AsgTpgAmY7PhCfg==",
	}

	signature, err := protocol._make_response_signature(query, headers)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if signature != "Y24MHhVti5UhSCW5qsUSDvT9SOk=" {
		t.Logf("signature shoud be %s, not %s", "Y24MHhVti5UhSCW5qsUSDvT9SOk=", signature)
		t.Fail()
	}
}

func Test_check_authorization(t *testing.T) {
	uid := "29j2NtzlUr8hjP8b"
	key := "8AKqXmNBkl85QK70cAOuH4bBd3gS0J"
	protocol := newProtocol(nil)
	protocol.Set(uid, key, "", "", "")
	query := "/ListTable"

	headers := DictString{
		"x-ots-date":        "Tue, 12 Aug 2014 10:23:03 GMT",
		"x-ots-requestid":   "0005006c-0e81-db74-4a34-ce0a5df229a1",
		"x-ots-contenttype": "protocol buffer",
		"x-ots-contentmd5":  "1B2M2Y8AsgTpgAmY7PhCfg==",
		"Authorization":     "OTS 29j2NtzlUr8hjP8b:Y24MHhVti5UhSCW5qsUSDvT9SOk=",
	}

	ok, err := protocol._check_authorization(query, headers)
	if err != nil {
		t.Fail()
	}
	if !ok {
		t.Fail()
	}
}

func Test_check_authorization_tobyzxj(t *testing.T) {
	uid := "0AkCEeXUWXeviDP6"
	key := "W8foPaZ53CB61C5H8JnTURsdXekWua"
	protocol := newProtocol(nil)
	protocol.Set(uid, key, "", "", "")
	query := "/ListTable"

	headers := DictString{
		"Date":              "Sat, 01 Nov 2014 08:49:24 GMT",
		"Connection":        "keep-alive",
		"Authorization":     "OTS 0AkCEeXUWXeviDP6:LQ1pIcPfC9NHZWMR9vrydXG4U4A=",
		"X-Ots-Contentmd5":  "IjpgaUwGKkfuEgyLaDq1mg==",
		"X-Ots-Contenttype": "protocol buffer",
		"X-Ots-Date":        "Sat, 01 Nov 2014 08:49:24 GMT",
		"X-Ots-requestid":   "000506c8-30ac-5974-1388-990a05bf1034",
	}

	ok, err := protocol._check_authorization(query, headers)
	if err != nil {
		t.Fail()
	}
	if !ok {
		t.Fail()
	}
}

func Test_make_request(t *testing.T) {
	uid := "29j2NtzlUr8hjP8b"
	key := "8AKqXmNBkl85QK70cAOuH4bBd3gS0J"
	protocol := newProtocol(nil)
	protocol.Set(uid, key, "", "", "")
	query := "/CreateTable"

	table_meta := OTSTableMeta{
		TableName: "myTable",
		SchemaOfPrimaryKey: OTSSchemaOfPrimaryKey{
			"gid": "INTEGER",
			"uid": "INTEGER",
		},
	}

	reserved_throughput := OTSReservedThroughput{
		OTSCapacityUnit{100, 100},
	}

	query, headers, body, err := protocol.make_request("CreateTable", &table_meta, &reserved_throughput)
	if err != nil {
		t.Fail()
	}

	t.Log("query:", query)
	t.Log("headers:", headers)
	t.Log("body:", body)

	// fmt.Println("query:", query)
	// fmt.Println("headers:", headers)
	// fmt.Println("body:", body)
}

func Test_get_request_id_string(t *testing.T) {
	headers := DictString{
		"x-ots-date":        "Tue, 12 Aug 2014 10:23:03 GMT",
		"x-ots-requestid":   "0005006c-0e81-db74-4a34-ce0a5df229a1",
		"x-ots-contenttype": "protocol buffer",
		"x-ots-contentmd5":  "1B2M2Y8AsgTpgAmY7PhCfg==",
		"Authorization":     "OTS 29j2NtzlUr8hjP8b:Y24MHhVti5UhSCW5qsUSDvT9SOk=",
	}

	protocol := newProtocol(nil)
	requestid := protocol._get_request_id_string(headers)
	if requestid != "0005006c-0e81-db74-4a34-ce0a5df229a1" {
		t.Fail()
	}
}
