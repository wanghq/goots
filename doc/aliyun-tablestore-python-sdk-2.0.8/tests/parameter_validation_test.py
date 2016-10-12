# -*- coding: utf8 -*-

import unittest
from lib.ots2_api_test_base import *
from ots2 import *
import lib.restriction as restriction
import lib.test_config as test_config
import copy
import time

class ParameterValidationTest(OTS2APITestBase):

    """参数合法性检查测试（不包含限制项）"""

    def _get_client(self, instance_name):
        client = OTSClient(
            test_config.OTS_ENDPOINT,
            test_config.OTS_ID,
            test_config.OTS_SECRET,
            instance_name
            )
        return client

    def _invalid_instance_or_table_name_op(self, client, table = None, error_code="", error_message=""):
        table_name = table if table != None else "table_test"
        table_meta = TableMeta(table_name, [("PK", "STRING")])
        reserved_throughput = ReservedThroughput(CapacityUnit(restriction.MaxReadWriteCapacityUnit, restriction.MaxReadWriteCapacityUnit))
        #create_table
        try:
            client.create_table(table_meta, reserved_throughput)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #delete_table
        try:
            client.delete_table(table_name)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)

        #list_table
        if table == None:
            try:
                client.list_table()
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)

        #update_table
        try:
            client.update_table(table_name, ReservedThroughput(CapacityUnit(1, 1)))
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #describe_table
        try:
            client.describe_table(table_name)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #get_row
        try:
            client.get_row(table_name, {"PK": "x"})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #put_row
        try:
            client.put_row(table_name, Condition("IGNORE"), {"PK": "x"}, {"COL": "x"})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #update_row
        try:
            client.update_row(table_name, Condition("IGNORE"), {"PK": "x"}, {'put':{"COL": "x"}})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #delete_row
        try:
            client.delete_row(table_name, Condition("IGNORE"), {"PK": "x"})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        
        #bacth_get_row
        batches = [(table_name, [{"PK": "x"}], [])]
        try:
            client.batch_get_row(batches)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)

        #batch_write_row
        put_row_item = {"table_name": table_name, "put": [PutRowItem(Condition("IGNORE"), {"PK": "x"}, {"COL": "x"})]}
        update_row_item = {"table_name": table_name, "update": [UpdateRowItem(Condition("IGNORE"), {"PK": "x"}, {'put':{"COL": "x"}})]}
        delete_row_item = {"table_name": table_name, "delete": [DeleteRowItem(Condition("IGNORE"), {"PK": "x"})]}
        batch_write_items = [put_row_item, update_row_item, delete_row_item]
        for item in batch_write_items:
            batches = [item]
            try:
                client.batch_write_row(batches)
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)
        #get_range
        try:
            client.get_range(table_name, "FORWARD", {"PK": INF_MIN},{"PK": INF_MAX})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        
    def test_update_row_columns_empty(self):
        """对于UpdateRow和BatchWriteRow的UpdateRow操作，column的个数为0，期望返回ErrorCode: OTSParameterInvalid """
        #create test table
        table_name = "table_test"
        table_meta = TableMeta(table_name, [('PK0', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(restriction.MaxReadWriteCapacityUnit,restriction.MaxReadWriteCapacityUnit))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        primary_keys = {"PK0": "x"}
        
        #update_row
        try:
            self.client_test.update_row(table_name, Condition("IGNORE"), primary_keys, {})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "No column specified while updating row.")
        #batch_write_row
        update_row_item = {"table_name": table_name, "update": [UpdateRowItem(Condition("IGNORE"), primary_keys, {})]}
        write_batches = [update_row_item]
        try:
            self.client_test.batch_write_row(write_batches)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "No attribute column specified to update row #%d in table: '%s'." % (0, table_name))

    def test_invalid_instance_name(self):
        """对于每个API，测试instance name为空，'*', '#', '中文', 'aa'的情形，期望返回ErrorCode: OTSParameterInvalid"""
        for instance_name in ["", "*", "#", "中文", "aa", "_aaa", "a###", "h***", "-aa", "aa-"]:
            client = self._get_client(instance_name)
            self._invalid_instance_or_table_name_op(client, error_code="OTSParameterInvalid", error_message="Invalid instance name: '%s'." % instance_name)
            
    def _valid_instance_or_table_name_op(self, client, table = None):
        table_name = table if table != None else "table_test"
        table_meta = TableMeta(table_name, [("PK", "STRING")])
        reserved_throughput = ReservedThroughput(CapacityUnit(1, 2))
        #create_table

        expect_increase_time = int(time.time())
        client.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        #describe_table
        response = client.describe_table(table_name)
        self.assert_DescribeTableResponse(response, CapacityUnit(1, 2), table_meta)
        #get_row
        consumed, primary_keys, columns = client.get_row(table_name, {"PK": "x"})
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(primary_keys, {})
        self.assert_equal(columns, {})
        #batch_get_row
        batches = [(table_name, [{"PK": "x"}], [])]
        response = client.batch_get_row(batches)
        expect_row_data_item = RowDataItem(True, "", "", CapacityUnit(1, 0), {}, {})
        expect_response = [[expect_row_data_item]]
        self.assert_RowDataItem_equal(response, expect_response)
        #get_range
        consumed, next_start_primary_keys, rows = client.get_range(table_name, 'FORWARD', {"PK": INF_MIN}, {"PK": INF_MAX})
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(next_start_primary_keys, None)
        self.assert_equal(rows, [])
        #put_row
        consumed = client.put_row(table_name, Condition("IGNORE"), {"PK": "x"}, {'COL': 'x'})
        self.assert_consumed(consumed, CapacityUnit(0, 1))
        #update_row
        consumed = client.update_row(table_name, Condition("IGNORE"), {"PK": "x"}, {'put':{"COL": "x1"}})
        self.assert_consumed(consumed, CapacityUnit(0, 1))
        #delete_row
        consumed = client.delete_row(table_name, Condition("IGNORE"), {"PK": "x"})
        self.assert_consumed(consumed, CapacityUnit(0, 1))
        #batch_write_row
        put_row_item = {"table_name": table_name, "put": [PutRowItem(Condition("IGNORE"), {"PK": "x"}, {'COL': 'x'})]}
        update_row_item = {"table_name": table_name, "update": [UpdateRowItem(Condition("IGNORE"), {"PK": "x"}, {'put':{'COL': 'x1'}})]}
        delete_row_item = {"table_name": table_name, "delete": [DeleteRowItem(Condition("IGNORE"), {"PK": "x"})]}
        op_list = [('put', put_row_item), ('update', update_row_item), ('delete', delete_row_item)]
        for op_type, item in op_list:
            write_batches = [item]
            response = client.batch_write_row(write_batches)
            expect_write_data_item = BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, 1))
            expect_response = [{op_type : [expect_write_data_item]}]
            self.assert_BatchWriteRowResponseItem(response, expect_response)
        #delete_table
        client.delete_table(table_name)
        #list_table
        if table == None:
            table_list = client.list_table()
            self.assert_equal(table_list, ())

    def test_instance_name_of_huge_size(self):
        """用一个长度为2K的instance name去访问OTS，期望返回OTSServiceError"""
        instance_name = 'X'  * (2 * 1024)
        client = self._get_client(instance_name)
        try:
            client.list_table()
            self.assert_false()
        except OTSServiceError:
            # i don't care what kind of error
            pass

    def test_invalid_table_name(self):
        """测试所有相关API中表名为空，'0', '#', '中文', 'T#', 'T中文', '3t', '-'的情况，期望返回ErrorCode: OTSParameterInvalid"""
        table_list = ['', '0', '#', '中文', 'T#', 'T中文', '3t', '-']
        for table_name in table_list:
            time.sleep(2)  # to avoid too frequently table operation
            self._invalid_instance_or_table_name_op(self.client_test, table_name, error_code="OTSParameterInvalid", error_message="Invalid table name: '%s'." % table_name)

    def test_valid_table_name(self):
        """测试所有相关API中表名为'_0', '_T', 'A0'的情况，期望操作成功"""
        table_list = ['_0', '_T', 'A0']
        for table_name in table_list:
            self._valid_instance_or_table_name_op(self.client_test, table_name)
    
    def _invalid_column_name_test(self,table_name, primary_keys, columns_name, error_code="", error_message=""):

        #get_row
        try:
            self.client_test.get_row(table_name, primary_keys, [columns_name])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)

        #put_row
        try:
            self.client_test.put_row(table_name, Condition("IGNORE"), primary_keys, {columns_name: "x"})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #update_row
        try:
            self.client_test.update_row(table_name, Condition("IGNORE"), primary_keys, {'put':{columns_name: "x1"}})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #batch_get_row
        try:
            get_batches = [(table_name, [primary_keys], [columns_name])]
            self.client_test.batch_get_row(get_batches)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)

        put_row_item = {"table_name": table_name, "put":[PutRowItem(Condition("IGNORE"), primary_keys, {columns_name: "x"})]}
        update_row_item = {"table_name": table_name, "update": [UpdateRowItem(Condition("IGNORE"), primary_keys, {'put':{columns_name: "x"}})]}
        batches_list = [put_row_item, update_row_item]
        for item in batches_list:
            write_batches = [item]
            try:
                self.client_test.batch_write_row(write_batches)
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)
        #get_range
        try:
            self.client_test.get_range(table_name, 'FORWARD', {"PK0": INF_MIN}, {"PK0": INF_MAX}, [columns_name], None)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
            
    def test_invalid_column_name(self):
        """测试所有相关API中column name为空，'0', '#', '中文', 'T#', 'T中文'的情况，期望返回ErrorCode: OTSParameterInvalid"""
        #create test table
        table_name = "table_test"
        table_meta = TableMeta(table_name, [('PK0', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(1, 1))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test')
        primary_keys = {"PK0": "x"}
        column_name_list = ['', '"', '"abc', '!', '!abc', '#abc', '#', '中文', 'T中文']
        for column_name in column_name_list:
            self._invalid_column_name_test(table_name, primary_keys, column_name, error_code="OTSParameterInvalid", error_message="Invalid column name: '%s'." % column_name)

    def _valid_column_name_test(self, table_name, columns_name):
        pk = {"PK": columns_name}
        #get_row
        consumed, primary_keys, columns = self.client_test.get_row(table_name, pk, [columns_name])
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(primary_keys, {})
        self.assert_equal(columns, {})
        #batch_get_row
        batches = [(table_name, [pk], [columns_name])]
        response = self.client_test.batch_get_row(batches)
        expect_row_data_item = RowDataItem(True, "", "", CapacityUnit(1, 0), {}, {})
        expect_response = [[expect_row_data_item]]
        self.assert_RowDataItem_equal(response, expect_response)
        #get_range
        consumed, next_start_primary_keys, rows = self.client_test.get_range(table_name, 'FORWARD', {"PK": INF_MIN}, {"PK": INF_MAX}, [columns_name], None)
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(next_start_primary_keys, None)
        self.assert_equal(rows, [])
        #put_row
        consumed = self.client_test.put_row(table_name, Condition("IGNORE"), pk, {columns_name: 'x'})
        self.assert_consumed(consumed, CapacityUnit(0, self.sum_CU_from_row(pk, {columns_name: 'x'})))
        #update_row
        consumed = self.client_test.update_row(table_name, Condition("IGNORE"), pk, {'put':{columns_name: "x1"}})
        self.assert_consumed(consumed, CapacityUnit(0, self.sum_CU_from_row(pk, {columns_name: "x1"})))

        #batch_write_row
        put_row_item = {"table_name": table_name, "put": [PutRowItem(Condition("IGNORE"), pk, {columns_name: 'x'})]}
        update_row_item = {"table_name": table_name, "update": [UpdateRowItem(Condition("IGNORE"), pk, {'put':{columns_name: 'x1'}})]}
        write_batches = [put_row_item]
        response = self.client_test.batch_write_row(write_batches)
        expect_write_data_item = {"put": [BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, self.sum_CU_from_row(pk, {columns_name: "x"})))]}
        expect_response = [expect_write_data_item]
        self.assert_BatchWriteRowResponseItem(response, expect_response) 
        write_batches = [update_row_item]
        response = self.client_test.batch_write_row(write_batches)
        expect_write_data_item = {"update": [BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, self.sum_CU_from_row(pk, {columns_name: "x1"})))]}
        expect_response = [expect_write_data_item]
        self.assert_BatchWriteRowResponseItem(response, expect_response) 

    def test_valid_column_name(self):
        """测试所有相关API中column name为'_0', '_T', 'A0'的情况，期望操作成功"""
        table_name = "table_test"
        table_meta = TableMeta(table_name, [("PK", "STRING")])
        reserved_throughput = ReservedThroughput(CapacityUnit(restriction.MaxReadWriteCapacityUnit, restriction.MaxReadWriteCapacityUnit))
        #create_table
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)
        column_name_list = ['_0', '_T', 'A0', '_a-b', 'a-b-c', 'waef-', '_--', '_%', '_a%b', 'abc%def', '/', '\\', '|', '&', '$', '~']
        for column_name in column_name_list:
            self._valid_column_name_test(table_name, column_name)

    def _invalid_primary_key_type_create_table(self, _type, error_code='', error_message=''):
        table_name = "table_test"
        PK_schema = []
        for i in range(restriction.MaxPKColumnNum):
            PK_schema.append(("PK%d" % i, "STRING"))

        for i in range(restriction.MaxPKColumnNum):
            PK = copy.copy(PK_schema)
            PK[i] = ("PK%d" % i, _type)
            table_meta = TableMeta(table_name, PK)
            try:
                self.client_test.create_table(table_meta, ReservedThroughput(CapacityUnit(10, 10)))
                self.assert_false()
            except Exception as e:
                self.assert_error(e, 400, error_code, error_message)

    def _invalid_pk_type_test(self, invalid_pk, is_range=None, error_code="", error_message="", error_message_for_range=""):
        table_name = "table_test"

        pk_schema, valid_primary_keys = self.get_primary_keys(restriction.MaxPKColumnNum, "STRING")
        
        for pk in valid_primary_keys.keys():
            primary_keys = copy.copy(valid_primary_keys)
            primary_keys[pk] = invalid_pk
            #get row
            try:
                self.client_test.get_row(table_name, primary_keys, [])
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)
            #put_row
            try:
                self.client_test.put_row(table_name, Condition("IGNORE"), primary_keys, {'COL': 'x'})
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)
            #update_row
            try:
                self.client_test.update_row(table_name, Condition("IGNORE"), primary_keys, {'put':{'COL': 'x'}})
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)
            #delete_row
            try:
                self.client_test.delete_row(table_name, Condition("IGNORE"), primary_keys)
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)
            #batch_get_row
            get_batches = [(table_name, [primary_keys], [])]
            try:
                self.client_test.batch_get_row(get_batches)
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)
            #batch_write_row
            put_row_item = {"table_name": table_name, "put":[PutRowItem(Condition("IGNORE"), primary_keys, {'COL': 'x'})]}
            update_row_item = {"table_name": table_name, "update": [UpdateRowItem(Condition("IGNORE"), primary_keys, {'put':{'COL': 'x'}})]}
            delete_row_item = {"table_name": table_name, "delete": [DeleteRowItem(Condition("IGNORE"), primary_keys)]}
            batches_list = [put_row_item, update_row_item, delete_row_item]
            for item in batches_list:
                write_batches = [item]
                try:
                    self.client_test.batch_write_row(write_batches)
                    self.assert_false()
                except OTSServiceError as e:
                    self.assert_error(e, 400, error_code, error_message)
        #get_range
        def assert_get_range(table_name, inclusive_start_primary_keys, exclusive_end_primary_keys):
            try:
                self.client_test.get_range(table_name, 'FORWARD', inclusive_start_primary_keys, exclusive_end_primary_keys, [], None)
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message_for_range)
        
        if is_range == None:
            pk_schema, exclusive_primary_keys = self.get_primary_keys(restriction.MaxPKColumnNum, "STRING", "PK", INF_MAX)
            pk_schema, inclusive_primary_keys = self.get_primary_keys(restriction.MaxPKColumnNum, "STRING", "PK", INF_MIN)
            for pk in valid_primary_keys.keys():
                primary_keys = copy.copy(valid_primary_keys)
                primary_keys[pk] = invalid_pk
                assert_get_range(table_name, inclusive_primary_keys, primary_keys)
                assert_get_range(table_name, primary_keys, exclusive_primary_keys)
        
    def test_PK_double(self):
        """测试所有相关API中PK为double的情况，期望返回OTSParameterInvalid"""
        double_value = 5.6
        self._invalid_primary_key_type_create_table("DOUBLE", error_code="OTSParameterInvalid", error_message="DOUBLE is an invalid type for the primary key.")
        self._invalid_pk_type_test(double_value, error_code="OTSParameterInvalid", error_message="DOUBLE is an invalid type for the primary key.", error_message_for_range="DOUBLE is an invalid type for the primary key in GetRange.")

    def test_PK_boolean(self):
        """测试所有相关API中 PK为boolean的情况，期望返回OTSParameterInvalid"""
        boolean_value = True
        self._invalid_primary_key_type_create_table("BOOLEAN", error_code="OTSParameterInvalid", error_message="BOOLEAN is an invalid type for the primary key.")
        self._invalid_pk_type_test(boolean_value, error_code="OTSParameterInvalid", error_message="BOOLEAN is an invalid type for the primary key.", error_message_for_range="BOOLEAN is an invalid type for the primary key in GetRange.")
        
    def test_table_not_exist(self):
        """测试除CreateTable以外所有API，操作的表不存在的情况，期望返回OTSObjectNotExist"""
        table_name = "table_name_for_table_not_exist_test"
        #delete_table
        try:
            self.client_test.delete_table(table_name)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")

        #update_table
        try:
            self.client_test.update_table(table_name, ReservedThroughput(CapacityUnit(10, 10)))
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")
        #describe_table
        try:
            self.client_test.describe_table(table_name)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")

        #get_row
        try:
            self.client_test.get_row(table_name, {"PK": "x"})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")
        #put_row
        try:
            self.client_test.put_row(table_name, Condition("IGNORE"), {"PK": "x"}, {"COL": "x"})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")
        #update_row
        try:
            self.client_test.update_row(table_name, Condition("IGNORE"), {"PK": "x"}, {'put':{"COL": "x"}})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")
        #delete_row
        try:
            self.client_test.delete_row(table_name, Condition("IGNORE"), {"PK": "x"})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")
        #bacth_get_row
        batches = [(table_name, [{"PK": "x"}], [])]
        response = self.client_test.batch_get_row(batches)
        expect_row_data_item = RowDataItem(False, "OTSObjectNotExist", "Requested table does not exist.", None, None, None) 
        expect_response = [[expect_row_data_item]]
        self.assert_RowDataItem_equal(response, expect_response)
       
        #batch_write_row
        put_row_item = {"table_name": table_name, "put": [PutRowItem(Condition("IGNORE"), {"PK": "x"}, {"COL": "x"})]}
        update_row_item = {"table_name": table_name, "update": [UpdateRowItem(Condition("IGNORE"), {"PK": "x"}, {'put':{"COL": "x"}})]}
        delete_row_item = {"table_name": table_name, "delete": [DeleteRowItem(Condition("IGNORE"), {"PK": "x"})]}
        batch_write_items = [put_row_item, update_row_item, delete_row_item]
        resp_key = ["put", "update", "delete"]
        for item, key in zip(batch_write_items, resp_key):
            batches = [item]
            response = self.client_test.batch_write_row(batches)
            expect_write_data_item = BatchWriteRowResponseItem(False, "OTSObjectNotExist", "Requested table does not exist.", None) 
            expect_response = [{key : [expect_write_data_item]}]
            self.assert_BatchWriteRowResponseItem(response, expect_response)

        #get_range
        try:
            self.client_test.get_range(table_name, "FORWARD", {"PK": INF_MIN},{"PK": INF_MAX})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")

    def _column_restriction_test_except(self, columns, error_code="", error_message=""):

        #create test table
        table_name = "table_test"
        table_meta = TableMeta(table_name, [('PK0', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(restriction.MaxReadWriteCapacityUnit,restriction.MaxReadWriteCapacityUnit))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test')

        primary_keys = {"PK0": "x"}

        #put_row
        try:
            self.client_test.put_row(table_name, Condition("IGNORE"), primary_keys, columns)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)
        #update_row
        try:
            self.client_test.update_row(table_name, Condition("IGNORE"), primary_keys, {'put':columns})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_message)

        put_row_item = {"table_name": table_name, "put": [PutRowItem(Condition("IGNORE"), primary_keys, columns)]}
        update_row_item = {"table_name": table_name, "update": [UpdateRowItem(Condition("IGNORE"), primary_keys, {'put':columns})]}
        batches_list = [put_row_item, update_row_item]
        for item in batches_list:
            write_batches = [item]
            try:
                self.client_test.batch_write_row(write_batches)
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, error_code, error_message)

 
    def test_column_type_INF_MIN(self):
        """测试所有API中，主键和列的类型（get_range的start和end除外）为INF_MIN的情况，期望返回OTSParameterInvalid"""
        self._invalid_pk_type_test(INF_MIN, True, error_code="OTSParameterInvalid", error_message="INF_MIN is an invalid type for the primary key.")
        self._column_restriction_test_except({"COL": INF_MIN}, "OTSParameterInvalid", "INF_MIN is an invalid type for the attribute column.")

    def test_column_type_INF_MAX(self):
        """测试所有API中，主键和列的类型（get_range的start和end除外）为INF_MAX的情况，期望返回OTSParameterInvalid"""
        self._invalid_pk_type_test(INF_MAX, True, error_code="OTSParameterInvalid", error_message="INF_MAX is an invalid type for the primary key.")
        self._column_restriction_test_except({"COL": INF_MAX}, "OTSParameterInvalid", "INF_MAX is an invalid type for the attribute column.")


if __name__ == '__main__':
    unittest.main()
