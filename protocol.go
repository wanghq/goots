// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// protocol for ots2
package goots

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	// "encoding/hex"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"sort"
	"strings"
	// "sync"
	"time"

	"code.google.com/p/goprotobuf/proto"
	. "github.com/GiterLab/goots/log"
	. "github.com/GiterLab/goots/otstype"
	. "github.com/GiterLab/goots/protobuf"
	"github.com/GiterLab/goots/protobuf/coder"
)

var API_VERSION = "2014-08-08"
var defaultProtocol = ots_protocol{
	api_version: API_VERSION,
	encoder:     coder.EncodeRequest,
	decoder:     coder.DecodeRequest,
}

func newProtocol(protocol *ots_protocol) *ots_protocol {
	if protocol == nil {
		return &defaultProtocol
	}

	protocol = new(ots_protocol)
	protocol.api_version = API_VERSION
	protocol.encoder = coder.EncodeRequest
	protocol.decoder = coder.DecodeRequest

	return protocol
}

var api_list = DictString{
	"CreateTable":   "",
	"ListTable":     "",
	"DeleteTable":   "",
	"DescribeTable": "",
	"UpdateTable":   "",
	"GetRow":        "",
	"PutRow":        "",
	"UpdateRow":     "",
	"DeleteRow":     "",
	"BatchGetRow":   "",
	"BatchWriteRow": "",
	"GetRange":      "",
}

type ots_protocol struct {
	api_version   string
	user_id       string
	user_key      string
	instance_name string
	encoding      string
	encoder       func(api_name string, args ...interface{}) (req []reflect.Value, err error)
	decoder       func(api_name string, args ...interface{}) (req []reflect.Value, err error)
	logger        string
}

func (o *ots_protocol) Set(user_id, user_key, instance_name, encoding, logger string) *ots_protocol {
	if user_id != "" {
		o.user_id = user_id
	}

	if user_key != "" {
		o.user_key = user_key
	}

	if instance_name != "" {
		o.instance_name = instance_name
	}

	if encoding != "" {
		o.encoding = encoding
	}

	if logger != "" {
		o.logger = logger
	}

	return o
}

func (o *ots_protocol) _make_headers_string(headers DictString) string {
	if len(headers) == 0 {
		return "\n"
	}

	// count headers
	count := 0
	for k, _ := range headers {
		if strings.HasPrefix(strings.ToLower(k), "x-ots-") && strings.ToLower(k) != "x-ots-signature" {
			count++
		}
	}
	if count == 0 {
		return "\n"
	}

	strslice := make([]string, count)
	i := 0
	for k, v := range headers {
		if strings.HasPrefix(strings.ToLower(k), "x-ots-") && strings.ToLower(k) != "x-ots-signature" {
			strslice[i] = fmt.Sprintf("%s:%s", strings.ToLower(k), strings.TrimSpace(v.(string)))
			i++
		}
	}
	sort.Strings(strslice)

	return strings.Join(strslice, "\n")
}

func (o *ots_protocol) _call_signature_method(signature_string string) string {
	// The signature method is supposed to be HmacSHA1
	// A switch case is required if there is other methods available
	signature := hmacSha1(o.user_key, []byte(signature_string))
	return base64Encode(signature)
}

func (o *ots_protocol) _make_request_signature(query string, headers DictString) (signature string, err error) {
	url_obj, err := url.Parse(query)
	if err != nil {
		return "", err
	}

	// TODO a special query should be input to test query sorting,
	// because none of the current APIs uses query map, but the sorting
	// is required in the protocol document.
	uri := url_obj.Path
	query_string := url_obj.Query()
	sorted_query := ""
	if len(query_string) != 0 {
		strslice := make([]string, len(query_string))
		i := 0
		for k, v := range query_string {
			strslice[i] = fmt.Sprintf("%s:%s", k, v[0])
			i++
		}
		sort.Strings(strslice)
		sorted_query = strings.Join(strslice, "&")
	}
	sorted_query = urlencode(sorted_query)
	signature_string := uri + "\n" + "POST" + "\n" + sorted_query + "\n"

	headers_string := o._make_headers_string(headers)
	signature_string = signature_string + headers_string + "\n"
	signature = o._call_signature_method(signature_string)

	return signature, nil
}

