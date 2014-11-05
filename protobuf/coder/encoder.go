// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// encoder for ots2
package coder

import (
	"fmt"
	"reflect"

	// "code.google.com/p/goprotobuf/proto"
	. "github.com/GiterLab/goots/log"
	. "github.com/GiterLab/goots/otstype"
	. "github.com/GiterLab/goots/protobuf"
)

const (
	INT32_MAX int32 = 2147483647
	INT32_MIN int32 = -2147483648
)

var api_encode_map = NewFuncmap()

type ots_proto_buffer_encoder map[string]interface{}

var api_encode_list = ots_proto_buffer_encoder{
	"CreateTable":   _encode_create_table,
	"DeleteTable":   _encode_delete_table,
	"ListTable":     _encode_list_table,
	"UpdateTable":   _encode_update_table,
	"DescribeTable": _encode_describe_table,
	"GetRow":        _encode_get_row,
	"PutRow":        _encode_put_row,
	"UpdateRow":     _encode_update_row,
	"DeleteRow":     _encode_delete_row,
	"BatchGetRow":   _encode_batch_get_row,
	"BatchWriteRow": _encode_batch_write_row,
	"GetRange":      _encode_get_range,
}

func init() {
	for k, v := range api_encode_list {
		api_encode_map.Bind(k, v)
	}
}

func _get_unicode(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	} else {
		panic(OTSClientError{}.Set("expect str or unicode type for string, not %v", reflect.TypeOf(value)))
	}

	return ""
}

func _get_int32(value interface{}) int32 {
	if v, ok := value.(int); ok {
		if int32(v) < INT32_MIN || int32(v) > INT32_MAX {
			panic(OTSClientError{}.Set("%d exceeds the range of int32", v))
		}
		return int32(v)
	} else if v, ok := value.(int32); ok {
		if v < INT32_MIN || v > INT32_MAX {
			panic(OTSClientError{}.Set("%d exceeds the range of int32", v))
		}
		return v
	}

	panic(OTSClientError{}.Set("expect int or long for the value, not %v", reflect.TypeOf(value)))
}

func _make_repeated_column_names(proto *[]string, columns_to_get []string) error {
	if columns_to_get == nil {
		// if no column name is given, get all primary_key_columns and attribute_columns.
		return nil
	}

	for _, column_name := range columns_to_get {
		*proto = append(*proto, _get_unicode(column_name))
	}

	// *proto = columns_to_get[:] // not used

	return nil
}

func _make_column_value(proto *ColumnValue, value interface{}) error {
	// you have to put 'int' under 'bool' in the switch case
	// because a bool is also a int !!!

	switch vtype := value.(type) {
	case string:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_STRING
		proto.Type = pcolumn_type
		proto.VString = NewString(_get_unicode(value.(string)))

	case bool:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_BOOLEAN
		proto.Type = pcolumn_type
		proto.VBool = NewBool(value.(bool))

	case int:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_INTEGER
		proto.Type = pcolumn_type
		proto.VInt = NewInt64(int64(value.(int)))

	case uint:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_INTEGER
		proto.Type = pcolumn_type
		proto.VInt = NewInt64(int64(value.(uint)))

	case int8:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_INTEGER
		proto.Type = pcolumn_type
		proto.VInt = NewInt64(int64(value.(int8)))

	case uint8:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_INTEGER
		proto.Type = pcolumn_type
		proto.VInt = NewInt64(int64(value.(uint8)))

	case int32:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_INTEGER
		proto.Type = pcolumn_type
		proto.VInt = NewInt64(int64(value.(int32)))

	case uint32:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_INTEGER
		proto.Type = pcolumn_type
		proto.VInt = NewInt64(int64(value.(uint32)))

	case int64:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_INTEGER
		proto.Type = pcolumn_type
		proto.VInt = NewInt64(value.(int64))

	case uint64:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_INTEGER
		proto.Type = pcolumn_type
		proto.VInt = NewInt64(int64(value.(uint64)))

	case float32:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_DOUBLE
		proto.Type = pcolumn_type
		proto.VDouble = NewFloat64(float64(value.(float32)))

	case float64:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_DOUBLE
		proto.Type = pcolumn_type
		proto.VDouble = NewFloat64(value.(float64))

	case []byte:
		pcolumn_type := new(ColumnType)
		*pcolumn_type = ColumnType_BINARY
		proto.Type = pcolumn_type
		proto.VBinary = value.([]byte)

	case ColumnType:
		v := value.(ColumnType)
		if v == ColumnType_INF_MIN {
			pcolumn_type := new(ColumnType)
			*pcolumn_type = ColumnType_INF_MIN
			proto.Type = pcolumn_type
		} else if v == ColumnType_INF_MAX {
			pcolumn_type := new(ColumnType)
			*pcolumn_type = ColumnType_INF_MAX
			proto.Type = pcolumn_type
		} else {
			return (OTSClientError{}.Set("don't expect the value of ColumnType"))
		}

	default:
		return (OTSClientError{}.Set("expect string, bool, (u)int, (u)int8, (u)int16, (u)int32, (u)int64, (u)float32 or (u)float64 for colum value, not %v", vtype))
	}

	return nil
}

func _get_column_type(type_str string) ColumnType {
	v, ok := ColumnType_value[type_str]
	if ok {
		return ColumnType(v)
	} else {
		panic(OTSClientError{}.Set("column_type should be one of [INF_MIN, INF_MAX, INTEGER, STRING, BOOLEAN, DOUBLE, BINARY], not %s", type_str))
	}
}

func _make_condition(proto *Condition, condition interface{}) error {
	switch condition.(type) {
	case Condition:
		exp := *condition.(Condition).RowExistence
		if v, ok := RowExistenceExpectation_name[int32(exp)]; ok {
			item := new(RowExistenceExpectation)
			*item = RowExistenceExpectation(RowExistenceExpectation_value[v])
			proto.RowExistence = item
		} else {
			return (OTSClientError{}.Set("condition value should be one of [IGNORE(0), EXPECT_EXIST(1), EXPECT_NOT_EXIST(2)], not %v", exp))
		}
	case *Condition:
		exp := *condition.(*Condition).RowExistence
		if v, ok := RowExistenceExpectation_name[int32(exp)]; ok {
			item := new(RowExistenceExpectation)
			*item = RowExistenceExpectation(RowExistenceExpectation_value[v])
			proto.RowExistence = item
		} else {
			return (OTSClientError{}.Set("condition value should be one of [IGNORE(0), EXPECT_EXIST(1), EXPECT_NOT_EXIST(2)], not %v", exp))
		}
	case string:
		exp := condition.(string)
		if v, ok := RowExistenceExpectation_value[exp]; ok {
			item := new(RowExistenceExpectation)
			*item = RowExistenceExpectation(v)
			proto.RowExistence = item
		} else {
			return (OTSClientError{}.Set("condition value should be one of [IGNORE(0), EXPECT_EXIST(1), EXPECT_NOT_EXIST(2)], not %v", exp))
		}
	default:
		return (OTSClientError{}.Set("condition should be one of [Condition, *Condition or string], not %v", reflect.TypeOf(condition)))
	}

	return nil
}

func _get_condition(condition_str string) Condition {
	if v, ok := RowExistenceExpectation_value[condition_str]; ok {
		item := new(RowExistenceExpectation)
		*item = RowExistenceExpectation(v)
		conditon := Condition{}
		conditon.RowExistence = item
	} else {
		panic(OTSClientError{}.Set("direction should be one of [IGNORE, EXPECT_EXIST, EXPECT_NOT_EXIST], not %s", condition_str))
	}

	return Condition{}
}

func _get_direction(direction_str string) Direction {
	v, ok := Direction_value[direction_str]
	if ok {
		return Direction(v)
	} else {
		panic(OTSClientError{}.Set("direction should be one of [FORWARD, BACKWARD], not %s", direction_str))
	}
}

func _make_column_schema(proto *ColumnSchema, schema_tuple interface{}) error {
	switch schema_tuple.(type) {
	case ColumnSchema:
		schema_name, schema_type := *schema_tuple.(ColumnSchema).Name, *schema_tuple.(ColumnSchema).Type
		proto.Name = new(string)
		proto.Type = new(ColumnType)
		*proto.Name = _get_unicode(schema_name)
		*proto.Type = schema_type
	case TupleString:
		schema_name, schema_type := schema_tuple.(TupleString).GetName(), schema_tuple.(TupleString).GetType()
		proto.Name = new(string)
		proto.Type = new(ColumnType)
		*proto.Name = _get_unicode(schema_name)
		if v, ok := schema_type.(string); ok {
			*proto.Type = _get_column_type(v)
		} else {
			return (OTSClientError{}.Set("schema_tuple should be (string, string), not (string, %v)", reflect.TypeOf(schema_type)))
		}
	default:
		return (OTSClientError{}.Set("type of schema_list is shoud be one of [ColumnSchema or TupleString]. not %v", reflect.TypeOf(schema_tuple)))
	}

	return nil
}

