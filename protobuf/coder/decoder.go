// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// encoder for ots2
package coder

import (
	// "fmt"
	"reflect"

	"code.google.com/p/goprotobuf/proto"
	. "github.com/GiterLab/goots/log"
	. "github.com/GiterLab/goots/otstype"
	. "github.com/GiterLab/goots/protobuf"
)

var api_decode_map = NewFuncmap()

type ots_proto_buffer_decoder map[string]interface{}

var api_decode_list = ots_proto_buffer_decoder{
	"CreateTable":   _decode_create_table,
	"DeleteTable":   _decode_delete_table,
	"ListTable":     _decode_list_table,
	"UpdateTable":   _decode_update_table,
	"DescribeTable": _decode_describe_table,
	"GetRow":        _decode_get_row,
	"PutRow":        _decode_put_row,
	"UpdateRow":     _decode_update_row,
	"DeleteRow":     _decode_delete_row,
	"BatchGetRow":   _decode_batch_get_row,
	"BatchWriteRow": _decode_batch_write_row,
	"GetRange":      _decode_get_range,
}

func init() {
	for k, v := range api_decode_list {
		api_decode_map.Bind(k, v)
	}
}

func _parse_string(str string) *string {
	if str == "" {
		return nil
	}

	return &str
}

func _parse_column_type() {

}

func _parse_value() {

}

func _parse_schema_list() {

}

func _parse_column_dict() {

}

func _parse_row() {

}

func _parse_row_list() {

}

func _parse_capacity_unit() {

}

func _parse_reserved_throughput_details() {

}

func _parse_get_row_item() {

}

func _parse_batch_get_row() {

}

func _parse_write_row_item() {

}

func _parse_batch_write_row() {

}

func _decode_create_table() {

}

func _decode_delete_table() {

}

func _decode_list_table(buf []byte) (list_tables *OTSListTableResponse, err error) {
	pb := &ListTableResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}

	list_tables = new(OTSListTableResponse)
	list_tables.TableNames = make([]string, len(pb.TableNames))
	copy(list_tables.TableNames, pb.TableNames)

	return list_tables, nil
}

func _decode_update_table() {

}

func _decode_describe_table() {

}

func _decode_get_row() {

}

func _decode_put_row() {

}

func _decode_update_row() {

}

func _decode_delete_row() {

}

func _decode_batch_get_row() {

}

func _decode_batch_write_row() {

}

func _decode_get_range() {

}

func decode_response() {

}

// request encode for ots2
func DecodeRequest(api_name string, args ...interface{}) (req []reflect.Value, err error) {
	if _, ok := api_decode_map[api_name]; !ok {
		return nil, (OTSClientError{}.Set("No PB decode method for API %s", api_name))
	}

	req, err = api_decode_map.Call(api_name, args...)
	if err != nil {
		return nil, err
	}

	return req, nil
}