func (o *ots_protocol) _make_headers(body []byte, query string) (headers DictString, err error) {
	// compose request headers and process request body if needed
	md5 := base64Encode(md5Encode(body))

	// rfc822
	// "Tue, 12 Aug 2014 10:23:03 GMT"
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	headers = DictString{
		"x-ots-date":         date,
		"x-ots-apiversion":   o.api_version,
		"x-ots-accesskeyid":  o.user_id,
		"x-ots-instancename": o.instance_name,
		"x-ots-contentmd5":   md5,
	}

	signature, err := o._make_request_signature(query, headers)
	if err != nil {
		return nil, err
	}
	headers["x-ots-signature"] = signature

	return headers, nil
}

func (o *ots_protocol) _make_response_signature(query string, headers DictString) (signature string, err error) {
	url_obj, err := url.Parse(query)
	if err != nil {
		return "", err
	}

	uri := url_obj.Path
	headers_string := o._make_headers_string(headers)
	signature_string := headers_string + "\n" + uri
	signature = o._call_signature_method(signature_string)

	return signature, nil
}

func (o *ots_protocol) _check_headers(headers DictString, body []byte) (ok bool, err error) {
	// check the response headers and process response body if needed.

	// 1, make sure we have all headers
	header_names := []string{
		"x-ots-contentmd5",
		"x-ots-requestid",
		"x-ots-date",
		"x-ots-contenttype",
	}

	for _, name := range header_names {
		if _, ok := headers[name]; !ok {
			return false, (OTSClientError{}.Set("\"%s\" is missing in response header", name))
		}
	}

	// 2, check md5
	md5 := base64Encode(md5Encode(body))
	if md5 != headers["x-ots-contentmd5"] {
		return false, (OTSClientError{}.Set("MD5 mismatch in response"))
	}

	// 3. check date
	server_time, err := time.Parse(("Mon, 02 Jan 2006 15:04:05 GMT"), headers["x-ots-date"].(string))
	if err != nil {
		return false, (OTSClientError{}.Set("Invalid date format in response - %s", err))
	}

	// 4, check date range
	server_unix_time := server_time.UTC()
	now_unix_time := time.Now().UTC()

	d := now_unix_time.Sub(server_unix_time)
	if math.Abs(float64(d.Seconds())) > float64(15*60*time.Second) {
		return false, (OTSClientError{}.Set("The difference between date in response and system time is more than 15 minutes"))
	}

	return true, nil
}

func (o *ots_protocol) _check_authorization(query string, headers DictString) (ok bool, err error) {
	auth, ok := headers["authorization"]
	if !ok {
		auth, ok = headers["Authorization"]
		if !ok {
			return false, (OTSClientError{}.Set("\"Authorization\" is missing in response header"))
		}
	}

	// 1, check authorization
	if !strings.HasPrefix(auth.(string), "OTS ") {
		return false, (OTSClientError{}.Set("Invalid Authorization in response"))
	}

	// 2, check accessid
	auth_string := auth.(string)[4:]
	auth_slice := strings.Split(auth_string, ":")
	if len(auth_slice) != 2 {
		return false, (OTSClientError{}.Set("Invalid Authorization in response"))
	}
	access_id := auth_slice[0]
	signature := auth_slice[1]
	if access_id != o.user_id {
		return false, (OTSClientError{}.Set("Invalid accesskeyid in response"))
	}

	// 3, check signature
	signature_src, err := o._make_response_signature(query, headers)
	if err != nil {
		return false, (OTSClientError{}.Set("Invalid signature in response - %s", err))
	}
	if signature != signature_src {
		return false, (OTSClientError{}.Set("Invalid signature in response - %s but %s", signature, signature_src))
	}

	return true, nil
}