// Deprecated
func _make_column_schema_python(proto *ColumnSchema, schema_tuple TupleString) error {
	schema_name, schema_type := schema_tuple.GetName(), schema_tuple.GetType()
	proto.Name = new(string)
	proto.Type = new(ColumnType)
	*proto.Name = _get_unicode(schema_name)
	if v, ok := schema_type.(string); ok {
		*proto.Type = _get_column_type(v)
	} else {
		return (OTSClientError{}.Set("schema_tuple should be (string, string), not (string, %v)", reflect.TypeOf(schema_type)))
	}

	return nil
}

func _make_schemas_with_list(proto *[]*ColumnSchema, schema_list interface{}) error {
	switch schema_list.(type) {
	case []ColumnSchema:
		if len(schema_list.([]ColumnSchema)) == 0 {
			return OTSClientError{}.Set("schema_list should not be empty")
		}
		*proto = make([]*ColumnSchema, len(schema_list.([]ColumnSchema)))
		for k, schema_tuple := range schema_list.([]ColumnSchema) {
			item := new(ColumnSchema)
			err := _make_column_schema(item, schema_tuple)
			if err != nil {
				return err
			}
			(*proto)[k] = item
		}

	case []*ColumnSchema:
		if len(schema_list.([]*ColumnSchema)) == 0 {
			return OTSClientError{}.Set("schema_list should not be empty")
		}
		*proto = make([]*ColumnSchema, len(schema_list.([]*ColumnSchema)))
		for k, schema_tuple := range schema_list.([]*ColumnSchema) {
			item := new(ColumnSchema)
			err := _make_column_schema(item, *schema_tuple)
			if err != nil {
				return err
			}
			(*proto)[k] = item
		}

	case []TupleString:
		if len(schema_list.([]TupleString)) == 0 {
			return OTSClientError{}.Set("schema_list should not be empty")
		}
		*proto = make([]*ColumnSchema, len(schema_list.([]TupleString)))
		for k, schema_tuple := range schema_list.([]TupleString) {
			item := new(ColumnSchema)
			// _make_column_schema_python(item, schema_tuple)
			err := _make_column_schema(item, schema_tuple)
			if err != nil {
				return err
			}
			(*proto)[k] = item
		}

	default:
		return (OTSClientError{}.Set("type of schema_list is shoud be one of [[]ColumnSchema []*ColumnSchema or []TupleString]. not %v", reflect.TypeOf(schema_list)))
	}
	return nil
}

func _make_columns_with_dict(proto *[]*Column, column_dict interface{}) error {
	switch column_dict.(type) {
	case []Column:
		if len(column_dict.([]Column)) == 0 {
			return OTSClientError{}.Set("column_dict should not be empty")
		}
		*proto = make([]*Column, len(column_dict.([]Column)))
		for k, column := range column_dict.([]Column) {
			item := new(Column)
			item.Name = new(string)
			*item.Name = _get_unicode(column.GetName())
			item.Value = column.GetValue()
			(*proto)[k] = item
		}

	case []*Column:
		if len(column_dict.([]*Column)) == 0 {
			return OTSClientError{}.Set("column_dict should not be empty")
		}
		*proto = make([]*Column, len(column_dict.([]*Column)))
		for k, column := range column_dict.([]*Column) {
			item := new(Column)
			item.Name = new(string)
			*item.Name = _get_unicode((*column).GetName())
			item.Value = (*column).GetValue()
			(*proto)[k] = item
		}

	case DictString:
		if len(column_dict.(DictString)) == 0 {
			return OTSClientError{}.Set("column_dict should not be empty")
		}
		*proto = make([]*Column, len(column_dict.(DictString)))
		i := 0
		for name, column := range column_dict.(DictString) {
			item := new(Column)
			item.Name = NewString(name)
			item.Value = new(ColumnValue)
			_make_column_value(item.Value, column)
			(*proto)[i] = item
			i++
		}

	default:
		return (OTSClientError{}.Set("type of schema_list is shoud be one of [[]Column []*Column or DictString]. not %v", reflect.TypeOf(column_dict)))
	}
	return nil
}

func _make_update_of_attribute_columns_with_dict(proto *[]*ColumnUpdate, column_dict interface{}) error {
	switch column_dict.(type) {
	case []ColumnUpdate:
		if len(column_dict.([]ColumnUpdate)) == 0 {
			return OTSClientError{}.Set("column_dict should not be empty")
		}
		*proto = make([]*ColumnUpdate, len(column_dict.([]ColumnUpdate)))
		for k, column_update := range column_dict.([]ColumnUpdate) {
			item := new(ColumnUpdate)
			item.Type = new(OperationType)
			*item.Type = column_update.GetType()
			item.Name = new(string)
			*item.Name = _get_unicode(column_update.GetName())
			item.Value = new(ColumnValue)
			item.Value = column_update.GetValue()
			(*proto)[k] = item
		}
	case []*ColumnUpdate:
		if len(column_dict.([]*ColumnUpdate)) == 0 {
			return OTSClientError{}.Set("column_dict should not be empty")
		}
		*proto = make([]*ColumnUpdate, len(column_dict.([]*ColumnUpdate)))
		for k, column_update := range column_dict.([]*ColumnUpdate) {
			item := new(ColumnUpdate)
			item.Type = new(OperationType)
			*item.Type = (*column_update).GetType()
			item.Name = new(string)
			*item.Name = _get_unicode((*column_update).GetName())
			item.Value = new(ColumnValue)
			item.Value = (*column_update).GetValue()
			(*proto)[k] = item
		}
	case DictString:
		// DictString --> map[string] DictString --> map[string] map[string]interface
		*proto = make([]*ColumnUpdate, 0, 10) // modify 10 to big value to fit your app
		for key, value := range column_dict.(DictString) {
			if key == "PUT" {
				// value.(DictString) or value.(OTSColumnsToPut)
				switch value.(type) {
				case DictString:
					for k, v := range value.(DictString) {
						item := new(ColumnUpdate)
						item.Type = new(OperationType)
						*item.Type = OperationType_PUT
						item.Name = new(string)
						*item.Name = _get_unicode(k)
						item.Value = new(ColumnValue)
						_make_column_value(item.Value, v)
						*proto = append(*proto, item)
					}
				case OTSColumnsToPut:
					for k, v := range value.(OTSColumnsToPut) {
						item := new(ColumnUpdate)
						item.Type = new(OperationType)
						*item.Type = OperationType_PUT
						item.Name = new(string)
						*item.Name = _get_unicode(k)
						item.Value = new(ColumnValue)
						_make_column_value(item.Value, v)
						*proto = append(*proto, item)
					}
				default:
					return (OTSClientError{}.Set("expect DictString  or OTSColumnsToPut for put operation in 'update_of_attribute_columns', not %v", reflect.TypeOf(value)))
				}

			} else if key == "DELETE" {
				switch value.(type) {
				case []string:
					for _, v := range value.([]string) {
						item := new(ColumnUpdate)
						item.Type = new(OperationType)
						*item.Type = OperationType_DELETE
						item.Name = new(string)
						*item.Name = _get_unicode(v)
						*proto = append(*proto, item)
					}
				case OTSColumnsToDelete:
					for _, v := range value.(OTSColumnsToDelete) {
						item := new(ColumnUpdate)
						item.Type = new(OperationType)
						*item.Type = OperationType_DELETE
						item.Name = new(string)
						*item.Name = _get_unicode(v)
						*proto = append(*proto, item)
					}
				default:
					return (OTSClientError{}.Set("expect list([]string or OTSColumnsToDelete) for delete operation in 'update_of_attribute_columns', not %v", reflect.TypeOf(value)))
				}

			} else {
				return (OTSClientError{}.Set("operation type in 'update_of_attribute_columns' should be 'PUT' or 'DELETE', not %s", key))
			}
		}

	default:
		return (OTSClientError{}.Set("expect DictString or []ColumnUpdate for 'update_of_attribute_columns', not %v", reflect.TypeOf(column_dict)))
	}
	return nil
}

