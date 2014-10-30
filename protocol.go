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
	"net/url"
	"sort"
	"strings"
	// "reflect"
	// "sync"
	"time"

	// . "github.com/GiterLab/goots/log"
	. "github.com/GiterLab/goots/otstype"
	// "github.com/GiterLab/goots/urllib"
	// "code.google.com/p/goprotobuf/proto"
)

var API_VERSION = "2014-08-08"
var defaultProtocol = ots_protocol{
	api_version: API_VERSION,
}

func newProtocol(protocol *ots_protocol) *ots_protocol {
	if protocol == nil {
		return &defaultProtocol
	}

	protocol = new(ots_protocol)
	protocol.api_version = API_VERSION

	return protocol
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
		if strings.HasPrefix(k, "x-ots-") && k != "x-ots-signature" {
			count++
		}
	}
	if count == 0 {
		return "\n"
	}

	strslice := make([]string, count)
	i := 0
	for k, v := range headers {
		if strings.HasPrefix(k, "x-ots-") && k != "x-ots-signature" {
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
	// fmt.Println("uri:", uri)
	// fmt.Println("sorted_query:", sorted_query)
	// fmt.Println("url_obj.Opaque:", url_obj.Opaque)
	// fmt.Println("signature_string, before:", signature_string)
	headers_string := o._make_headers_string(headers)
	// fmt.Println("headers_string:", headers_string)
	signature_string = signature_string + headers_string + "\n"
	// fmt.Println("signature_string, after:", signature_string)
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