func (o *ots_protocol) make_request(api_name string, args ...interface{}) (query string, headers DictString, body []byte, err error) {
	if _, ok := api_list[api_name]; !ok {
		return "", DictString{}, nil, (OTSClientError{}.Set("API %s is not supported", api_name))
	}

	proto_obj, err := o.encoder(api_name, args...)
	if err != nil {
		return "", DictString{}, nil, err
	}

	if len(proto_obj) < 2 {
		return "", DictString{}, nil, (OTSClientError{}.Set("Not enough params"))
	} else {
		err_index := len(proto_obj)
		if proto_obj[err_index-1].Interface() != nil {
			err, ok := proto_obj[err_index-1].Interface().(error)
			if ok {
				return "", DictString{}, nil, err
			}
		}
	}

	// TODO:
	// MAKR BY TOBYZXJ
	// "CreateTable"
	// "ListTable"
	// "DeleteTable"
	// "DescribeTable"
	// "UpdateTable"
	// "GetRow"
	// "PutRow"
	// "UpdateRow"
	// "DeleteRow"
	// "BatchGetRow"
	// "BatchWriteRow"
	// "GetRange"
	switch t := proto_obj[0].Interface().(type) {
	case *CreateTableRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*CreateTableRequest)))
	case *ListTableRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*ListTableRequest)))
	case *DeleteTableRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*DeleteTableRequest)))
	case *DescribeTableRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*DescribeTableRequest)))
	case *UpdateTableRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*UpdateTableRequest)))
	case *GetRowRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*GetRowRequest)))
	case *PutRowRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*PutRowRequest)))
	case *UpdateRowRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*UpdateRowRequest)))
	case *DeleteRowRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*DeleteRowRequest)))
	case *BatchGetRowRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*BatchGetRowRequest)))
	case *BatchWriteRowRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*BatchWriteRowRequest)))
	case *GetRangeRequest:
		body = []byte(proto.MarshalTextString(proto_obj[0].Interface().(*GetRangeRequest)))

	default:
		return "", DictString{}, nil, fmt.Errorf("Unknown type: %v", t)
	}

	query = "/" + api_name
	headers, err = o._make_headers(body, query)
	if err != nil {
		return "", DictString{}, nil, err
	}
	if OTSDebugEnable {
		// prevent MessageToString from happening
		// when no log is going to be actually printed
		// since it's very time consuming
		OTSError{}.Set("OTS request, API: %s, Protobuf: %v", api_name, proto_obj[0].Interface())
	}

	return query, headers, body, nil
}

func (o *ots_protocol) _get_request_id_string(headers DictString) string {
	request_id, ok := headers["x-ots-requestid"]
	if ok {
		return request_id.(string)
	}

	return ""
}

func (o *ots_protocol) parse_response(api_name, status string, headers DictString, body []byte) (ok bool, err error) {

	return true, nil
}

func (o *ots_protocol) handle_error(api_name, query, status, reason string, headers DictString, body []byte) OTSError {

	return OTSError{}
}

///////////////////////////////////////
////       COMMON TOOLS            ////
///////////////////////////////////////

// create md5 string
func md5Encode(src []byte) []byte {
	h := md5.New()
	h.Write(src)
	return h.Sum(nil)
}

// hmacsha1
func hmacSha1(key string, src []byte) []byte {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write(src)
	return mac.Sum(nil)
}

// base64 encode
func base64Encode(src []byte) string {
	coder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=")
	return coder.EncodeToString(src)
}

// base64 decode
func base64Decode(src []byte) ([]byte, error) {
	coder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=")
	return coder.DecodeString(string(src))
}

// urlencode
func urlencode(s string) (result string) {
	for _, c := range s {
		if c <= 0x7f { // single byte
			result += fmt.Sprintf("%%%X", c)
		} else if c > 0x1fffff { // quaternary byte
			result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
				0xf0+((c&0x1c0000)>>18),
				0x80+((c&0x3f000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else if c > 0x7ff { // triple byte
			result += fmt.Sprintf("%%%X%%%X%%%X",
				0xe0+((c&0xf000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else { // double byte
			result += fmt.Sprintf("%%%X%%%X",
				0xc0+((c&0x7c0)>>6),
				0x80+(c&0x3f),
			)
		}
	}

	return result
}