func _make_table_meta(proto *TableMeta, table_meta interface{}) error {
	switch table_meta.(type) {
	case TableMeta:
		proto.TableName = new(string)
		*proto.TableName = _get_unicode(*table_meta.(TableMeta).TableName)
		primary_key := new([]*ColumnSchema)
		err := _make_schemas_with_list(primary_key, table_meta.(TableMeta).PrimaryKey)
		if err != nil {
			return err
		}
		proto.PrimaryKey = (*primary_key)[:]
	case TupleString:
		proto.TableName = new(string)
		*proto.TableName = _get_unicode(table_meta.(TupleString).GetName())
		if v, ok := table_meta.(TupleString).V.([]TupleString); ok {
			primary_key := new([]*ColumnSchema)
			err := _make_schemas_with_list(primary_key, v)
			if err != nil {
				return err
			}
			proto.PrimaryKey = (*primary_key)[:]
		} else {
			return (OTSClientError{}.Set("table_meta.V should be an instance of []TupleString, not %v", reflect.TypeOf(table_meta.(TupleString).V)))
		}
	case OTSTableMeta:
		proto.TableName = NewString(table_meta.(OTSTableMeta).TableName)
		primary_key := new([]*ColumnSchema)

		// change map[string]string to []TupleString
		tuple_string := make([]TupleString, len(table_meta.(OTSTableMeta).SchemaOfPrimaryKey))
		i := 0
		for k, v := range table_meta.(OTSTableMeta).SchemaOfPrimaryKey {
			tuple_string[i].K = k
			tuple_string[i].V = v
			i++
		}
		err := _make_schemas_with_list(primary_key, tuple_string)
		if err != nil {
			return err
		}
		proto.PrimaryKey = (*primary_key)[:]
	default:
		return (OTSClientError{}.Set("table_meta should be an instance of TableMeta, OTSTableMeta or TupleString, not %v", reflect.TypeOf(table_meta)))
	}

	return nil
}

func _make_capacity_unit(proto *CapacityUnit, capacity_unit interface{}) error {
	switch capacity_unit.(type) {
	case CapacityUnit:
		if capacity_unit.(CapacityUnit).Read == nil || capacity_unit.(CapacityUnit).Write == nil {
			return (OTSClientError{}.Set("both of read and write of CapacityUnit are required"))
		}
		proto.Read = NewInt32(_get_int32(*capacity_unit.(CapacityUnit).Read))
		proto.Write = NewInt32(_get_int32(*capacity_unit.(CapacityUnit).Write))

	case OTSCapacityUnit:
		if capacity_unit.(OTSCapacityUnit).Read == 0 || capacity_unit.(OTSCapacityUnit).Write == 0 {
			return (OTSClientError{}.Set("both of read and write of CapacityUnit are required"))
		}
		proto.Read = NewInt32(_get_int32(capacity_unit.(OTSCapacityUnit).Read))
		proto.Write = NewInt32(_get_int32(capacity_unit.(OTSCapacityUnit).Write))
	}

	return nil
}

func _make_reserved_throughput(proto *ReservedThroughput, reserved_throughput interface{}) error {
	switch reserved_throughput.(type) {
	case ReservedThroughput:
		capacity_unit := *reserved_throughput.(ReservedThroughput).CapacityUnit
		proto.CapacityUnit = new(CapacityUnit)
		err := _make_capacity_unit(proto.CapacityUnit, capacity_unit)
		if err != nil {
			return err
		}

	case OTSReservedThroughput:
		capacity_unit := reserved_throughput.(OTSReservedThroughput).CapacityUnit
		proto.CapacityUnit = new(CapacityUnit)
		err := _make_capacity_unit(proto.CapacityUnit, capacity_unit)
		if err != nil {
			return err
		}

	default:
		return (OTSClientError{}.Set("reserved_throughput should be an instance of [ReservedThroughput, OTSTableMeta or OTSReservedThroughput], not %v", reflect.TypeOf(reserved_throughput)))
	}

	return nil
}

func _make_update_capacity_unit(proto *CapacityUnit, capacity_unit interface{}) error {
	switch capacity_unit.(type) {
	case CapacityUnit:
		if capacity_unit.(CapacityUnit).Read == nil && capacity_unit.(CapacityUnit).Write == nil {
			return (OTSClientError{}.Set("at least one of read or write of CapacityUnit is required"))
		}

		if capacity_unit.(CapacityUnit).Read != nil {
			proto.Read = NewInt32(_get_int32(*capacity_unit.(CapacityUnit).Read))
		}

		if capacity_unit.(CapacityUnit).Write != nil {
			proto.Write = NewInt32(_get_int32(*capacity_unit.(CapacityUnit).Write))
		}

	case OTSCapacityUnit:
		if capacity_unit.(OTSCapacityUnit).Read == 0 && capacity_unit.(OTSCapacityUnit).Write == 0 {
			return (OTSClientError{}.Set("at least one of read or write of CapacityUnit is required"))
		}

		if capacity_unit.(OTSCapacityUnit).Read != 0 {
			proto.Read = NewInt32(_get_int32(capacity_unit.(OTSCapacityUnit).Read))
		}

		if capacity_unit.(OTSCapacityUnit).Write != 0 {
			proto.Write = NewInt32(_get_int32(capacity_unit.(OTSCapacityUnit).Write))
		}
	}

	return nil
}

func _make_update_reserved_throughput(proto *ReservedThroughput, reserved_throughput interface{}) error {
	switch reserved_throughput.(type) {
	case ReservedThroughput:
		capacity_unit := *reserved_throughput.(ReservedThroughput).CapacityUnit
		proto.CapacityUnit = new(CapacityUnit)
		err := _make_update_capacity_unit(proto.CapacityUnit, capacity_unit)
		if err != nil {
			return err
		}

	case OTSReservedThroughput:
		capacity_unit := reserved_throughput.(OTSReservedThroughput).CapacityUnit
		proto.CapacityUnit = new(CapacityUnit)
		err := _make_update_capacity_unit(proto.CapacityUnit, capacity_unit)
		if err != nil {
			return err
		}

	default:
		return (OTSClientError{}.Set("reserved_throughput should be an instance of ReservedThroughput, OTSTableMeta or OTSReservedThroughput, not %v", reflect.TypeOf(reserved_throughput)))
	}

	return nil
}

func _make_batch_get_row(proto *BatchGetRowRequest, batch_list interface{}) error {
	switch batch_list.(type) {
	case []TableInBatchGetRowRequest:
		list_len := len(batch_list.([]TableInBatchGetRowRequest))
		proto.Tables = make([]*TableInBatchGetRowRequest, list_len)
		for i, v := range batch_list.([]TableInBatchGetRowRequest) {
			table_item := new(TableInBatchGetRowRequest)
			// table_name
			table_name := _get_unicode(*v.TableName)
			table_item.TableName = new(string)
			*table_item.TableName = table_name
			// columns_to_get
			columns_to_get := new([]string)
			_make_repeated_column_names(columns_to_get, v.ColumnsToGet)
			table_item.ColumnsToGet = *columns_to_get
			// row_list
			table_item.Rows = make([]*RowInBatchGetRowRequest, len(v.Rows))
			for i1, v1 := range v.Rows {
				row := new(RowInBatchGetRowRequest)
				primary_key := new([]*Column)
				err := _make_columns_with_dict(primary_key, v1.PrimaryKey)
				if err != nil {
					return err
				}
				row.PrimaryKey = *primary_key
				table_item.Rows[i1] = row
			}
			proto.Tables[i] = table_item
		}

	case []*TableInBatchGetRowRequest:
		list_len := len(batch_list.([]*TableInBatchGetRowRequest))
		proto.Tables = make([]*TableInBatchGetRowRequest, list_len)
		for i, v := range batch_list.([]*TableInBatchGetRowRequest) {
			table_item := new(TableInBatchGetRowRequest)
			// table_name
			table_name := _get_unicode(*(*v).TableName)
			table_item.TableName = new(string)
			*table_item.TableName = table_name
			// columns_to_get
			columns_to_get := new([]string)
			_make_repeated_column_names(columns_to_get, (*v).ColumnsToGet)
			table_item.ColumnsToGet = *columns_to_get
			// row_list
			table_item.Rows = make([]*RowInBatchGetRowRequest, len(v.Rows))
			for i1, v1 := range v.Rows {
				row := new(RowInBatchGetRowRequest)
				primary_key := new([]*Column)
				err := _make_columns_with_dict(primary_key, v1.PrimaryKey)
				if err != nil {
					return err
				}
				row.PrimaryKey = *primary_key
				table_item.Rows[i1] = row
			}
			proto.Tables[i] = table_item
		}
	case OTSBatchGetRowRequest:
		list_len := len(batch_list.(OTSBatchGetRowRequest))
		proto.Tables = make([]*TableInBatchGetRowRequest, list_len)
		for i, v := range batch_list.(OTSBatchGetRowRequest) {
			table_item := new(TableInBatchGetRowRequest)
			// table_name
			table_item.TableName = NewString(v.TableName)
			// columns_to_get
			columns_to_get := new([]string)
			_make_repeated_column_names(columns_to_get, v.ColumnsToGet)
			table_item.ColumnsToGet = *columns_to_get
			// row_list
			table_item.Rows = make([]*RowInBatchGetRowRequest, len(v.Rows))
			for i1, v1 := range v.Rows {
				row := new(RowInBatchGetRowRequest)
				primary_key := new([]*Column)
				err := _make_columns_with_dict(primary_key, DictString(v1))
				if err != nil {
					return err
				}
				row.PrimaryKey = *primary_key
				table_item.Rows[i1] = row
			}
			proto.Tables[i] = table_item
		}

	default:
		return (OTSClientError{}.Set("batch_list should be an instance of [[]TableInBatchGetRowRequest, []*TableInBatchGetRowRequest or OTSBatchGetRowRequest], not %v", reflect.TypeOf(batch_list)))
	}

	return nil
}

