// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// encoder for ots2
package coder

import (
	// "fmt"
	"reflect"
	"time"

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

func _parse_column_type(column_type_enum ColumnType) string {
	return ColumnType_name[int32(column_type_enum)]
}

func _parse_value(value *ColumnValue) interface{} {
	switch value.GetType() {
	case ColumnType_INTEGER:
		return value.GetVInt()
	case ColumnType_STRING:
		return value.GetVString()
	case ColumnType_BOOLEAN:
		return value.GetVBool()
	case ColumnType_DOUBLE:
		return value.GetVDouble()
	case ColumnType_BINARY:
		return value.GetVBinary()
	default:
		panic(OTSClientError{}.Set("invalid column value type: %d", value.GetType()))
	}

	return nil
}

func _parse_schema_list(primary_key []*ColumnSchema) OTSSchemaOfPrimaryKey {
	schema_of_primary_key := make(OTSSchemaOfPrimaryKey, len(primary_key))
	for _, v := range primary_key {
		key := v.GetName()
		value := _parse_column_type(v.GetType())
		schema_of_primary_key[key] = value
	}

	return schema_of_primary_key
}

func _parse_column_dict(colum []*Column) DictString {
	if len(colum) == 0 {
		return nil
	}

	dict := make(DictString, len(colum))
	for _, v := range colum {
		dict[v.GetName()] = _parse_value(v.GetValue())
	}

	return dict
}

func _parse_row(row *Row) *OTSRow {
	if row == nil {
		return nil
	}

	ots_row := new(OTSRow)
	ots_row.PrimaryKeyColumns = (OTSPrimaryKey)(_parse_column_dict(row.GetPrimaryKeyColumns()))
	ots_row.AttributeColumns = (OTSAttribute)(_parse_column_dict(row.GetAttributeColumns()))

	return ots_row
}

func _parse_row_list(rows []*Row) OTSRows {
	if len(rows) == 0 {
		return nil
	}

	ots_rows := make(OTSRows, len(rows))
	for i, v := range rows {
		ots_rows[i] = _parse_row(v)
	}

	return ots_rows
}

func _parse_table_meta(table_meta *TableMeta) *OTSTableMeta {
	if table_meta == nil {
		return nil
	}

	pobj := new(OTSTableMeta)
	pobj.TableName = table_meta.GetTableName()
	pobj.SchemaOfPrimaryKey = _parse_schema_list(table_meta.PrimaryKey)

	return pobj
}

func _parse_capacity_unit(capacity_unit *CapacityUnit) *OTSCapacityUnit {
	if capacity_unit == nil {
		return nil
	}

	pobj := new(OTSCapacityUnit)
	pobj.Read = capacity_unit.GetRead()
	pobj.Write = capacity_unit.GetWrite()

	return pobj
}

func _parse_reserved_throughput_details(reserved_throughput_details *ReservedThroughputDetails) *OTSReservedThroughputDetails {
	if reserved_throughput_details == nil {
		return nil
	}

	pobj := new(OTSReservedThroughputDetails)
	pobj.CapacityUnit = _parse_capacity_unit(reserved_throughput_details.GetCapacityUnit())
	pobj.LastDecreaseTime = time.Unix(reserved_throughput_details.GetLastDecreaseTime(), 0)
	pobj.LastIncreaseTime = time.Unix(reserved_throughput_details.GetLastIncreaseTime(), 0)
	pobj.NumberOfDecreasesToday = reserved_throughput_details.GetNumberOfDecreasesToday()

	return pobj
}

func _parse_get_row_item() {

}

func _parse_batch_get_row() {

}

func _parse_write_row_item() {

}

func _parse_batch_write_row() {

}

func _decode_create_table(buf []byte) (err error) {
	pb := &CreateTableResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return err
	}

	return nil
}

func _decode_delete_table(buf []byte) (err error) {
	pb := &DeleteTableResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return err
	}

	return nil
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

func _decode_update_table(buf []byte) (update_table_response *OTSUpdateTableResponse, err error) {
	pb := &UpdateTableResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}

	update_table_response = new(OTSUpdateTableResponse)
	update_table_response.ReservedThroughputDetails = _parse_reserved_throughput_details(pb.GetReservedThroughputDetails())

	return update_table_response, nil
}

func _decode_describe_table(buf []byte) (describe_table_response *OTSDescribeTableResponse, err error) {
	pb := &DescribeTableResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}

	describe_table_response = new(OTSDescribeTableResponse)
	describe_table_response.TableMeta = _parse_table_meta(pb.GetTableMeta())
	describe_table_response.ReservedThroughputDetails = _parse_reserved_throughput_details(pb.GetReservedThroughputDetails())

	return describe_table_response, nil
}

func _decode_get_row(buf []byte) (get_row_response *OTSGetRowResponse, err error) {
	pb := &GetRowResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}

	get_row_response = new(OTSGetRowResponse)
	get_row_response.Row = _parse_row(pb.GetRow())
	get_row_response.Consumed = _parse_capacity_unit(pb.GetConsumed().GetCapacityUnit())

	return get_row_response, nil
}

func _decode_put_row(buf []byte) (put_row_response *OTSPutRowResponse, err error) {
	pb := &PutRowResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}

	put_row_response = new(OTSPutRowResponse)
	put_row_response.Consumed = _parse_capacity_unit(pb.GetConsumed().GetCapacityUnit())

	return put_row_response, nil
}

func _decode_update_row(buf []byte) (update_row_response *OTSUpdateRowResponse, err error) {
	pb := &UpdateRowResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}

	update_row_response = new(OTSUpdateRowResponse)
	update_row_response.Consumed = _parse_capacity_unit(pb.GetConsumed().GetCapacityUnit())

	return update_row_response, nil
}

func _decode_delete_row(buf []byte) (delete_row_response *OTSDeleteRowResponse, err error) {
	pb := &DeleteRowResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}

	delete_row_response = new(OTSDeleteRowResponse)
	delete_row_response.Consumed = _parse_capacity_unit(pb.GetConsumed().GetCapacityUnit())

	return delete_row_response, nil
}

func _decode_batch_get_row(buf []byte) {

}

func _decode_batch_write_row(buf []byte) {

}

func _decode_get_range(buf []byte) {

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
