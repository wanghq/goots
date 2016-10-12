# -*- coding: utf8 -*-

import unittest
from lib.ots2_api_test_base import OTS2APITestBase
import lib.restriction as restriction
from ots2 import *
from ots2.error import *
import time
import logging

class TableOperationTest(OTS2APITestBase):

    """表级别操作测试"""

    def test_delete_existing_table(self):
        """删除一个存在的表，期望成功, list_table()确认表已经删除, describe_table()返回异常OTSObjectNotExist"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_test_delete_existing'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.client_test.delete_table(table_name)
        self.assert_equal(False, table_name in self.client_test.list_table())

        try:
            self.client_test.describe_table(table_name)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 404, "OTSObjectNotExist", "Requested table does not exist.")

    def test_create_table_already_exist(self):
        """创建一个表，表名与现有表重复，期望返回ErrorCode: OTSObjectAlreadyExist, list_table()确认没有2个重名的表"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_test_already_exist'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        self.client_test.create_table(table_meta, reserved_throughput)

        table_meta_new = TableMeta(table_name, [('PK2', 'STRING'), ('PK3', 'STRING')])
        try:
            self.client_test.create_table(table_meta_new, reserved_throughput)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 409, "OTSObjectAlreadyExist", "Requested table already exists.")

        table_list = self.client_test.list_table()
        self.assert_equal(table_list, (table_name,))


    def test_duplicate_PK_name_in_table_meta(self):
        """创建表的时候，TableMeta中有2个PK列，列名重复，期望返回OTSParameterInvalid，list_table()确认没有这个表"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_test_duplicate_PK'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK0', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        try:
            self.client_test.create_table(table_meta, reserved_throughput)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "The name of primary key must be unique.")

        self.assert_equal(False, table_name in self.client_test.list_table())


    def test_PK_type_STRING(self):
        """创建表的时候，TableMeta中有4个PK列，都为STRING类型，期望正常，describe_table()获取信息与创表参数一致"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_PK_type_STRING'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'STRING'), ('PK2', 'STRING'), ('PK3', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        self.client_test.create_table(table_meta, reserved_throughput)

        describe_response = self.client_test.describe_table(table_name)
        self.assert_DescribeTableResponse(describe_response, reserved_throughput.capacity_unit, table_meta)

    def test_PK_type_INTEGER(self):
        """创建表的时候，TableMeta中有4个PK列，都为INTEGER类型，期望正常，describe_table()获取信息与创表参数一致"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_PK_type_INTEGER'
        table_meta = TableMeta(table_name, [('PK0', 'INTEGER'), ('PK1', 'INTEGER'), ('PK2', 'INTEGER'), ('PK3', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        self.client_test.create_table(table_meta, reserved_throughput)

        describe_response = self.client_test.describe_table(table_name)
        self.assert_DescribeTableResponse(describe_response, reserved_throughput.capacity_unit, table_meta)

    def test_PK_type_BINARY(self):
        """创建表的时候，TableMeta中有4个PK列，都为BINARY类型，期望正常，describe_table()获取信息与创表参数一致"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_PK_type_BINARY'
        table_meta = TableMeta(table_name, [('PK0', 'BINARY'), ('PK1', 'BINARY'), ('PK2', 'BINARY'), ('PK3', 'BINARY')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        self.client_test.create_table(table_meta, reserved_throughput)

        describe_response = self.client_test.describe_table(table_name)
        self.assert_DescribeTableResponse(describe_response, reserved_throughput.capacity_unit, table_meta)

    def test_PK_type_invalid(self):
        """测试创建表时，第1，2，3，4个PK列type分别为DOUBLE, BOOLEAN, BINARY，期望返回OTSParameterInvalid, list_table()确认创建失败"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_PK_type_invalid'
        pk_list = []
        er_col = []
        pk_list.append([('PK0', 'DOUBLE'), ('PK1', 'STRING'), ('PK2', 'STRING'), ('PK3', 'STRING')])
        er_col.append(0)
        pk_list.append([('PK0', 'STRING'), ('PK1', 'DOUBLE'), ('PK2', 'STRING'), ('PK3', 'STRING')])
        er_col.append(1)
        pk_list.append([('PK0', 'STRING'), ('PK1', 'STRING'), ('PK2', 'DOUBLE'), ('PK3', 'STRING')])
        er_col.append(2)
        pk_list.append([('PK0', 'STRING'), ('PK1', 'STRING'), ('PK2', 'STRING'), ('PK3', 'DOUBLE')])
        er_col.append(3)
        for pk_schema, i in zip(pk_list, er_col):
            table_meta = TableMeta(table_name, pk_schema)
            reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
            try:
                self.client_test.create_table(table_meta, reserved_throughput)
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, "OTSParameterInvalid", pk_schema[i][1] + " is an invalid type for the primary key.")
            self.assert_equal(False, table_name in self.client_test.list_table())

    def test_create_table_again(self):
        """创建一个表，设置CU(1, 1), 删除它，然后用同样的Name，不同的PK创建表，设置CU为(2, 2)，期望成功，describe_table()获取信息与创建表一致，CU为(2,2)，操作验证CU"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_create_again'
        table_meta = TableMeta(table_name, [('PK0', 'INTEGER'), ('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(1, 1))
        self.client_test.create_table(table_meta, reserved_throughput)

        self.client_test.delete_table(table_name)

        table_meta_new = TableMeta(table_name, [('PK0_new', 'INTEGER'), ('PK1', 'STRING')])
        reserved_throughput_new = ReservedThroughput(CapacityUnit(2, 2))
        self.client_test.create_table(table_meta_new, reserved_throughput_new)
        self.wait_for_partition_load('table_create_again')

        describe_response = self.client_test.describe_table(table_name)
        self.assert_DescribeTableResponse(describe_response, reserved_throughput_new.capacity_unit, table_meta_new)

        pk_dict_exist = {'PK0_new': 3, 'PK1':'1'}
        pk_dict_not_exist = {'PK0_new': 5, 'PK1':'2'}
        self.check_CU_by_consuming(table_name, pk_dict_exist,  pk_dict_not_exist, reserved_throughput_new.capacity_unit)

    def test_CU_doesnot_messed_up_with_two_tables(self):
        """创建2个表，分别设置CU为(1, 2)和(2, 1)，操作验证CU，describe_table()确认设置成功"""
        time.sleep(1) # to avoid too frequently table operation
        table_name_1 = 'table1_CU_mess_up_test'
        table_meta_1 = TableMeta(table_name_1, [('PK0', 'STRING'), ('PK1', 'STRING')])
        reserved_throughput_1 = ReservedThroughput(CapacityUnit(1, 2))
        table_name_2 = 'table2_CU_mess_up_test'
        table_meta_2 = TableMeta(table_name_2, [('PK0', 'STRING'), ('PK1', 'STRING')])
        reserved_throughput_2 = ReservedThroughput(CapacityUnit(2, 1))
        pk_dict_exist = {'PK0':'a', 'PK1':'1'}
        pk_dict_not_exist = {'PK0':'B', 'PK1':'2'}
        self.client_test.create_table(table_meta_1, reserved_throughput_1)
        self.client_test.create_table(table_meta_2, reserved_throughput_2)
        self.wait_for_partition_load('table1_CU_mess_up_test')
        self.wait_for_partition_load('table2_CU_mess_up_test')

        describe_response_1 = self.client_test.describe_table(table_name_1)
        self.assert_DescribeTableResponse(describe_response_1, reserved_throughput_1.capacity_unit, table_meta_1)
        self.check_CU_by_consuming(table_name_1, pk_dict_exist,  pk_dict_not_exist, reserved_throughput_1.capacity_unit)
        describe_response_2 = self.client_test.describe_table(table_name_2)
        self.assert_DescribeTableResponse(describe_response_2, reserved_throughput_2.capacity_unit, table_meta_2)
        self.check_CU_by_consuming(table_name_2, pk_dict_exist,  pk_dict_not_exist, reserved_throughput_2.capacity_unit)

    def test_create_table_with_CU_0_0(self):
        """创建1个表，CU是(0, 0)，describe_table()确认设置成功"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_cu_0_0'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        describe_response = self.client_test.describe_table(table_name)
        self.assert_DescribeTableResponse(describe_response, reserved_throughput.capacity_unit, table_meta)

    def test_create_table_with_CU_0_1(self):
        """创建1个表，CU是(0, 1)，describe_table()确认设置成功"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_cu_0_1'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(0, 1))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        describe_response = self.client_test.describe_table(table_name)
        self.assert_DescribeTableResponse(describe_response, reserved_throughput.capacity_unit, table_meta)

    def test_create_table_with_CU_1_0(self):
        """创建1个表，CU是(1, 0)，describe_table()确认设置成功"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_cu_1_0'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(1, 0))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        describe_response = self.client_test.describe_table(table_name)
        self.assert_DescribeTableResponse(describe_response, reserved_throughput.capacity_unit, table_meta)

    def test_create_table_with_CU_1_1(self):
        """创建1个表，CU是(1, 1)，describe_table()确认设置成功"""
        time.sleep(1) # to avoid too frequently table operation
        table_name = 'table_cu_1_1'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(1, 1))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        describe_response = self.client_test.describe_table(table_name)
        self.assert_DescribeTableResponse(describe_response, reserved_throughput.capacity_unit, table_meta)

if __name__ == '__main__':
    unittest.main()