func _make_put_row_item(proto *PutRowInBatchWriteRowRequest, put_row_item interface{}) error {
	switch put_row_item.(type) {
	case PutRowInBatchWriteRowRequest:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, *put_row_item.(PutRowInBatchWriteRowRequest).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, put_row_item.(PutRowInBatchWriteRowRequest).PrimaryKey)
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key
		attribute_columns := new([]*Column)
		err = _make_columns_with_dict(attribute_columns, put_row_item.(PutRowInBatchWriteRowRequest).AttributeColumns)
		if err != nil {
			return err
		}
		proto.AttributeColumns = *attribute_columns

	case *PutRowInBatchWriteRowRequest:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, *put_row_item.(*PutRowInBatchWriteRowRequest).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, put_row_item.(*PutRowInBatchWriteRowRequest).PrimaryKey)
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key
		attribute_columns := new([]*Column)
		err = _make_columns_with_dict(attribute_columns, put_row_item.(*PutRowInBatchWriteRowRequest).AttributeColumns)
		if err != nil {
			return err
		}
		proto.AttributeColumns = *attribute_columns

	case OTSPutRowItem:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, put_row_item.(OTSPutRowItem).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, DictString(put_row_item.(OTSPutRowItem).PrimaryKey))
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key
		attribute_columns := new([]*Column)
		err = _make_columns_with_dict(attribute_columns, DictString(put_row_item.(OTSPutRowItem).AttributeColumns))
		if err != nil {
			return err
		}
		proto.AttributeColumns = *attribute_columns

	default:
		return (OTSClientError{}.Set("put_row_item should be an instance of [PutRowInBatchWriteRowRequest, *PutRowInBatchWriteRowRequest or OTSPutRowItem], not %v", reflect.TypeOf(put_row_item)))
	}

	return nil
}

func _make_update_row_item(proto *UpdateRowInBatchWriteRowRequest, update_row_item interface{}) error {
	switch update_row_item.(type) {
	case UpdateRowInBatchWriteRowRequest:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, *update_row_item.(UpdateRowInBatchWriteRowRequest).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, update_row_item.(UpdateRowInBatchWriteRowRequest).PrimaryKey)
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key
		attribute_columns := new([]*ColumnUpdate)
		err = _make_update_of_attribute_columns_with_dict(attribute_columns, update_row_item.(UpdateRowInBatchWriteRowRequest).AttributeColumns)
		if err != nil {
			return err
		}
		proto.AttributeColumns = *attribute_columns

	case *UpdateRowInBatchWriteRowRequest:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, *update_row_item.(*UpdateRowInBatchWriteRowRequest).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, update_row_item.(*UpdateRowInBatchWriteRowRequest).PrimaryKey)
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key
		attribute_columns := new([]*ColumnUpdate)
		err = _make_update_of_attribute_columns_with_dict(attribute_columns, update_row_item.(*UpdateRowInBatchWriteRowRequest).AttributeColumns)
		if err != nil {
			return err
		}
		proto.AttributeColumns = *attribute_columns

	case OTSUpdateRowItem:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, update_row_item.(OTSUpdateRowItem).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, DictString(update_row_item.(OTSUpdateRowItem).PrimaryKey))
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key
		attribute_columns := new([]*ColumnUpdate)
		err = _make_update_of_attribute_columns_with_dict(attribute_columns, DictString(update_row_item.(OTSUpdateRowItem).UpdateOfAttributeColumns))
		if err != nil {
			return err
		}
		proto.AttributeColumns = *attribute_columns

	default:
		return (OTSClientError{}.Set("update_row_item should be an instance of [UpdateRowInBatchWriteRowRequest, *UpdateRowInBatchWriteRowRequest or OTSUpdateRowItem], not %v", reflect.TypeOf(update_row_item)))
	}

	return nil
}

func _make_delete_row_item(proto *DeleteRowInBatchWriteRowRequest, delete_row_item interface{}) error {
	switch delete_row_item.(type) {
	case DeleteRowInBatchWriteRowRequest:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, *delete_row_item.(DeleteRowInBatchWriteRowRequest).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, delete_row_item.(DeleteRowInBatchWriteRowRequest).PrimaryKey)
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key

	case *DeleteRowInBatchWriteRowRequest:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, *delete_row_item.(*DeleteRowInBatchWriteRowRequest).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, delete_row_item.(*DeleteRowInBatchWriteRowRequest).PrimaryKey)
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key

	case OTSDeleteRowItem:
		proto.Condition = new(Condition)
		err := _make_condition(proto.Condition, delete_row_item.(OTSDeleteRowItem).Condition)
		if err != nil {
			return err
		}
		primary_key := new([]*Column)
		err = _make_columns_with_dict(primary_key, DictString(delete_row_item.(OTSDeleteRowItem).PrimaryKey))
		if err != nil {
			return err
		}
		proto.PrimaryKey = *primary_key

	default:
		return (OTSClientError{}.Set("delete_row_item should be an instance of [DeleteRowInBatchWriteRowRequest, *DeleteRowInBatchWriteRowRequest or OTSDeleteRowItem], not %v", reflect.TypeOf(delete_row_item)))
	}

	return nil
}

