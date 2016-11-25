// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// encoder for ots2
package coder

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	. "github.com/GiterLab/goots/otstype"
	. "github.com/GiterLab/goots/protobuf"
	"github.com/golang/protobuf/proto"
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
		panic(errors.New(fmt.Sprintf("invalid column value type: %d", value.GetType())))
	}

	return nil
}

func _parse_schema_list(primary_key []*ColumnSchema) OTSSchemaOfPrimaryKey {
	schema_of_primary_key := make(OTSSchemaOfPrimaryKey, len(primary_key))
	for i, v := range primary_key {
		key := v.GetName()
		value := _parse_column_type(v.GetType())
		schema_of_primary_key[i].SetKey(key)
		schema_of_primary_key[i].SetValue(value)
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

func _parse_get_row_item(row_list []*RowInBatchGetRowResponse) []*OTSRowInBatchGetRowResponseItem {
	if len(row_list) == 0 {
		return nil
	}

	pobj := make([]*OTSRowInBatchGetRowResponseItem, len(row_list))
	for i, v := range row_list {
		row_item := new(OTSRowInBatchGetRowResponseItem)
		if v.GetIsOk() {
			row_item.IsOk = v.GetIsOk()
			row_item.ErrorCode = "None"
			row_item.ErrorMessage = "None"
			row_item.Consumed = _parse_capacity_unit(v.GetConsumed().GetCapacityUnit())
			row_item.Row = _parse_row(v.GetRow())

		} else {
			row_item.IsOk = v.GetIsOk()
			row_item.ErrorCode = v.GetError().GetCode()
			row_item.ErrorMessage = v.GetError().GetMessage()
			row_item.Consumed = nil
			row_item.Row = nil
		}

		pobj[i] = row_item
	}

	return pobj
}

func _parse_batch_get_row(table_list []*TableInBatchGetRowResponse) []*OTSTableInBatchGetRowResponseItem {
	if len(table_list) == 0 {
		return nil
	}

	pobj := make([]*OTSTableInBatchGetRowResponseItem, len(table_list))
	for i, v := range table_list {
		table_item := new(OTSTableInBatchGetRowResponseItem)
		table_item.TableName = v.GetTableName()
		table_item.Rows = _parse_get_row_item(v.GetRows())
		pobj[i] = table_item
	}

	return pobj
}

func _parse_write_row_item(row_list []*RowInBatchWriteRowResponse) []*OTSRowInBatchWriteRowResponseItem {
	if len(row_list) == 0 {
		return nil
	}

	pobj := make([]*OTSRowInBatchWriteRowResponseItem, len(row_list))
	for i, v := range row_list {
		row_item := new(OTSRowInBatchWriteRowResponseItem)
		if v.GetIsOk() {
			row_item.IsOk = v.GetIsOk()
			row_item.ErrorCode = "None"
			row_item.ErrorMessage = "None"
			row_item.Consumed = _parse_capacity_unit(v.GetConsumed().GetCapacityUnit())

		} else {
			row_item.IsOk = v.GetIsOk()

			row_item.ErrorCode = v.GetError().GetCode()
			row_item.ErrorMessage = v.GetError().GetMessage()
			row_item.Consumed = nil
		}

		pobj[i] = row_item
	}

	return pobj
}

func _parse_batch_write_row(table_list []*TableInBatchWriteRowResponse) []*OTSTableInBatchWriteRowResponseItem {
	if len(table_list) == 0 {
		return nil
	}

	pobj := make([]*OTSTableInBatchWriteRowResponseItem, len(table_list))
	for i, v := range table_list {
		table_item := new(OTSTableInBatchWriteRowResponseItem)
		table_item.TableName = v.GetTableName()
		table_item.PutRows = _parse_write_row_item(v.GetPutRows())
		table_item.UpdateRows = _parse_write_row_item(v.GetUpdateRows())
		table_item.DeleteRows = _parse_write_row_item(v.GetDeleteRows())
		pobj[i] = table_item
	}

	return pobj
}

func _decode_create_table(buf []byte) (err error) {
	pb := &CreateTableResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return err
	}
	print_response_message(pb)

	return nil
}

func _decode_delete_table(buf []byte) (err error) {
	pb := &DeleteTableResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return err
	}
	print_response_message(pb)

	return nil
}

func _decode_list_table(buf []byte) (list_tables *OTSListTableResponse, err error) {
	pb := &ListTableResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}
	print_response_message(pb)

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
	print_response_message(pb)

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
	print_response_message(pb)

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
	print_response_message(pb)

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
	print_response_message(pb)

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
	print_response_message(pb)

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
	print_response_message(pb)

	delete_row_response = new(OTSDeleteRowResponse)
	delete_row_response.Consumed = _parse_capacity_unit(pb.GetConsumed().GetCapacityUnit())

	return delete_row_response, nil
}

func _decode_batch_get_row(buf []byte) (response_item_list *OTSBatchGetRowResponse, err error) {
	pb := &BatchGetRowResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}
	print_response_message(pb)

	response_item_list = new(OTSBatchGetRowResponse)
	response_item_list.Tables = _parse_batch_get_row(pb.GetTables())

	return response_item_list, nil
}

func _decode_batch_write_row(buf []byte) (response_item_list *OTSBatchWriteRowResponse, err error) {
	pb := &BatchWriteRowResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}
	print_response_message(pb)

	response_item_list = new(OTSBatchWriteRowResponse)
	response_item_list.Tables = _parse_batch_write_row(pb.GetTables())

	return response_item_list, nil
}

func _decode_get_range(buf []byte) (response_row_list *OTSGetRangeResponse, err error) {
	pb := &GetRangeResponse{}
	err = proto.Unmarshal(buf, pb)
	if err != nil {
		return nil, err
	}
	print_response_message(pb)

	response_row_list = new(OTSGetRangeResponse)
	response_row_list.Consumed = _parse_capacity_unit(pb.GetConsumed().GetCapacityUnit())
	response_row_list.NextStartPrimaryKey = (OTSPrimaryKey)(_parse_column_dict(pb.GetNextStartPrimaryKey()))
	response_row_list.Rows = _parse_row_list(pb.GetRows())

	return response_row_list, nil
}

// request encode for ots2
func DecodeRequest(api_name string, args ...interface{}) (req []reflect.Value, err error) {
	if _, ok := api_decode_map[api_name]; !ok {
		return nil, errors.New(fmt.Sprintf("No PB decode method for API %s" + api_name))
	}

	req, err = api_decode_map.Call(api_name, args...)
	if err != nil {
		return nil, err
	}

	return req, nil
}
