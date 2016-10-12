# -*- coding: utf8 -*-

from ots2.error import *

__all__ = [
    'INF_MIN',
    'INF_MAX',
    'TableMeta',
    'CapacityUnit',
    'ReservedThroughput',
    'ReservedThroughputDetails',
    'UpdateTableResponse',
    'DescribeTableResponse',
    'RowDataItem',
    'Condition',
    'PutRowItem',
    'UpdateRowItem',
    'DeleteRowItem',
    'BatchWriteRowResponseItem'
]


class TableMeta(object):

    def __init__(self, table_name, schema_of_primary_key):
        # schema_of_primary_key: [('PK0', 'STRING'), ('PK1', 'INTEGER'), ...]
        self.table_name = table_name
        self.schema_of_primary_key = schema_of_primary_key


class CapacityUnit(object):

    def __init__(self, read=0, write=0):
        self.read = read
        self.write = write


class ReservedThroughput(object):

    def __init__(self, capacity_unit):
        self.capacity_unit = capacity_unit


class ReservedThroughputDetails(object):
    
    def __init__(self, capacity_unit, last_increase_time, last_decrease_time, number_of_decreases_today):
        self.capacity_unit = capacity_unit
        self.last_increase_time = last_increase_time
        self.last_decrease_time = last_decrease_time
        self.number_of_decreases_today = number_of_decreases_today


class UpdateTableResponse(object):

    def __init__(self, reserved_throughput_details):
        self.reserved_throughput_details = reserved_throughput_details


class DescribeTableResponse(object):

    def __init__(self, table_meta, reserved_throughput_details):
        self.table_meta = table_meta
        self.reserved_throughput_details = reserved_throughput_details


class RowDataItem(object):

    def __init__(self, is_ok, error_code, error_message, consumed, primary_key_columns, attribute_columns):
        # is_ok can be True or False
        # when is_ok is False,
        #     error_code & error_message are available
        # when is_ok is True,
        #     consumed & primary_key_columns & attribute_columns are available
        self.is_ok = is_ok
        self.error_code = error_code
        self.error_message = error_message
        self.consumed = consumed
        self.primary_key_columns = primary_key_columns
        self.attribute_columns = attribute_columns


class Condition(object):

    def __init__(self, row_existence_expectation):
        self.row_existence_expectation = row_existence_expectation 


class PutRowItem(object):

    def __init__(self, condition, primary_key, attribute_columns):
        self.condition = condition
        self.primary_key = primary_key
        self.attribute_columns = attribute_columns


class UpdateRowItem(object):
    
    def __init__(self, condition, primary_key, update_of_attribute_columns):
        self.condition = condition
        self.primary_key = primary_key
        self.update_of_attribute_columns = update_of_attribute_columns


class DeleteRowItem(object):
    
    def __init__(self, condition, primary_key):
        self.condition = condition
        self.primary_key = primary_key


class BatchWriteRowResponseItem(object):

    def __init__(self, is_ok, error_code, error_message, consumed):
        self.is_ok = is_ok
        self.error_code = error_code
        self.error_message = error_message
        self.consumed = consumed


class INF_MIN(object):
    # for get_range
    pass


class INF_MAX(object):
    # for get_range
    pass