func _make_batch_write_row(proto *BatchWriteRowRequest, batch_list interface{}) error {
	switch batch_list.(type) {
	case []TableInBatchWriteRowRequest:
		list_len := len(batch_list.([]TableInBatchWriteRowRequest))
		proto.Tables = make([]*TableInBatchWriteRowRequest, list_len)
		for i, v := range batch_list.([]TableInBatchWriteRowRequest) {
			table_item := new(TableInBatchWriteRowRequest)
			// table_name
			table_name := _get_unicode(*v.TableName)
			table_item.TableName = new(string)
			*table_item.TableName = table_name

			// PutRows
			table_item.PutRows = make([]*PutRowInBatchWriteRowRequest, len(v.PutRows))
			for i1, v1 := range v.PutRows {
				put_rows_item := new(PutRowInBatchWriteRowRequest)
				err := _make_put_row_item(put_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.PutRows[i1] = put_rows_item
			}

			// UpdateRows
			table_item.UpdateRows = make([]*UpdateRowInBatchWriteRowRequest, len(v.UpdateRows))
			for i1, v1 := range v.UpdateRows {
				update_rows_item := new(UpdateRowInBatchWriteRowRequest)
				err := _make_update_row_item(update_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.UpdateRows[i1] = update_rows_item
			}

			// DeleteRows
			table_item.DeleteRows = make([]*DeleteRowInBatchWriteRowRequest, len(v.DeleteRows))
			for i1, v1 := range v.DeleteRows {
				delete_rows_item := new(DeleteRowInBatchWriteRowRequest)
				err := _make_delete_row_item(delete_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.DeleteRows[i1] = delete_rows_item
			}
			proto.Tables[i] = table_item
		}

	case []*TableInBatchWriteRowRequest:
		list_len := len(batch_list.([]*TableInBatchWriteRowRequest))
		proto.Tables = make([]*TableInBatchWriteRowRequest, list_len)
		for i, v := range batch_list.([]*TableInBatchWriteRowRequest) {
			table_item := new(TableInBatchWriteRowRequest)
			// table_name
			table_name := _get_unicode(*v.TableName)
			table_item.TableName = new(string)
			*table_item.TableName = table_name

			// PutRows
			table_item.PutRows = make([]*PutRowInBatchWriteRowRequest, len(v.PutRows))
			for i1, v1 := range v.PutRows {
				put_rows_item := new(PutRowInBatchWriteRowRequest)
				err := _make_put_row_item(put_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.PutRows[i1] = put_rows_item
			}

			// UpdateRows
			table_item.UpdateRows = make([]*UpdateRowInBatchWriteRowRequest, len(v.UpdateRows))
			for i1, v1 := range v.UpdateRows {
				update_rows_item := new(UpdateRowInBatchWriteRowRequest)
				err := _make_update_row_item(update_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.UpdateRows[i1] = update_rows_item
			}

			// DeleteRows
			table_item.DeleteRows = make([]*DeleteRowInBatchWriteRowRequest, len(v.DeleteRows))
			for i1, v1 := range v.DeleteRows {
				delete_rows_item := new(DeleteRowInBatchWriteRowRequest)
				err := _make_delete_row_item(delete_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.DeleteRows[i1] = delete_rows_item
			}
			proto.Tables[i] = table_item
		}

	case OTSBatchWriteRowRequest:
		list_len := len(batch_list.(OTSBatchWriteRowRequest))
		proto.Tables = make([]*TableInBatchWriteRowRequest, list_len)
		for i, v := range batch_list.(OTSBatchWriteRowRequest) {
			table_item := new(TableInBatchWriteRowRequest)
			// table_name
			table_name := _get_unicode(v.TableName)
			table_item.TableName = new(string)
			*table_item.TableName = table_name

			// PutRows
			table_item.PutRows = make([]*PutRowInBatchWriteRowRequest, len(v.PutRows))
			for i1, v1 := range v.PutRows {
				put_rows_item := new(PutRowInBatchWriteRowRequest)
				err := _make_put_row_item(put_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.PutRows[i1] = put_rows_item
			}

			// UpdateRows
			table_item.UpdateRows = make([]*UpdateRowInBatchWriteRowRequest, len(v.UpdateRows))
			for i1, v1 := range v.UpdateRows {
				update_rows_item := new(UpdateRowInBatchWriteRowRequest)
				err := _make_update_row_item(update_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.UpdateRows[i1] = update_rows_item
			}

			// DeleteRows
			table_item.DeleteRows = make([]*DeleteRowInBatchWriteRowRequest, len(v.DeleteRows))
			for i1, v1 := range v.DeleteRows {
				delete_rows_item := new(DeleteRowInBatchWriteRowRequest)
				err := _make_delete_row_item(delete_rows_item, v1)
				if err != nil {
					return err
				}
				table_item.DeleteRows[i1] = delete_rows_item
			}
			proto.Tables[i] = table_item
		}

	default:
		return (OTSClientError{}.Set("batch_list should be an instance of [[]TableInBatchWriteRowRequest, []*TableInBatchWriteRowRequest or OTSBatchWriteRowRequest], not %v", reflect.TypeOf(batch_list)))
	}

	return nil
}

// simple testing for above encoder functions
// just for tobyzxj
// forgot it
func TestEncoder() {
	fmt.Println("Encoder test...")

	// func _make_repeated_column_names(proto []string, columns_to_get []string) error
	// var proto []string
	// _make_repeated_column_names(&proto, []string{})
	// fmt.Println(proto)
	// _make_repeated_column_names(&proto, []string{"toby", "allen"})
	// fmt.Println(proto)
	// _make_repeated_column_names(&proto, []string{"toby1", "allen2"})
	// fmt.Println(proto)

	// func _make_column_value(proto *ColumnValue, value interface{})
	// proto := ColumnValue{}
	// _make_column_value(&proto, "tobyzxj")
	// fmt.Println(proto, proto.GetVString())
	// _make_column_value(&proto, 123)
	// fmt.Println(proto, proto.GetVInt())

	// func _get_column_type(type_str string) ColumnType
	// proto := _get_column_type("INF_MIN")
	// fmt.Println(proto)
	// proto = _get_column_type("tobyzxj")
	// fmt.Println(proto)

	// func _make_condition(proto *Condition, condition interface{}) error
	// [1]
	// proto := Condition{}
	// fmt.Println(proto)
	// prowExistence := new(RowExistenceExpectation)
	// *prowExistence = RowExistenceExpectation_IGNORE
	// _make_condition(&proto, Condition{RowExistence: prowExistence})
	// fmt.Println(proto)
	//
	// [2]
	// proto := Condition{}
	// fmt.Println(proto)
	// prowExistence := new(RowExistenceExpectation)
	// *prowExistence = RowExistenceExpectation_IGNORE
	// _make_condition(&proto, &Condition{RowExistence: prowExistence})
	// fmt.Println(proto)
	//
	// [3]
	// proto := Condition{}
	// fmt.Println(proto)
	// _make_condition(&proto, "IGNORE")
	// fmt.Println(proto)
	//
	// [4]
	// proto := Condition{}
	// fmt.Println(proto)
	// _make_condition(&proto, "NO_IGNORE")
	// fmt.Println(proto)

	// func _make_column_schema(proto *ColumnSchema, schema_tuple ColumnSchema) error
	// proto := new(ColumnSchema)
	// fmt.Println(proto)
	// column_type := new(ColumnType)
	// *column_type = ColumnType_INF_MIN
	// schema_tuple := ColumnSchema{
	// 	Name: NewString("toby"),
	// 	Type: column_type,
	// }
	// _make_column_schema(proto, schema_tuple)
	// fmt.Println(proto)

	// func _make_schemas_with_list(proto *[]*ColumnSchema, schema_list interface{}) error
	// [1]
	// var proto = new([]*ColumnSchema)
	// fmt.Println(proto)
	// column_type := new(ColumnType)
	// *column_type = ColumnType_INF_MIN
	// schema_list := []ColumnSchema{
	// 	{Name: NewString("toby1"), Type: column_type},
	// 	{Name: NewString("toby2"), Type: column_type},
	// 	{Name: NewString("toby3"), Type: column_type},
	// }
	// err := _make_schemas_with_list(proto, schema_list)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }
	// fmt.Println(proto)
	//
	// func _make_schemas_with_list(proto *[]*ColumnSchema, schema_list interface{}) error
	// [2]
	// var proto = new([]*ColumnSchema)
	// fmt.Println(proto)
	// column_type := new(ColumnType)
	// *column_type = ColumnType_INF_MIN
	// schema_list := []TupleString{
	// 	{"toby1", "STRING"},
	// 	{"toby2", "INTEGER"},
	// 	{"toby3", "STRING"},
	// }
	// err := _make_schemas_with_list(proto, schema_list)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }
	// fmt.Println(proto)

	// func _make_columns_with_dict(proto *[]*Column, column_dict interface{}) error
	// [1]
	// proto := new([]*Column)
	// fmt.Println(proto)
	// column_value := new(ColumnValue)
	// _make_column_value(column_value, 123)
	// column_dict := []Column{
	// 	{Name: NewString("tobyzxj1"), Value: column_value},
	// 	{Name: NewString("tobyzxj2"), Value: column_value},
	// }
	// _make_columns_with_dict(proto, column_dict)
	// fmt.Println(proto)
	//
	// [2]
	// proto := new([]*Column)
	// fmt.Println(proto)
	// column_value := new(ColumnValue)
	// _make_column_value(column_value, 123)
	// column_dict := []*Column{
	// 	&Column{
	// 		Name:  NewString("tobyzxj1"),
	// 		Value: column_value,
	// 	},
	// 	&Column{
	// 		Name:  NewString("tobyzxj2"),
	// 		Value: column_value,
	// 	},
	// }
	// _make_columns_with_dict(proto, column_dict)
	// fmt.Println(proto)
	//
	// [3]
	// proto := new([]*Column)
	// fmt.Println(proto)
	// column_dict := DictString{
	// 	"tobyzxj1": 123,
	// 	"tobyzxj2": "i'm here",
	// }
	// _make_columns_with_dict(proto, column_dict)
	// fmt.Println(proto)

	// func _make_update_of_attribute_columns_with_dict(proto *[]*ColumnUpdate, column_dict interface{}) error
	// [1]
	// proto := new([]*ColumnUpdate)
	// fmt.Println(proto)
	// column_update_type1 := new(OperationType)
	// *column_update_type1 = OperationType_PUT
	// column_update_type2 := new(OperationType)
	// *column_update_type2 = OperationType_DELETE
	// column_update_value := new(ColumnValue)
	// _make_column_value(column_update_value, 123)
	// column_dict := []ColumnUpdate{
	// 	{Type: column_update_type1, Name: NewString("tobyzxj1"), Value: column_update_value},
	// 	{Type: column_update_type2, Name: NewString("tobyzxj2")},
	// }
	// _make_update_of_attribute_columns_with_dict(proto, column_dict)
	// fmt.Println(proto)
	//
	// [2]
	// proto := new([]*ColumnUpdate)
	// fmt.Println(proto)
	// column_update_type1 := new(OperationType)
	// *column_update_type1 = OperationType_PUT
	// column_update_type2 := new(OperationType)
	// *column_update_type2 = OperationType_DELETE
	// column_update_value := new(ColumnValue)
	// _make_column_value(column_update_value, 123)
	// column_dict := []*ColumnUpdate{
	// 	&ColumnUpdate{
	// 		Type:  column_update_type1,
	// 		Name:  NewString("tobyzxj1"),
	// 		Value: column_update_value,
	// 	},
	// 	&ColumnUpdate{
	// 		Type: column_update_type2,
	// 		Name: NewString("tobyzxj2"),
	// 	},
	// }
	// _make_update_of_attribute_columns_with_dict(proto, column_dict)
	// fmt.Println(proto)
	//
	// [3]
	// proto := new([]*ColumnUpdate)
	// fmt.Println(proto)
	// column_dict := DictString{
	// 	"PUT": DictString{
	// 		"put_name1": "value1",
	// 		"put_name2": 123,
	// 	},
	// 	"DELETE": []string{
	// 		"del_name1", "del_name2",
	// 	},
	// }
	// _make_update_of_attribute_columns_with_dict(proto, column_dict)
	// fmt.Println(proto)

	// func _make_table_meta(proto *TableMeta, table_meta interface{}) error
	// [1]
	// proto := new(TableMeta)
	// fmt.Println(proto)
	// column_interge_type := new(ColumnType)
	// *column_interge_type = ColumnType_INTEGER
	// column_string_type := new(ColumnType)
	// *column_string_type = ColumnType_STRING
	// table_meta := TableMeta{
	// 	TableName: NewString("tablename"),
	// 	PrimaryKey: []*ColumnSchema{
	// 		&ColumnSchema{
	// 			Name: NewString("tobyzxj1"),
	// 			Type: column_interge_type,
	// 		},
	// 		&ColumnSchema{
	// 			Name: NewString("tobyzxj2"),
	// 			Type: column_string_type,
	// 		},
	// 	},
	// }
	// _make_table_meta(proto, table_meta)
	// fmt.Println(proto)
	//
	// [2]
	// proto := new(TableMeta)
	// fmt.Println(proto)
	// table_meta := TupleString{
	// 	"tablename", []TupleString{
	// 		{"PK1", "INTEGER"},
	// 		{"PK2", "STRING"},
	// 	},
	// }
	// _make_table_meta(proto, table_meta)
	// fmt.Println(proto)
	//
	// [3]
	// proto := new(TableMeta)
	// fmt.Println(proto)
	// table_meta := OTSTableMeta{
	// 	"tablename",
	// 	OTSSchemaOfPrimaryKey{
	// 		"PK1": "INTEGER",
	// 		"PK2": "STRING",
	// 	},
	// }
	// _make_table_meta(proto, table_meta)
	// fmt.Println(proto)

	// func _make_capacity_unit(proto *CapacityUnit, capacity_unit interface{}) error
	// proto := new(CapacityUnit)
	// fmt.Println(proto)
	// capacity_unit := OTSCapacityUnit{100, 100}
	// _make_capacity_unit(proto, capacity_unit)
	// fmt.Println(proto)

	// func _make_reserved_throughput(proto *ReservedThroughput, reserved_throughput interface{}) error
	// proto := new(ReservedThroughput)
	// fmt.Println(proto)
	// reserved_throughput := OTSReservedThroughput{
	// 	OTSCapacityUnit{100, 100},
	// }
	// _make_reserved_throughput(proto, reserved_throughput)
	// fmt.Println(proto)

	// func _make_update_capacity_unit(proto *CapacityUnit, capacity_unit interface{}) error
	// proto := new(CapacityUnit)
	// fmt.Println(proto)
	// capacity_unit := OTSCapacityUnit{Read: 100}
	// _make_update_capacity_unit(proto, capacity_unit)
	// fmt.Println(proto)

	// func _make_update_reserved_throughput(proto *ReservedThroughput, reserved_throughput interface{}) error
	// proto := new(ReservedThroughput)
	// fmt.Println(proto)
	// reserved_throughput := OTSReservedThroughput{
	// 	OTSCapacityUnit{Read: 100, Write: 100},
	// }
	// _make_update_reserved_throughput(proto, reserved_throughput)
	// fmt.Println(proto)

	// func _make_batch_get_row(proto *BatchGetRowRequest, batch_list interface{}) error
	// [1]
	// proto := new(BatchGetRowRequest)
	// fmt.Println(proto)
	// column := new([]*Column)
	// column_value := new(ColumnValue)
	// _make_column_value(column_value, 123)
	// column_dict := DictString{
	// 	"PK1": 123,
	// 	"PK2": "STRING",
	// }
	// _make_columns_with_dict(column, column_dict)
	// // fmt.Println("columns:", column)
	//
	// batch_list := []TableInBatchGetRowRequest{
	// 	{
	// 		TableName: NewString("table_name0"),
	// 		Rows: []*RowInBatchGetRowRequest{
	// 			&RowInBatchGetRowRequest{
	// 				PrimaryKey: *column,
	// 			},
	// 		},
	// 		ColumnsToGet: []string{
	// 			"column_name0",
	// 			"column_name1",
	// 		},
	// 	},
	// 	{
	// 		TableName: NewString("table_name1"),
	// 		Rows: []*RowInBatchGetRowRequest{
	// 			&RowInBatchGetRowRequest{
	// 				PrimaryKey: *column,
	// 			},
	// 		},
	// 		ColumnsToGet: []string{
	// 			"column_name0",
	// 			"column_name1",
	// 		},
	// 	},
	// }
	// // fmt.Println("batch_list:", batch_list)
	// _make_batch_get_row(proto, batch_list)
	// fmt.Println(proto)
	//
	// [2]
	// proto := new(BatchGetRowRequest)
	// fmt.Println(proto)
	// column := new([]*Column)
	// column_value := new(ColumnValue)
	// _make_column_value(column_value, 123)
	// column_dict := DictString{
	// 	"PK1": 123,
	// 	"PK2": "STRING",
	// }
	// _make_columns_with_dict(column, column_dict)
	// // fmt.Println("columns:", column)

	// batch_list := []*TableInBatchGetRowRequest{
	// 	&TableInBatchGetRowRequest{
	// 		TableName: NewString("table_name0"),
	// 		Rows: []*RowInBatchGetRowRequest{
	// 			&RowInBatchGetRowRequest{
	// 				PrimaryKey: *column,
	// 			},
	// 		},
	// 		ColumnsToGet: []string{
	// 			"column_name0",
	// 			"column_name1",
	// 		},
	// 	},
	// 	&TableInBatchGetRowRequest{
	// 		TableName: NewString("table_name1"),
	// 		Rows: []*RowInBatchGetRowRequest{
	// 			&RowInBatchGetRowRequest{
	// 				PrimaryKey: *column,
	// 			},
	// 		},
	// 		ColumnsToGet: []string{
	// 			"column_name0",
	// 			"column_name1",
	// 		},
	// 	},
	// }
	// // fmt.Println("batch_list:", batch_list)
	// _make_batch_get_row(proto, batch_list)
	// fmt.Println(proto)
	//
	// [3]
	// proto := new(BatchGetRowRequest)
	// fmt.Println(proto)

	// batch_list := OTSBatchGetRowRequest{
	// 	{
	// 		// TableName
	// 		TableName: "table_name0",
	// 		// PrimaryKey
	// 		Rows: OTSPrimaryKeyRows{
	// 			{"gid": 1, "uid": 101},
	// 			{"gid": 2, "uid": 202},
	// 			{"gid": 3, "uid": 303},
	// 		},
	// 		// ColumnsToGet
	// 		ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
	// 	},
	// 	{
	// 		// TableName
	// 		TableName: "notExistTable",
	// 		// PrimaryKey
	// 		Rows: OTSPrimaryKeyRows{
	// 			{"gid": 1, "uid": 101},
	// 			{"gid": 2, "uid": 202},
	// 			{"gid": 3, "uid": 303},
	// 		},
	// 		// ColumnsToGet
	// 		ColumnsToGet: OTSColumnsToGet{"name", "address", "mobile", "age"},
	// 	},
	// }
	//
	// batch_list = OTSBatchGetRowRequest{} // test none
	// fmt.Println("batch_list:", batch_list)
	// _make_batch_get_row(proto, batch_list)
	// fmt.Println(proto)

	// func _make_put_row_item(proto *PutRowInBatchWriteRowRequest, put_row_item interface{}) error
	// [1]
	// proto := new(PutRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// condition := new(Condition)
	// prowExistence := new(RowExistenceExpectation)
	// *prowExistence = RowExistenceExpectation_IGNORE
	// _make_condition(condition, Condition{RowExistence: prowExistence})
	// // fmt.Println(condition)
	// pk_column := new([]*Column)
	// pk_column_dict := DictString{
	// 	"gid": 2,
	// 	"uid": 202,
	// }
	// _make_columns_with_dict(pk_column, pk_column_dict)
	// attr_column := new([]*Column)
	// attr_column_dict := DictString{
	// 	"name":    "李四",
	// 	"address": "中国某地",
	// 	"age":     20,
	// }
	// _make_columns_with_dict(attr_column, attr_column_dict)
	// put_row_item := PutRowInBatchWriteRowRequest{
	// 	Condition:        condition,
	// 	PrimaryKey:       *pk_column,
	// 	AttributeColumns: *attr_column,
	// }
	// _make_put_row_item(proto, put_row_item)
	// fmt.Println(proto)
	//
	// [2]
	// proto := new(PutRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// condition := new(Condition)
	// prowExistence := new(RowExistenceExpectation)
	// *prowExistence = RowExistenceExpectation_IGNORE
	// _make_condition(condition, Condition{RowExistence: prowExistence})
	// // fmt.Println(condition)
	// pk_column := new([]*Column)
	// pk_column_dict := DictString{
	// 	"gid": 2,
	// 	"uid": 202,
	// }
	// _make_columns_with_dict(pk_column, pk_column_dict)
	// attr_column := new([]*Column)
	// attr_column_dict := DictString{
	// 	"name":    "李四",
	// 	"address": "中国某地",
	// 	"age":     20,
	// }
	// _make_columns_with_dict(attr_column, attr_column_dict)
	// put_row_item := &PutRowInBatchWriteRowRequest{
	// 	Condition:        condition,
	// 	PrimaryKey:       *pk_column,
	// 	AttributeColumns: *attr_column,
	// }
	// _make_put_row_item(proto, put_row_item)
	// fmt.Println(proto)
	//
	// [3]
	// proto := new(PutRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// put_row_item := OTSPutRowItem{
	// 	Condition: OTSCondition_EXPECT_EXIST,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 2,
	// 		"uid": 202,
	// 	},
	// 	AttributeColumns: OTSAttribute{
	// 		"name":    "李四",
	// 		"address": "中国某地",
	// 		"age":     20,
	// 	},
	// }
	// _make_put_row_item(proto, put_row_item)
	// fmt.Println(proto)

	// func _make_update_row_item(proto *UpdateRowInBatchWriteRowRequest, update_row_item interface{}) error
	//
	// [1]
	// proto := new(UpdateRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// condition := new(Condition)
	// prowExistence := new(RowExistenceExpectation)
	// *prowExistence = RowExistenceExpectation_IGNORE
	// _make_condition(condition, Condition{RowExistence: prowExistence})
	// // fmt.Println(condition)
	// pk_column := new([]*Column)
	// pk_column_dict := DictString{
	// 	"gid": 2,
	// 	"uid": 202,
	// }
	// _make_columns_with_dict(pk_column, pk_column_dict)
	// attr_column := new([]*ColumnUpdate)
	// attr_column_dict := DictString{
	// 	"PUT": DictString{
	// 		"name":    "李四",
	// 		"address": "中国某地",
	// 	},
	// 	"DELETE": []string{
	// 		"age",
	// 	},
	// }
	// _make_update_of_attribute_columns_with_dict(attr_column, attr_column_dict)
	// update_row_item := UpdateRowInBatchWriteRowRequest{
	// 	Condition:        condition,
	// 	PrimaryKey:       *pk_column,
	// 	AttributeColumns: *attr_column,
	// }
	// _make_update_row_item(proto, update_row_item)
	// fmt.Println(proto)
	//
	// [2]
	//
	// proto := new(UpdateRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// condition := new(Condition)
	// prowExistence := new(RowExistenceExpectation)
	// *prowExistence = RowExistenceExpectation_IGNORE
	// _make_condition(condition, Condition{RowExistence: prowExistence})
	// // fmt.Println(condition)
	// pk_column := new([]*Column)
	// pk_column_dict := DictString{
	// 	"gid": 2,
	// 	"uid": 202,
	// }
	// _make_columns_with_dict(pk_column, pk_column_dict)
	// attr_column := new([]*ColumnUpdate)
	// attr_column_dict := DictString{
	// 	"PUT": DictString{
	// 		"name":    "李四",
	// 		"address": "中国某地",
	// 	},
	// 	"DELETE": []string{
	// 		"age",
	// 	},
	// }
	// _make_update_of_attribute_columns_with_dict(attr_column, attr_column_dict)
	// update_row_item := &UpdateRowInBatchWriteRowRequest{
	// 	Condition:        condition,
	// 	PrimaryKey:       *pk_column,
	// 	AttributeColumns: *attr_column,
	// }
	// _make_update_row_item(proto, update_row_item)
	// fmt.Println(proto)
	//
	// [3]
	// proto := new(UpdateRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// update_row_item := OTSUpdateRowItem{
	// 	Condition: OTSCondition_EXPECT_EXIST,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 2,
	// 		"uid": 202,
	// 	},
	// 	UpdateOfAttributeColumns: OTSUpdateOfAttribute{
	// 		"PUT": DictString{
	// 			"name":    "李四",
	// 			"address": "中国某地",
	// 		},
	// 		"DELETE": []string{
	// 			"age",
	// 		},
	// 	},
	// }
	// _make_update_row_item(proto, update_row_item)
	// fmt.Println(proto)
	//
	//
	// [4]
	// proto := new(UpdateRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// update_row_item := OTSUpdateRowItem{
	// 	Condition: OTSCondition_EXPECT_EXIST,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 2,
	// 		"uid": 202,
	// 	},
	// 	UpdateOfAttributeColumns: OTSUpdateOfAttribute{
	// 		OTSOperationType_PUT: OTSColumnsToPut{
	// 			"name":    "李四",
	// 			"address": "中国某地",
	// 		},
	// 		OTSOperationType_DELETE: OTSColumnsToDelete{
	// 			"age",
	// 		},
	// 	},
	// }
	// _make_update_row_item(proto, update_row_item)
	// fmt.Println(proto)

	// func _make_delete_row_item(proto *DeleteRowInBatchWriteRowRequest, delete_row_item interface{}) error
	//
	// [1]
	// proto := new(DeleteRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// condition := new(Condition)
	// prowExistence := new(RowExistenceExpectation)
	// *prowExistence = RowExistenceExpectation_IGNORE
	// _make_condition(condition, Condition{RowExistence: prowExistence})
	// // fmt.Println(condition)
	// pk_column := new([]*Column)
	// pk_column_dict := DictString{
	// 	"gid": 2,
	// 	"uid": 202,
	// }
	// _make_columns_with_dict(pk_column, pk_column_dict)
	// delete_row_item := DeleteRowInBatchWriteRowRequest{
	// 	Condition:  condition,
	// 	PrimaryKey: *pk_column,
	// }
	// _make_delete_row_item(proto, delete_row_item)
	// fmt.Println(proto)
	//
	// [2]
	// proto := new(DeleteRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// condition := new(Condition)
	// prowExistence := new(RowExistenceExpectation)
	// *prowExistence = RowExistenceExpectation_IGNORE
	// _make_condition(condition, Condition{RowExistence: prowExistence})
	// // fmt.Println(condition)
	// pk_column := new([]*Column)
	// pk_column_dict := DictString{
	// 	"gid": 2,
	// 	"uid": 202,
	// }
	// _make_columns_with_dict(pk_column, pk_column_dict)
	// delete_row_item := &DeleteRowInBatchWriteRowRequest{
	// 	Condition:  condition,
	// 	PrimaryKey: *pk_column,
	// }
	// _make_delete_row_item(proto, delete_row_item)
	// fmt.Println(proto)
	//
	// [3]
	// proto := new(DeleteRowInBatchWriteRowRequest)
	// fmt.Println(proto)
	// delete_row_item := OTSDeleteRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 2,
	// 		"uid": 202,
	// 	},
	// }
	// _make_delete_row_item(proto, delete_row_item)
	// fmt.Println(proto)

	// func _make_batch_write_row(proto *BatchWriteRowRequest, batch_list interface{}) error
	//
	// [1]
	// proto := new(BatchWriteRowRequest)
	// fmt.Println(proto)
	// // put row
	// put_row := new(PutRowInBatchWriteRowRequest)
	// put_row_item := OTSPutRowItem{
	// 	Condition: OTSCondition_EXPECT_NOT_EXIST,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 2,
	// 		"uid": 202,
	// 	},
	// 	AttributeColumns: OTSAttribute{
	// 		"name":    "李四",
	// 		"address": "中国某地",
	// 		"age":     20,
	// 	},
	// }
	// _make_put_row_item(put_row, put_row_item)
	// // fmt.Println("put_row:", put_row)
	// // update_row
	// update_row := new(UpdateRowInBatchWriteRowRequest)
	// update_row_item := OTSUpdateRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 3,
	// 		"uid": 303,
	// 	},
	// 	UpdateOfAttributeColumns: OTSUpdateOfAttribute{
	// 		OTSOperationType_PUT: OTSColumnsToPut{
	// 			"name":    "李四",
	// 			"address": "中国某地",
	// 		},
	// 		OTSOperationType_DELETE: OTSColumnsToDelete{
	// 			"mobile", "age",
	// 		},
	// 	},
	// }
	// _make_update_row_item(update_row, update_row_item)
	// // fmt.Println("update_row:", update_row)
	// // delete_row
	// delete_row := new(DeleteRowInBatchWriteRowRequest)
	// delete_row_item := OTSDeleteRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 4,
	// 		"uid": 404,
	// 	},
	// }
	// _make_delete_row_item(delete_row, delete_row_item)
	// // fmt.Println("delete_row:", delete_row)
	// batch_list := []TableInBatchWriteRowRequest{
	// 	{
	// 		TableName: NewString("myTable"),
	// 		PutRows: []*PutRowInBatchWriteRowRequest{
	// 			put_row,
	// 		},
	// 		UpdateRows: []*UpdateRowInBatchWriteRowRequest{
	// 			update_row,
	// 		},
	// 		DeleteRows: []*DeleteRowInBatchWriteRowRequest{
	// 			delete_row,
	// 		},
	// 	},
	// 	{
	// 		TableName: NewString("notExistTable"),
	// 		PutRows: []*PutRowInBatchWriteRowRequest{
	// 			put_row,
	// 		},
	// 		UpdateRows: []*UpdateRowInBatchWriteRowRequest{
	// 			update_row,
	// 		},
	// 		DeleteRows: []*DeleteRowInBatchWriteRowRequest{
	// 			delete_row,
	// 		},
	// 	},
	// }
	// _make_batch_write_row(proto, batch_list)
	// fmt.Println(proto)
	//
	// [2]
	// proto := new(BatchWriteRowRequest)
	// fmt.Println(proto)
	// // put row
	// put_row := new(PutRowInBatchWriteRowRequest)
	// put_row_item := OTSPutRowItem{
	// 	Condition: OTSCondition_EXPECT_NOT_EXIST,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 2,
	// 		"uid": 202,
	// 	},
	// 	AttributeColumns: OTSAttribute{
	// 		"name":    "李四",
	// 		"address": "中国某地",
	// 		"age":     20,
	// 	},
	// }
	// _make_put_row_item(put_row, put_row_item)
	// // fmt.Println("put_row:", put_row)
	// // update_row
	// update_row := new(UpdateRowInBatchWriteRowRequest)
	// update_row_item := OTSUpdateRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 3,
	// 		"uid": 303,
	// 	},
	// 	UpdateOfAttributeColumns: OTSUpdateOfAttribute{
	// 		OTSOperationType_PUT: OTSColumnsToPut{
	// 			"name":    "李四",
	// 			"address": "中国某地",
	// 		},
	// 		OTSOperationType_DELETE: OTSColumnsToDelete{
	// 			"mobile", "age",
	// 		},
	// 	},
	// }
	// _make_update_row_item(update_row, update_row_item)
	// // fmt.Println("update_row:", update_row)
	// // delete_row
	// delete_row := new(DeleteRowInBatchWriteRowRequest)
	// delete_row_item := OTSDeleteRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 4,
	// 		"uid": 404,
	// 	},
	// }
	// _make_delete_row_item(delete_row, delete_row_item)
	// // fmt.Println("delete_row:", delete_row)
	// batch_list := []*TableInBatchWriteRowRequest{
	// 	&TableInBatchWriteRowRequest{
	// 		TableName: NewString("myTable"),
	// 		PutRows: []*PutRowInBatchWriteRowRequest{
	// 			put_row,
	// 		},
	// 		UpdateRows: []*UpdateRowInBatchWriteRowRequest{
	// 			update_row,
	// 		},
	// 		DeleteRows: []*DeleteRowInBatchWriteRowRequest{
	// 			delete_row,
	// 		},
	// 	},
	// 	&TableInBatchWriteRowRequest{
	// 		TableName: NewString("notExistTable"),
	// 		PutRows: []*PutRowInBatchWriteRowRequest{
	// 			put_row,
	// 		},
	// 		UpdateRows: []*UpdateRowInBatchWriteRowRequest{
	// 			update_row,
	// 		},
	// 		DeleteRows: []*DeleteRowInBatchWriteRowRequest{
	// 			delete_row,
	// 		},
	// 	},
	// }
	// _make_batch_write_row(proto, batch_list)
	// fmt.Println(proto)
	//
	// [3]
	// proto := new(BatchWriteRowRequest)
	// fmt.Println(proto)
	// // put row
	// put_row_item := OTSPutRowItem{
	// 	Condition: OTSCondition_EXPECT_NOT_EXIST,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 2,
	// 		"uid": 202,
	// 	},
	// 	AttributeColumns: OTSAttribute{
	// 		"name":    "李四",
	// 		"address": "中国某地",
	// 		"age":     20,
	// 	},
	// }
	// // update_row
	// update_row_item := OTSUpdateRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 3,
	// 		"uid": 303,
	// 	},
	// 	UpdateOfAttributeColumns: OTSUpdateOfAttribute{
	// 		OTSOperationType_PUT: OTSColumnsToPut{
	// 			"name":    "李四",
	// 			"address": "中国某地",
	// 		},
	// 		OTSOperationType_DELETE: OTSColumnsToDelete{
	// 			"mobile", "age",
	// 		},
	// 	},
	// }

	// // delete_row
	// delete_row_item := OTSDeleteRowItem{
	// 	Condition: OTSCondition_IGNORE,
	// 	PrimaryKey: OTSPrimaryKey{
	// 		"gid": 4,
	// 		"uid": 404,
	// 	},
	// }

	// batch_list := OTSBatchWriteRowRequest{
	// 	{
	// 		TableName: "myTable",
	// 		PutRows: OTSPutRows{
	// 			put_row_item,
	// 		},
	// 		UpdateRows: OTSUpdateRows{
	// 			update_row_item,
	// 		},
	// 		DeleteRows: OTSDeleteRows{
	// 			delete_row_item,
	// 		},
	// 	},
	// 	{
	// 		TableName: "notExistTable",
	// 		PutRows: OTSPutRows{
	// 			put_row_item,
	// 		},
	// 		UpdateRows: OTSUpdateRows{
	// 			update_row_item,
	// 		},
	// 		DeleteRows: OTSDeleteRows{
	// 			delete_row_item,
	// 		},
	// 	},
	// }
	// _make_batch_write_row(proto, batch_list)
	// fmt.Println(proto)
}

func _encode_create_table(table_meta *OTSTableMeta, reserved_throughput *OTSReservedThroughput) (req *CreateTableRequest, err error) {
	proto := new(CreateTableRequest)
	proto.TableMeta = new(TableMeta)
	proto.ReservedThroughput = new(ReservedThroughput)
	err = _make_table_meta(proto.TableMeta, *table_meta)
	if err != nil {
		return nil, err
	}
	err = _make_reserved_throughput(proto.ReservedThroughput, *reserved_throughput)
	if err != nil {
		return nil, err
	}

	return proto, nil
}

func _encode_delete_table(table_name string) (req *DeleteTableRequest, err error) {
	proto := new(DeleteTableRequest)
	proto.TableName = NewString(table_name)

	return proto, nil
}

func _encode_list_table() (req *ListTableRequest, err error) {
	proto := new(ListTableRequest)

	return proto, nil
}

func _encode_update_table(table_name string, reserved_throughput *OTSReservedThroughput) (req *UpdateTableRequest, err error) {
	proto := new(UpdateTableRequest)
	proto.TableName = NewString(table_name)
	proto.ReservedThroughput = new(ReservedThroughput)
	err = _make_update_reserved_throughput(proto.ReservedThroughput, *reserved_throughput)
	if err != nil {
		return nil, err
	}

	return proto, nil
}

func _encode_describe_table(table_name string) (req *DescribeTableRequest, err error) {
	proto := new(DescribeTableRequest)
	proto.TableName = NewString(table_name)

	return proto, nil
}

func _encode_get_row(table_name string, primary_key *OTSPrimaryKey, columns_to_get *OTSColumnsToGet) (req *GetRowRequest, err error) {
	proto := new(GetRowRequest)
	proto.TableName = NewString(table_name)
	_primary_key := new([]*Column)
	err = _make_columns_with_dict(_primary_key, DictString(*primary_key))
	if err != nil {
		return nil, err
	}
	proto.PrimaryKey = *_primary_key
	if columns_to_get != nil {
		_columns_to_get := new([]string)
		err = _make_repeated_column_names(_columns_to_get, []string(*columns_to_get))
		if err != nil {
			return nil, err
		}
		proto.ColumnsToGet = *_columns_to_get
	} else {
		proto.ColumnsToGet = nil
	}

	return proto, nil
}

func _encode_put_row(table_name string, condition string, primary_key *OTSPrimaryKey, attribute_columns *OTSAttribute) (req *PutRowRequest, err error) {
	proto := new(PutRowRequest)
	proto.TableName = NewString(table_name)
	proto.Condition = new(Condition)
	err = _make_condition(proto.Condition, condition)
	if err != nil {
		return nil, err
	}
	_primary_key := new([]*Column)
	err = _make_columns_with_dict(_primary_key, DictString(*primary_key))
	if err != nil {
		return nil, err
	}
	proto.PrimaryKey = *_primary_key
	_attribute_columns := new([]*Column)
	err = _make_columns_with_dict(_attribute_columns, DictString(*attribute_columns))
	if err != nil {
		return nil, err
	}
	proto.AttributeColumns = *_attribute_columns

	return proto, nil
}

func _encode_update_row() {

}

func _encode_delete_row() {

}

func _encode_batch_get_row() {

}

func _encode_batch_write_row() {

}

func _encode_get_range() {

}

// request encode for ots2
func EncodeRequest(api_name string, args ...interface{}) (req []reflect.Value, err error) {
	if _, ok := api_encode_map[api_name]; !ok {
		return nil, (OTSClientError{}.Set("No PB encode method for API %s", api_name))
	}

	req, err = api_encode_map.Call(api_name, args...)
	if err != nil {
		return nil, err
	}

	return req, nil
}
