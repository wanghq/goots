# -*- coding: utf8 -*-

import google.protobuf.text_format as text_format

from ots2.metadata import *
import ots2.protobuf.ots_protocol_2_pb2 as pb2

class OTSProtoBufferDecoder:

    def __init__(self, encoding):
        self.encoding = encoding

        self.api_decode_map = {
            'CreateTable'       : self._decode_create_table,
            'ListTable'         : self._decode_list_table,
            'DeleteTable'       : self._decode_delete_table,
            'DescribeTable'     : self._decode_describe_table,
            'UpdateTable'       : self._decode_update_table,
            'GetRow'            : self._decode_get_row,
            'PutRow'            : self._decode_put_row,
            'UpdateRow'         : self._decode_update_row,
            'DeleteRow'         : self._decode_delete_row,
            'BatchGetRow'       : self._decode_batch_get_row,
            'BatchWriteRow'     : self._decode_batch_write_row,
            'GetRange'          : self._decode_get_range,
        }

    def _parse_string(self, string):
        if string == '':
            return None
        else:
            return string

    def _parse_column_type(self, column_type_enum):
        reverse_enum_map = {
            pb2.INF_MIN : 'INF_MIN',
            pb2.INF_MAX : 'INF_MAX',
            pb2.INTEGER : 'INTEGER',
            pb2.STRING  : 'STRING',
            pb2.BOOLEAN : 'BOOLEAN',
            pb2.DOUBLE  : 'DOUBLE',
            pb2.BINARY  : 'BINARY'
        }
        if column_type_enum in reverse_enum_map:
            return reverse_enum_map[column_type_enum]
        else:
            raise OTSClientError("invalid value for column type: %s" % str(column_type_enum))

    def _parse_value(self, proto):
        if proto.type == pb2.INTEGER:
            return proto.v_int
        elif proto.type == pb2.STRING:
            return proto.v_string
        elif proto.type == pb2.BOOLEAN:
            return proto.v_bool
        elif proto.type == pb2.DOUBLE:
            return proto.v_double
        elif proto.type == pb2.BINARY:
            return bytearray(proto.v_binary)
        else:
            raise OTSClientError("invalid column value type: %s" % str(proto.type))

    def _parse_schema_list(self, proto):
        ret = []
        for item in proto:
            ret.append((item.name, self._parse_column_type(item.type)))
        return ret

    def _parse_column_dict(self, proto):
        ret = {}
        for item in proto:
            ret[item.name] = self._parse_value(item.value)
        return ret

    def _parse_row(self, proto):
        return (
            self._parse_column_dict(proto.primary_key_columns),
            self._parse_column_dict(proto.attribute_columns)
        )

    def _parse_row_list(self, proto):
        row_list = [] 
        for row_item in proto:
            row_list.append(self._parse_row(row_item))
        return row_list

    def _parse_capacity_unit(self, proto):
        if proto is None:
            capacity_unit = None
        else:
            cu_read = proto.read if proto.HasField('read') else 0
            cu_write = proto.write if proto.HasField('write') else 0
            capacity_unit = CapacityUnit(cu_read, cu_write)
        return capacity_unit

    def _parse_reserved_throughput_details(self, proto):
        last_decrease_time = proto.last_decrease_time if proto.HasField('last_decrease_time') else None
        capacity_unit = self._parse_capacity_unit(proto.capacity_unit)

        reserved_throughput_details = ReservedThroughputDetails(
            capacity_unit,
            proto.last_increase_time, 
            last_decrease_time,
            proto.number_of_decreases_today
        )
        return reserved_throughput_details

    def _parse_get_row_item(self, proto):
        row_list = []
        for row_item in proto:
            if row_item.is_ok:
                error_code = None
                error_message = None
                capacity_unit = self._parse_capacity_unit(row_item.consumed.capacity_unit)
                primary_key_columns = self._parse_column_dict(row_item.row.primary_key_columns)
                attribute_columns = self._parse_column_dict(row_item.row.attribute_columns)
            else:
                error_code = row_item.error.code
                error_message = row_item.error.message if row_item.error.HasField('message') else ''
                if row_item.HasField('consumed'):
                    capacity_unit = self._parse_capacity_unit(row_item.consumed.capacity_unit)
                else:
                    capacity_unit = None
                primary_key_columns = None
                attribute_columns = None

            row_data_item = RowDataItem(
                row_item.is_ok, error_code, error_message,
                capacity_unit, primary_key_columns, attribute_columns
            )
            row_list.append(row_data_item)
        
        return row_list

    def _parse_batch_get_row(self, proto):
        rows = []
        for table_item in proto:
            rows.append(self._parse_get_row_item(table_item.rows)) 
        return rows

    def _parse_write_row_item(self, proto):
        row_list = []
        for row_item in proto:
            if row_item.is_ok:
                error_code = None
                error_message = None
                consumed = self._parse_capacity_unit(row_item.consumed.capacity_unit)
            else:
                error_code = row_item.error.code
                error_message = row_item.error.message if row_item.error.HasField('message') else ''
                if row_item.HasField('consumed'):
                    consumed = self._parse_capacity_unit(row_item.consumed.capacity_unit)
                else:
                    consumed = None

            write_row_item = BatchWriteRowResponseItem(
                row_item.is_ok, error_code, error_message, consumed
            )
            row_list.append(write_row_item)
        
        return row_list

    def _parse_batch_write_row(self, proto):
        result_list = []
        for table_item in proto:
            table_dict = {}
            if table_item.put_rows:
                put_list = self._parse_write_row_item(table_item.put_rows)
                table_dict['put'] = put_list
            if table_item.update_rows:
                update_list = self._parse_write_row_item(table_item.update_rows)
                table_dict['update'] = update_list
            if table_item.delete_rows:
                delete_list = self._parse_write_row_item(table_item.delete_rows)
                table_dict['delete'] = delete_list
            result_list.append(table_dict)
        return result_list

    def _decode_create_table(self, body):
        proto = pb2.CreateTableResponse()
        proto.ParseFromString(body)
        return None, proto

    def _decode_list_table(self, body):
        proto = pb2.ListTableResponse()
        proto.ParseFromString(body)
        names = tuple(proto.table_names)
        return names, proto

    def _decode_delete_table(self, body):
        proto = pb2.DeleteTableResponse()
        proto.ParseFromString(body)
        return None, proto
        
    def _decode_describe_table(self, body):
        proto = pb2.DescribeTableResponse()
        proto.ParseFromString(body)

        table_meta = TableMeta(
            proto.table_meta.table_name,
            self._parse_schema_list(
                proto.table_meta.primary_key
            )
        )
        
        reserved_throughput_details = self._parse_reserved_throughput_details(proto.reserved_throughput_details)
        describe_table_response = DescribeTableResponse(table_meta, reserved_throughput_details)
        return describe_table_response, proto

    def _decode_update_table(self, body):
        proto = pb2.UpdateTableResponse()
        proto.ParseFromString(body)

        reserved_throughput_details = self._parse_reserved_throughput_details(proto.reserved_throughput_details)
        update_table_response = UpdateTableResponse(reserved_throughput_details)

        return update_table_response, proto

    def _decode_get_row(self, body):
        proto = pb2.GetRowResponse()
        proto.ParseFromString(body)

        primary_key_columns, attribute_columns = self._parse_row(proto.row)
        consumed = self._parse_capacity_unit(proto.consumed.capacity_unit)
        return (consumed, primary_key_columns, attribute_columns), proto

    def _decode_put_row(self, body):
        proto = pb2.PutRowResponse()
        proto.ParseFromString(body)

        consumed = self._parse_capacity_unit(proto.consumed.capacity_unit)
        return consumed, proto

    def _decode_update_row(self, body):
        proto = pb2.UpdateRowResponse()
        proto.ParseFromString(body)

        consumed = self._parse_capacity_unit(proto.consumed.capacity_unit)
        return consumed, proto

    def _decode_delete_row(self, body):
        proto = pb2.DeleteRowResponse()
        proto.ParseFromString(body)

        consumed = self._parse_capacity_unit(proto.consumed.capacity_unit)
        return consumed, proto

    def _decode_batch_get_row(self, body):
        proto = pb2.BatchGetRowResponse()
        proto.ParseFromString(body)

        rows = self._parse_batch_get_row(proto.tables)
        return rows, proto

    def _decode_batch_write_row(self, body):
        proto = pb2.BatchWriteRowResponse()
        proto.ParseFromString(body)

        rows = self._parse_batch_write_row(proto.tables)
        return rows, proto

    def _decode_get_range(self, body):
        proto = pb2.GetRangeResponse()
        proto.ParseFromString(body)
        
        capacity_unit = self._parse_capacity_unit(proto.consumed.capacity_unit)
        next_start_pk = self._parse_column_dict(proto.next_start_primary_key)
        if not next_start_pk:
            next_start_pk = None
        row_list = self._parse_row_list(proto.rows)
        return (capacity_unit, next_start_pk, row_list), proto

    def decode_response(self, api_name, response_body):
        if api_name not in self.api_decode_map:
            raise OTSClientError("No PB decode method for API %s" % api_name)

        handler = self.api_decode_map[api_name]
        return handler(response_body)

