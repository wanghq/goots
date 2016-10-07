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

	. "github.com/GiterLab/goots/log"
	. "github.com/GiterLab/goots/otstype"
	. "github.com/GiterLab/goots/protobuf"
	"github.com/GiterLab/goots/protobuf/coder"
	"github.com/golang/protobuf/proto"
)

var API_VERSION = "2014-08-08"
var defaultProtocol = ots_protocol{
	api_version: API_VERSION,
	encoder:     coder.EncodeRequest,
	decoder:     coder.DecodeRequest,
}

func newProtocol(protocol *ots_protocol) *ots_protocol {
	if OTSDebugEnable {
		coder.DebugEncoderEnable = true
		coder.DebugDecoderEnable = true
	}

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

func (o *ots_protocol) _check_headers(headers DictString, body []byte, status int) (ok bool, err error) {
	// check the response headers and process response body if needed.

	// 1, make sure we have all headers
	header_names := []string{
		"x-ots-contentmd5",
		"x-ots-requestid",
		"x-ots-date",
		"x-ots-contenttype",
	}

	if status >= 200 && status < 300 {
		for _, name := range header_names {
			if _, ok := headers[name]; !ok {
				return false, (OTSClientError{}.Set("\"%s\" is missing in response header", name))
			}
		}
	}

	// 2, check md5
	if _, ok := headers["x-ots-contentmd5"]; ok {
		md5 := base64Encode(md5Encode(body))
		if md5 != headers["x-ots-contentmd5"] {
			return false, (OTSClientError{}.Set("MD5 mismatch in response"))
		}
	}

	// 3. check date
	if _, ok := headers["x-ots-date"]; ok {
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

	pb, err := o.encoder(api_name, args...)
	if err != nil {
		return "", DictString{}, nil, err
	}

	if len(pb) < 2 {
		return "", DictString{}, nil, (OTSClientError{}.Set("Not enough params"))
	} else {
		err_index := len(pb)
		if pb[err_index-1].Interface() != nil {
			err, ok := pb[err_index-1].Interface().(error)
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
	switch t := pb[0].Interface().(type) {
	case *CreateTableRequest:
		body, err = proto.Marshal(pb[0].Interface().(*CreateTableRequest))
	case *ListTableRequest:
		body, err = proto.Marshal(pb[0].Interface().(*ListTableRequest))
	case *DeleteTableRequest:
		body, err = proto.Marshal(pb[0].Interface().(*DeleteTableRequest))
	case *DescribeTableRequest:
		body, err = proto.Marshal(pb[0].Interface().(*DescribeTableRequest))
	case *UpdateTableRequest:
		body, err = proto.Marshal(pb[0].Interface().(*UpdateTableRequest))
	case *GetRowRequest:
		body, err = proto.Marshal(pb[0].Interface().(*GetRowRequest))
	case *PutRowRequest:
		body, err = proto.Marshal(pb[0].Interface().(*PutRowRequest))
	case *UpdateRowRequest:
		body, err = proto.Marshal(pb[0].Interface().(*UpdateRowRequest))
	case *DeleteRowRequest:
		body, err = proto.Marshal(pb[0].Interface().(*DeleteRowRequest))
	case *BatchGetRowRequest:
		body, err = proto.Marshal(pb[0].Interface().(*BatchGetRowRequest))
	case *BatchWriteRowRequest:
		body, err = proto.Marshal(pb[0].Interface().(*BatchWriteRowRequest))
	case *GetRangeRequest:
		body, err = proto.Marshal(pb[0].Interface().(*GetRangeRequest))

	default:
		return "", DictString{}, nil, fmt.Errorf("Unknown type: %v", t)
	}

	query = "/" + api_name
	headers, err = o._make_headers(body, query)
	if err != nil {
		return "", DictString{}, nil, err
	}

	// prevent MessageToString from happening
	// when no log is going to be actually printed
	// since it's very time consuming
	OTSError{}.Log(OTSLoggerEnable, "OTS request, API: %s, Headers: %s, Protobuf: %v", api_name, fmt.Sprintf("%v", headers), pb[0].Interface())

	return query, headers, body, nil
}

func (o *ots_protocol) _get_request_id_string(headers DictString) string {
	request_id, ok := headers["x-ots-requestid"]
	if ok {
		return request_id.(string)
	}

	return ""
}

func (o *ots_protocol) parse_response(api_name, reason string, status int, headers DictString, body []byte) (ret []reflect.Value, ots_service_err *OTSServiceError) {
	ots_service_err = new(OTSServiceError)
	if _, ok := api_list[api_name]; !ok {
		return nil, (ots_service_err.SetErrorMessage("API %s is not supported", api_name))
	}

	ret, err := o.decoder(api_name, body)
	if err != nil {
		request_id := o._get_request_id_string(headers)
		error_message := fmt.Sprintf("Response format is invalid, %s, RequestID: %s, HTTP status: %s, Body: %v.", err, request_id, reason, body)
		return nil, ots_service_err.SetErrorMessage(error_message).SetHttpStatus(reason).SetRequestId(request_id).SetErrorCode(fmt.Sprintf("%d", status))
	}

	// prevent MessageToString from happening
	// when no log is going to be actually printed
	// since it's very time consuming
	request_id := o._get_request_id_string(headers)
	OTSError{}.Log(OTSLoggerEnable, "OTS request, API: %s, RequestID: %s, Protobuf: %v", api_name, request_id, ret[0].Interface())

	return ret, nil
}

func (o *ots_protocol) handle_error(api_name, query, reason string, status int, headers DictString, body []byte) (ots_service_err *OTSServiceError) {
	ots_service_err = new(OTSServiceError)
	request_id := o._get_request_id_string(headers)
	if _, ok := api_list[api_name]; !ok {
		return ots_service_err.SetErrorMessage("API %s is not supported", api_name).SetHttpStatus(reason).SetErrorCode(fmt.Sprintf("%d", status)).SetRequestId(request_id)
	}

	// 1. check headers & _check authorization
	if ok, err := o._check_headers(headers, body, status); !ok {
		return ots_service_err.SetErrorMessage("check headers failed - %s", err).SetHttpStatus(reason).SetErrorCode(fmt.Sprintf("%d", status)).SetRequestId(request_id)
	}
	if status != 403 {
		if ok, err := o._check_authorization(query, headers); !ok {
			return ots_service_err.SetErrorMessage("check authorization failed - %s", err).SetHttpStatus(reason).SetErrorCode(fmt.Sprintf("%d", status)).SetRequestId(request_id)
		}
	}

	// 2. ok
	if status >= 200 && status < 300 {
		return nil
	} else {
		// prevent MessageToString from happening
		// when no log is going to be actually printed
		// since it's very time consuming
		pb_err := &Error{}
		proto.Unmarshal(body, pb_err)
		if pb_err.Code != nil && pb_err.Message != nil {
			OTSError{}.Log(OTSLoggerEnable, "OTS request failed, API: %s, HTTPStatus: %s, ErrorCode: %s, ErrorMessage: %s,  RequestID: %s", api_name, reason, pb_err.GetCode(), pb_err.GetMessage(), request_id)
			return ots_service_err.SetErrorMessage(pb_err.GetMessage()).SetHttpStatus(reason).SetErrorCode(pb_err.GetCode()).SetRequestId(request_id)
		}

		OTSError{}.Log(OTSLoggerEnable, "OTS request failed, API: %s, HTTPStatus: %s, ErrorCode: %d, ErrorMessage: %v,  RequestID: %s", api_name, reason, status, body, request_id)
		return ots_service_err.SetErrorMessage(strings.TrimSpace(string(body))).SetHttpStatus(reason).SetErrorCode(fmt.Sprintf("%d", status)).SetRequestId(request_id)
	}

	return nil
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
	coder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	return coder.EncodeToString(src)
}

// base64 decode
func base64Decode(src []byte) ([]byte, error) {
	coder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
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
