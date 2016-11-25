# -*- coding: utf8 -*-

import google.protobuf.text_format as text_format

from ots2.error import *
from ots2.metadata import *
import ots2.protobuf.ots_protocol_2_pb2 as pb2

INT32_MAX = 2147483647
INT32_MIN = -2147483648

class OTSProtoBufferEncoder:

    def __init__(self, encoding):
        self.encoding = encoding

        self.api_encode_map = {
            'CreateTable'       : self._encode_create_table, 
            'DeleteTable'       : self._encode_delete_table, 
            'ListTable'         : self._encode_list_table,
            'UpdateTable'       : self._encode_update_table,
            'DescribeTable'     : self._encode_describe_table,
            'GetRow'            : self._encode_get_row,
            'PutRow'            : self._encode_put_row,
            'UpdateRow'         : self._encode_update_row,
            'DeleteRow'         : self._encode_delete_row,
            'BatchGetRow'       : self._encode_batch_get_row,
            'BatchWriteRow'     : self._encode_batch_write_row,
            'GetRange'          : self._encode_get_range
        }

    def _get_unicode(self, value):
        if isinstance(value, str):
            return value.decode(self.encoding)
        elif isinstance(value, unicode):
            return value
        else:
            raise OTSClientError(
                "expect str or unicode type for string, not %s: %s" % (
                value.__class__.__name__, str(value))
            )

    def _get_int32(self, int32):
        if isinstance(int32, int) or isinstance(int32, long):
            if int32 < INT32_MIN or int32 > INT32_MAX:
                raise OTSClientError("%s exceeds the range of int32" % int32)
            return int32
        else:
            raise OTSClientError(
                "expect int or long for the value, not %s"
                % int32.__class__.__name__
            )

    def _make_repeated_column_names(self, proto, columns_to_get):
        if columns_to_get is None:
            # if no column name is given, get all primary_key_columns and attribute_columns.
            return

        if not isinstance(columns_to_get, list) and not isinstance(columns_to_get, tuple):
            raise OTSClientError(
                "expect list or tuple for columns_to_get, not %s"
                % columns_to_get.__class__.__name__
            )

        for column_name in columns_to_get:
            proto.append(self._get_unicode(column_name))

    def _make_column_value(self, proto, value):
        # you have to put 'int' under 'bool' in the switch case
        # because a bool is also a int !!!

        if isinstance(value, str) or isinstance(value, unicode):
            string = self._get_unicode(value)
            proto.type = pb2.STRING
            proto.v_string = string
        elif isinstance(value, bool):
            proto.type = pb2.BOOLEAN
            proto.v_bool = value
        elif isinstance(value, int) or isinstance(value, long):
            proto.type = pb2.INTEGER
            proto.v_int = value
        elif isinstance(value, float):
            proto.type = pb2.DOUBLE
            proto.v_double = value
        elif isinstance(value, bytearray):
            proto.type = pb2.BINARY
            proto.v_binary = bytes(value)
        elif value is INF_MIN:
            proto.type = pb2.INF_MIN
        elif value is INF_MAX:
            proto.type = pb2.INF_MAX
        else:
            raise OTSClientError(
                "expect str, unicode, int, long, bool or float for colum value, not %s"
                % value.__class__.__name__
            )

    def _get_column_type(self, type_str):
        enum_map = {
            'INF_MIN'   : pb2.INF_MIN,
            'INF_MAX'   : pb2.INF_MAX,
            'INTEGER'   : pb2.INTEGER,
            'STRING'    : pb2.STRING,
            'BOOLEAN'   : pb2.BOOLEAN,
            'DOUBLE'    : pb2.DOUBLE,
            'BINARY'    : pb2.BINARY,
        }

        if type_str in enum_map:
            return enum_map[type_str]
        else:
            raise OTSClientError(
                "column_type should be one of [%s], not %s" % (
                    ", ".join(enum_map.keys()), str(type_str)
                )
            )

    def _make_condition(self, proto, condition):

        if not isinstance(condition, Condition):
            raise OTSClientError(
                "condition should be an instance of Condition, not %s" %
                condition.__class__.__name__
            )

        enum_map = {
            'IGNORE'            : pb2.IGNORE,
            'EXPECT_EXIST'      : pb2.EXPECT_EXIST,
            'EXPECT_NOT_EXIST'  : pb2.EXPECT_NOT_EXIST
        }

        expectation_str = self._get_unicode(condition.row_existence_expectation) 
        if expectation_str in enum_map:
            proto.row_existence = enum_map[expectation_str]
        else:
            raise OTSClientError(
                "row_existence_expectation should be one of [%s], not %s" % (
                    ", ".join(enum_map.keys()), str(expectation_str)
                )
            )

    def _get_direction(self, direction_str):
        enum_map = {
            'FORWARD'           : pb2.FORWARD,
            'BACKWARD'          : pb2.BACKWARD
        }

        if direction_str in enum_map:
            return enum_map[direction_str]
        else:
            raise OTSClientError(
                "direction should be one of [%s], not %s" % (
                    ", ".join(enum_map.keys()), str(direction_str)
                )
            )

    def _make_column_schema(self, proto, schema_tuple):
        (schema_name, schema_type) = schema_tuple
        proto.name = self._get_unicode(schema_name)
        proto.type = self._get_column_type(schema_type)

    def _make_schemas_with_list(self, proto, schema_list):
        for schema_tuple in schema_list:
            if not isinstance(schema_tuple, tuple):
                raise OTSClientError(
                    "all schemas of primary keys should be tuple, not %s" % (
                        schema_tuple.__class__.__name__
                    )
                )
            schema_proto = proto.add()
            self._make_column_schema(schema_proto, schema_tuple)

    def _make_columns_with_dict(self, proto, column_dict):
        for name, value in column_dict.iteritems():
            item = proto.add()
            item.name = self._get_unicode(name)
            self._make_column_value(item.value, value)

    def _make_update_of_attribute_columns_with_dict(self, proto, column_dict):
    
        if not isinstance(column_dict, dict):
            raise OTSClientError(
                "expect dict for 'update_of_attribute_columns', not %s" % (
                    column_dict.__class__.__name__
                )
            )

        for key, value in column_dict.iteritems():
            if key == 'put':
                if not isinstance(column_dict[key], dict):
                    raise OTSClientError(
                        "expect dict for put operation in 'update_of_attribute_columns', not %s" % (
                            column_dict[key].__class__.__name__
                        )
                    )
                for name, value in column_dict[key].iteritems():
                    item = proto.add()
                    item.type = pb2.PUT
                    item.name = self._get_unicode(name)
                    self._make_column_value(item.value, value)
            elif key == 'delete':
                if not isinstance(column_dict[key], list):
                    raise OTSClientError(
                        "expect list for delete operation in 'update_of_attribute_columns', not %s" % (
                            column_dict[key].__class__.__name__
                        )
                    )
                for name in column_dict[key]:
                    item = proto.add()
                    item.type = pb2.DELETE
                    item.name = self._get_unicode(name)
            else:
                raise OTSClientError(
                    "operation type in 'update_of_attribute_columns' should be 'put' or 'delete', not %s" % (
                        key
                    )
                )

    def _make_table_meta(self, proto, table_meta):
        if not isinstance(table_meta, TableMeta):
            raise OTSClientError(
                "table_meta should be an instance of TableMeta, not %s" 
                % table_meta.__class__.__name__
            )

        proto.table_name = self._get_unicode(table_meta.table_name)
        
        self._make_schemas_with_list(
            proto.primary_key,
            table_meta.schema_of_primary_key,
        )

    def _make_capacity_unit(self, proto, capacity_unit):

        if not isinstance(capacity_unit, CapacityUnit):
            raise OTSClientError(
                "capacity_unit should be an instance of CapacityUnit, not %s" 
                % capacity_unit.__class__.__name__
            )
        
        if capacity_unit.read is None or capacity_unit.write is None:
            raise OTSClientError("both of read and write of CapacityUnit are required")
        proto.read = self._get_int32(capacity_unit.read)
        proto.write = self._get_int32(capacity_unit.write)

    def _make_reserved_throughput(self, proto, reserved_throughput):

        if not isinstance(reserved_throughput, ReservedThroughput):
            raise OTSClientError(
                "reserved_throughput should be an instance of ReservedThroughput, not %s" 
                % reserved_throughput.__class__.__name__
            )
        
        self._make_capacity_unit(proto.capacity_unit, reserved_throughput.capacity_unit)

    def _make_update_capacity_unit(self, proto, capacity_unit):
        if not isinstance(capacity_unit, CapacityUnit):
            raise OTSClientError(
                "capacity_unit should be an instance of CapacityUnit, not %s" 
                % capacity_unit.__class__.__name__
            )
        
        if capacity_unit.read is None and capacity_unit.write is None:
            raise OTSClientError("at least one of read or write of CapacityUnit is required")
        if capacity_unit.read is not None:
            proto.read = self._get_int32(capacity_unit.read)
        if capacity_unit.write is not None:
            proto.write = self._get_int32(capacity_unit.write)

    def _make_update_reserved_throughput(self, proto, reserved_throughput):

        if not isinstance(reserved_throughput, ReservedThroughput):
            raise OTSClientError(
                "reserved_throughput should be an instance of ReservedThroughput, not %s" 
                % reserved_throughput.__class__.__name__
            )
        
        self._make_update_capacity_unit(proto.capacity_unit, reserved_throughput.capacity_unit)

    def _make_batch_get_row(self, proto, batch_list):
        if not isinstance(batch_list, list):
            raise OTSClientError(
                "batch_list should be a list, not %s" 
                % batch_list.__class__.__name__
            )
        
        for (table_name, row_list, columns_to_get) in batch_list:
            table_item = proto.tables.add()
            table_item.table_name = self._get_unicode(table_name)
            self._make_repeated_column_names(table_item.columns_to_get, columns_to_get)
            for primary_key in row_list:
                if isinstance(primary_key, dict):
                    row = table_item.rows.add()
                    self._make_columns_with_dict(row.primary_key, primary_key)
                else:
                    raise OTSClientError(
                        "The row should be a dict, not %s" 
                        % row_item.__class__.__name__
                    )

    def _make_put_row_item(self, proto, put_row_item):
        self._make_condition(proto.condition, put_row_item.condition)
        self._make_columns_with_dict(proto.primary_key, put_row_item.primary_key)
        self._make_columns_with_dict(proto.attribute_columns, put_row_item.attribute_columns)

    def _make_update_row_item(self, proto, update_row_item):
        self._make_condition(proto.condition, update_row_item.condition)
        self._make_columns_with_dict(proto.primary_key, update_row_item.primary_key)
        self._make_update_of_attribute_columns_with_dict(proto.attribute_columns, update_row_item.update_of_attribute_columns)

    def _make_delete_row_item(self, proto, delete_row_item):
        self._make_condition(proto.condition, delete_row_item.condition)
        self._make_columns_with_dict(proto.primary_key, delete_row_item.primary_key)

    def _make_batch_write_row(self, proto, batch_list):
        if not isinstance(batch_list, list):
            raise OTSClientError(
                "batch_list should be a list, not %s" 
                % batch_list.__class__.__name__
            )
        
        enum_map = {
            'put': PutRowItem, 
            'update': UpdateRowItem, 
            'delete': DeleteRowItem
        }
        
        for table_dict in batch_list:
            if not isinstance(table_dict, dict):
                raise OTSClientError(
                    "every item in batch_list should be a dict, not %s" 
                    % table_dict.__class__.__name__
                )

            table_name = table_dict.get('table_name')
            table_item = proto.tables.add()
            table_item.table_name = self._get_unicode(table_name)

            for key,row_list in table_dict.iteritems():
                if key is 'table_name':
                    continue
                if not key in enum_map:
                    raise OTSClientError(
                        "operation type must be one of [%s], not %s" % (
                        ", ".join(enum_map.keys()), str(key))
                    )
                if not isinstance(row_list, list):
                    raise OTSClientError(
                        "rows to write should be a list, not %s" 
                        % row_list.__class__.__name__
                    )
                for row_item in row_list:
                    if not isinstance(row_item, enum_map[key]):
                        raise OTSClientError(
                            "row should be an instance of %s, not %s" % (
                            enum_map[key].__name__, row_item.__class__.__name__)
                        )
                    if key is 'put':
                        row = table_item.put_rows.add()
                        self._make_put_row_item(row, row_item)
                    elif key is 'update':
                        row = table_item.update_rows.add()
                        self._make_update_row_item(row, row_item)
                    elif key is 'delete':
                        row = table_item.delete_rows.add()
                        self._make_delete_row_item(row, row_item)

    def _encode_create_table(self, table_meta, reserved_throughput):
        proto = pb2.CreateTableRequest()
        self._make_table_meta(proto.table_meta, table_meta)
        self._make_reserved_throughput(proto.reserved_throughput, reserved_throughput)
        return proto

    def _encode_delete_table(self, table_name):
        proto = pb2.DeleteTableRequest()
        proto.table_name = self._get_unicode(table_name)
        return proto

    def _encode_list_table(self):
        proto = pb2.ListTableRequest()
        return proto

    def _encode_update_table(self, table_name, reserved_throughput):
        proto = pb2.UpdateTableRequest()
        proto.table_name = self._get_unicode(table_name)
        self._make_update_reserved_throughput(proto.reserved_throughput, reserved_throughput)
        return proto

    def _encode_describe_table(self, table_name):
        proto = pb2.DescribeTableRequest()
        proto.table_name = self._get_unicode(table_name)
        return proto

    def _encode_get_row(self, table_name, primary_key, columns_to_get):
        proto = pb2.GetRowRequest()
        proto.table_name = self._get_unicode(table_name)
        self._make_columns_with_dict(proto.primary_key, primary_key)
        self._make_repeated_column_names(proto.columns_to_get, columns_to_get)
        return proto

    def _encode_put_row(self, table_name, condition, primary_key, attribute_columns):
        proto = pb2.PutRowRequest()
        proto.table_name = self._get_unicode(table_name)
        self._make_condition(proto.condition, condition)
        self._make_columns_with_dict(proto.primary_key, primary_key)
        self._make_columns_with_dict(proto.attribute_columns, attribute_columns)
        return proto

    def _encode_update_row(self, table_name, condition, primary_key, update_of_attribute_columns):
        proto = pb2.UpdateRowRequest()
        proto.table_name = self._get_unicode(table_name)
        self._make_condition(proto.condition, condition)
        self._make_columns_with_dict(proto.primary_key, primary_key)
        self._make_update_of_attribute_columns_with_dict(proto.attribute_columns, update_of_attribute_columns)
        return proto

    def _encode_delete_row(self, table_name, condition, primary_key):
        proto = pb2.DeleteRowRequest()
        proto.table_name = self._get_unicode(table_name)
        self._make_condition(proto.condition, condition)
        self._make_columns_with_dict(proto.primary_key, primary_key)
        return proto

    def _encode_batch_get_row(self, batch_list):
        proto = pb2.BatchGetRowRequest()
        self._make_batch_get_row(proto, batch_list)
        return proto

    def _encode_batch_write_row(self, batch_list):
        proto = pb2.BatchWriteRowRequest()
        self._make_batch_write_row(proto, batch_list)
        return proto

    def _encode_get_range(self, table_name, direction, 
                inclusive_start_primary_key, exclusive_end_primary_key, 
                columns_to_get, limit):
        proto = pb2.GetRangeRequest()
        proto.table_name = self._get_unicode(table_name)
        proto.direction = self._get_direction(direction)
        self._make_columns_with_dict(proto.inclusive_start_primary_key, inclusive_start_primary_key)
        self._make_columns_with_dict(proto.exclusive_end_primary_key, exclusive_end_primary_key)
        self._make_repeated_column_names(proto.columns_to_get, columns_to_get)
        if limit is not None:
            proto.limit = self._get_int32(limit)
        return proto

    def encode_request(self, api_name, *args, **kwargs):
        if api_name not in self.api_encode_map:
            raise OTSClientError("No PB encode method for API %s" % api_name)

        handler = self.api_encode_map[api_name]
        return handler(*args, **kwargs)
