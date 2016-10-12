# -*- coding: utf8 -*-

import unittest
from lib.ots2_api_test_base import OTS2APITestBase
from ots2 import *
import lib.restriction as restriction 
import copy
from collections import OrderedDict
from ots2.error import *
import math
import time


class RowOpTest(OTS2APITestBase):

    """行读写测试"""
    # 行操作API： GetRow, PutRow, UpdateRow, BatchGetRow, BatchWriteRow, GetRange
    # 测试每一个写操作，都要用GetRow或BatchGetRow验证数据操作符合预期
    def _check_all_row_op_with_exception_meta_not_match(self, wrong_pk):
        try:
            self.client_test.get_row('XX', wrong_pk)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSInvalidPK", "Primary key schema mismatch.")

        try:
            self.client_test.put_row('XX', Condition("IGNORE"), wrong_pk, {})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSInvalidPK", "Primary key schema mismatch.")

        try:
            self.client_test.update_row('XX', Condition("IGNORE"), wrong_pk, {'put':{'C1': 'V'}})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSInvalidPK", "Primary key schema mismatch.")

        try:
            self.client_test.delete_row('XX', Condition("IGNORE"), wrong_pk)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSInvalidPK", "Primary key schema mismatch.")
        
        response = self.client_test.batch_get_row([('XX', [wrong_pk], [])])
        eresponse = [[RowDataItem(False, "OTSInvalidPK", "Primary key schema mismatch.", None, None, None)]]
        self.assert_RowDataItem_equal(response, eresponse)

        wrong_pk_item = PutRowItem(Condition("IGNORE"), wrong_pk, {})
        response = self.client_test.batch_write_row([{'table_name': 'XX', 'put': [wrong_pk_item]}])
        eresponse = [{'put': [BatchWriteRowResponseItem(False, "OTSInvalidPK", "Primary key schema mismatch.", None)]}]
        self.assert_BatchWriteRowResponseItem(response, eresponse)

        wrong_pk_item = UpdateRowItem(Condition("IGNORE"), wrong_pk, {'put':{'C1': 'V'}})
        response = self.client_test.batch_write_row([{'table_name': 'XX', 'update': [wrong_pk_item]}])
        eresponse = [{'update': [BatchWriteRowResponseItem(False, "OTSInvalidPK", "Primary key schema mismatch.", None)]}]
        self.assert_BatchWriteRowResponseItem(response, eresponse)
        
        wrong_pk_item = DeleteRowItem(Condition("IGNORE"), wrong_pk)
        response = self.client_test.batch_write_row([{'table_name': 'XX', 'delete': [wrong_pk_item]}])
        eresponse = [{'delete': [BatchWriteRowResponseItem(False, "OTSInvalidPK", "Primary key schema mismatch.", None)]}]
        self.assert_BatchWriteRowResponseItem(response, eresponse)
        
        get_range_end = {}
        for k, v in wrong_pk.iteritems():
            get_range_end[k] = INF_MAX
        try:
            self.client_test.get_range('XX', 'FORWARD', wrong_pk, get_range_end)
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSInvalidPK", "Primary key schema mismatch.")

    def _check_all_row_op_without_exception(self, pks, cols):
        ercu = 1
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ercu, 0))
        self.assert_equal(rpks, {})
        self.assert_equal(rcols, {})

        ewcu = self.sum_CU_from_row(pks, cols)
        cu = self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)
        self.assert_consumed(cu, CapacityUnit(0, ewcu))
        ercu = ewcu
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ercu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)
        

        cu = self.client_test.update_row('XX', Condition("IGNORE"), pks, {'put':cols})
        self.assert_consumed(cu, CapacityUnit(0, ewcu))
        ercu = ewcu
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ercu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        cu = self.client_test.delete_row('XX', Condition("IGNORE"), pks)
        self.assert_consumed(cu, CapacityUnit(0, ewcu))
        ercu = 1
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ercu, 0))
        self.assert_equal(rpks, {})
        self.assert_equal(rcols, {})
        
        response = self.client_test.batch_get_row([('XX', [pks], [])])
        eresponse = [[RowDataItem(True, '', '', CapacityUnit(1, 0), {}, {})]]
        self.assert_RowDataItem_equal(response, eresponse)

        pks_item = PutRowItem(Condition("IGNORE"), pks, cols)
        ewcu = self.sum_CU_from_row(pks, cols)
        response = self.client_test.batch_write_row([{'table_name': 'XX', 'put': [pks_item]}])
        eresponse = [{'put': [BatchWriteRowResponseItem(True, '', '', CapacityUnit(0, ewcu))]}]
        self.assert_BatchWriteRowResponseItem(response, eresponse)
        ercu = ewcu
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ercu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        pks_item = UpdateRowItem(Condition("IGNORE"), pks, {'put':cols})
        ewcu = self.sum_CU_from_row(pks, cols)
        response = self.client_test.batch_write_row([{'table_name': 'XX', 'update': [pks_item]}])
        eresponse = [{'update': [BatchWriteRowResponseItem(True, '', '', CapacityUnit(0, ewcu))]}]
        self.assert_BatchWriteRowResponseItem(response, eresponse)
        ercu = ewcu
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ercu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        pks_item = DeleteRowItem(Condition("IGNORE"), pks)
        ewcu = self.sum_CU_from_row(pks, {})
        response = self.client_test.batch_write_row([{'table_name': 'XX', 'delete': [pks_item]}])
        eresponse = [{'delete': [BatchWriteRowResponseItem(True, '', '', CapacityUnit(0, ewcu))]}]
        self.assert_BatchWriteRowResponseItem(response, eresponse)
        ercu = ewcu
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ercu, 0))
        self.assert_equal(rpks, {})
        self.assert_equal(rcols, {})

        get_range_end = {}
        for k, v in pks.iteritems():
            get_range_end[k] = INF_MAX
        self.client_test.get_range('XX', 'FORWARD', pks, get_range_end)
 
    def test_pk_name_not_match(self):
        """建一个表，PK为[('PK1', 'STRING')]，测试所有行操作API，PK为{'PK2' : 'blah'}，期望返回OTSMetaNotMatch"""
        table_meta = TableMeta('XX', [('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')
 
        wrong_pk = {'PK2': 'blah'}
        self._check_all_row_op_with_exception_meta_not_match(wrong_pk)

    def test_pk_value_type_not_match(self):
        """建一个表，PK为[('PK1', 'STRING')]，测试所有行操作API，PK为{'PK2' : 123}，期望返回OTSMetaNotMatch"""
        table_meta = TableMeta('XX', [('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')
 
        wrong_pk = {'PK2': 123}
        self._check_all_row_op_with_exception_meta_not_match(wrong_pk)

    def test_pk_num_not_match(self):
        """建一个表，PK为[('PK1', 'STRING'), ('PK2', 'INTEGER')]，测试所有行操作API，PK为{'PK1' : 'blah'}或PK为{'PK1' : 'blah', 'PK2' : 123, 'PK3' : 'blah'}, 期望返回OTSMetaNotMatch"""
        table_meta = TableMeta('XX', [('PK1', 'STRING'), ('PK2', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')
 
        wrong_pk = {'PK1': 'blah'}
        self._check_all_row_op_with_exception_meta_not_match(wrong_pk)

        wrong_pk = {'PK1': 'blah', 'PK2': 123, 'PK3': 'blah'}
        self._check_all_row_op_with_exception_meta_not_match(wrong_pk)

    def test_pk_order_insensitive(self):
        """建一个表，PK为[('PK1', 'STRING'), ('PK2', 'INTEGER')]，测试所有行操作API，PK为OrderedDict([('PK2', 123), ('PK1', 'blah')])，期望相应的操作成功"""
        table_meta = TableMeta('XX', [('PK1', 'STRING'), ('PK2', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')
 
        pks = OrderedDict([('PK2', 123), ('PK1', 'blah')])
        self._check_all_row_op_without_exception(pks, {'C1': 'V'})

    def test_all_the_types(self):
        """建一个表，PK为[('PK1', 'STRING'), ('PK2', 'INTEGER')], 让所有的行操作API都操作行({'PK1' : 'blah', 'PK2' : 123}, {'C1' : 'blah', 'C2' : 123, 'C3' : True, 'C4' : False, 'C5' : 3.14, 'C6' : bytearray(1)})，期望相应的操作成功"""
        table_meta = TableMeta('XX', [('PK1', 'STRING'), ('PK2', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')
 
        pks = {'PK1' : 'blah', 'PK2' : 123}
        cols = {'C1' : 'blah', 'C2' : 123, 'C3' : True, 'C4' : False, 'C5' : 3.14, 'C6' : bytearray(1)}
        self._check_all_row_op_without_exception(pks, cols)

    def test_row_op_with_binary_type_pk(self):
        """建一个表[('PK1', 'BINARY'), ('PK2', 'BINARY')], 向表中预先插入一些数据，测试所有读写API对这些数据的读写操作，验证PK类型支持Binary之后所有API都是兼容的。"""
        table_meta = TableMeta('XX', [('PK1', 'BINARY'), ('PK2', 'BINARY')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')
 
        pks = {'PK1' : bytearray(3), 'PK2' : bytearray(2)}
        cols = {'C1' : 'blah', 'C2' : 123, 'C3' : True, 'C4' : False, 'C5' : 3.14, 'C6' : bytearray(1)}
        self._check_all_row_op_without_exception(pks, cols)

    def test_columns_to_get_semantics(self):
        """有两个表，A，B，有4个行A0, A1, B0, B1，分别在这2个表上，数据都为({'PK1' : 'blah', 'PK2' : 123}, {'C1' : 'blah', 'C2' : 123, 'C3' : True, 'C4' : False, 'C5' : 3.14, 'C6' : bytearray(1)})，让GetRow读取A0，BatchGetRow读A0, A1, B0, B1，GetRange读取A0, A1，columns_to_get参数分别为空，['C1'], ['PK1'], ['blah']，期望返回符合预期，验证CU符合预期"""
        table_meta_a = TableMeta('AA', [('PK1', 'STRING'), ('PK2', 'INTEGER')])
        table_meta_b = TableMeta('BB', [('PK1', 'STRING'), ('PK2', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit, 
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta_a, reserved_throughput)
        self.client_test.create_table(table_meta_b, reserved_throughput)
        self.wait_for_partition_load('AA')
        self.wait_for_partition_load('BB')

        pks0 = {'PK1': '0', 'PK2': 123}
        pks1 = {'PK1': '1', 'PK2': 123}
        cols = {'C1': 'blah', 'C2': 123, 'C3': True, 'C4': False, 'C5': 3.14, 'C6': bytearray(1)}
        putrow_item0 = PutRowItem(Condition("IGNORE"), pks0, cols)
        putrow_item1 = PutRowItem(Condition("IGNORE"), pks1, cols)
        self.client_test.batch_write_row([
            {'table_name': 'AA', 'put': [putrow_item0, putrow_item1]}, 
            {'table_name': 'BB', 'put': [putrow_item0, putrow_item1]}
        ])

        cu, rpks, rcols = self.client_test.get_row('AA', pks0)
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(rpks, pks0)
        self.assert_equal(rcols, cols)

        cu, rpks, rcols = self.client_test.get_row('AA', pks0, ['C1'])
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(rpks, {})
        self.assert_equal(rcols, {'C1': 'blah'})

        cu, rpks, rcols = self.client_test.get_row('AA', pks0, ['PK1'])
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(rpks, {'PK1': '0'})
        self.assert_equal(rcols, {})

        cu, rpks, rcols = self.client_test.get_row('AA', pks0, ['blah'])
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(rpks, {})
        self.assert_equal(rcols, {})
        
        response = self.client_test.batch_get_row([
            ('AA', [pks0, pks1], []), 
            ('BB', [pks0, pks1], [])
        ])
        rowdata_item0 = RowDataItem(True, '', '', CapacityUnit(1, 0), pks0, cols)
        rowdata_item1 = RowDataItem(True, '', '', CapacityUnit(1, 0), pks1, cols)
        eresponse = [[rowdata_item0, rowdata_item1], [rowdata_item0, rowdata_item1]]
        self.assert_RowDataItem_equal(response, eresponse)

        response = self.client_test.batch_get_row([
            ('AA', [pks0, pks1], ['C1']), 
            ('BB', [pks0, pks1], ['C1'])
        ])
        rowdata_item = RowDataItem(True, '', '', CapacityUnit(1, 0), {}, {'C1': 'blah'})
        eresponse = [[rowdata_item, rowdata_item], [rowdata_item, rowdata_item]]
        self.assert_RowDataItem_equal(response, eresponse)

        response = self.client_test.batch_get_row([
            ('AA', [pks0, pks1], ['PK1']), 
            ('BB', [pks0, pks1], ['PK1'])
        ])
        rowdata_item0 = RowDataItem(True, '', '', CapacityUnit(1, 0), {'PK1': '0'}, {})
        rowdata_item1 = RowDataItem(True, '', '', CapacityUnit(1, 0), {'PK1': '1'}, {})
        eresponse = [[rowdata_item0, rowdata_item1], [rowdata_item0, rowdata_item1]]
        self.assert_RowDataItem_equal(response, eresponse)

        response = self.client_test.batch_get_row([
            ('AA', [pks0, pks1], ['blah']), 
            ('BB', [pks0, pks1], ['blah'])
        ])
        rowdata_item = RowDataItem(True, '', '', CapacityUnit(1, 0), {}, {})
        eresponse = [[rowdata_item, rowdata_item], [rowdata_item, rowdata_item]]
        self.assert_RowDataItem_equal(response, eresponse)

        cu, next_pks, rows = self.client_test.get_range('AA', 'FORWARD', 
                {'PK1': INF_MIN, 'PK2': INF_MIN}, 
                {'PK1': INF_MAX, 'PK2': INF_MAX},
                []
        )
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(next_pks, None)
        self.assert_equal(rows, [(pks0, cols), (pks1, cols)])

        cu, next_pks, rows = self.client_test.get_range('AA', 'FORWARD', 
                {'PK1': INF_MIN, 'PK2': INF_MIN}, 
                {'PK1': INF_MAX, 'PK2': INF_MAX},
                ['C1']
        )
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(next_pks, None)
        self.assert_equal(rows, [({}, {'C1': 'blah'}), ({}, {'C1': 'blah'})])

        cu, next_pks, rows = self.client_test.get_range('AA', 'FORWARD', 
                {'PK1': INF_MIN, 'PK2': INF_MIN}, 
                {'PK1': INF_MAX, 'PK2': INF_MAX},
                ['PK1']
        )
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(next_pks, None)
        self.assert_equal(rows, [({'PK1': '0'}, {}), ({'PK1': '1'}, {})])
           
        cu, next_pks, rows = self.client_test.get_range('AA', 'FORWARD', 
                {'PK1': INF_MIN, 'PK2': INF_MIN}, 
                {'PK1': INF_MAX, 'PK2': INF_MAX},
                ['PK1', 'C1']
        )
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(next_pks, None)
        self.assert_equal(rows, [({'PK1': '0'}, {'C1': 'blah'}), ({'PK1': '1'}, {'C1': 'blah'})])
           
        cu, next_pks, rows = self.client_test.get_range('AA', 'FORWARD', 
                {'PK1': INF_MIN, 'PK2': INF_MIN}, 
                {'PK1': INF_MAX, 'PK2': INF_MAX},
                ['blah']
        )
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(next_pks, None)
        self.assert_equal(rows, [])

    def test_CU_consumed_for_whole_row(self):
        """有一行，数据为({'PK0' : 'blah'}, {'C1' : 500B, 'C2' : 500B})，读整行，或者只读PK0, C1, C2，期望消耗的CU为(2, 0)或(1, 0)"""
        table_meta = TableMeta('XX', [('PK0', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')

        pks = {'PK0': 'blah'}
        cols = {'C1': 'V' * 512, 'C2': 'X' * 512}
        self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)

        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        cu, rpks, rcols = self.client_test.get_row('XX', pks, ['PK0'])
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, {})

        cu, rpks, rcols = self.client_test.get_row('XX', pks, ['C1'])
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(rpks, {})
        self.assert_equal(rcols, {'C1': 'V' * 512})

        cu, rpks, rcols = self.client_test.get_row('XX', pks, ['C2'])
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(rpks, {})
        self.assert_equal(rcols, {'C2': 'X' * 512})

    def test_get_row_miss(self):
        """GetRow读一个不存在的行, 期望返回为空, 验证消耗CU(1, 0)"""
        table_meta = TableMeta('XX', [('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')
        
        cu_expect = CapacityUnit(1, 0)
        cu, rpks, rcols = self.client_test.get_row('XX', {'PK1': 'b'})
        self.assert_consumed(cu, cu_expect)
        self.assert_equal(rpks, {})
        self.assert_equal(rcols, {})

    def _test_get_row_error_with_CU_consumed(self):
        """GetRow在后端出错（How?），从后端验证消耗了CU(1, 0)"""
        # 这个case暂时不实现
        raise NotImplementedError

    def test_update_row_when_row_exist(self):
        """原有的行包含列 {'C0' : 2k, 'C1' : 2k}（2k指的是不重复的2k STRING），数据量小于8K大于4K，分别测试UpdateRow，row existence expectation为EXIST, IGNORE时，列值为分别为：{'C0' : 2k, 'C1' : 2k}（覆盖），{'C2' : 2k, 'C3' : 2k}（添加），{'C0' : Delete, 'C1' : Delete}（删除），{'C0' : 2k, 'C1' : Delete, 'C2' : 2k}(交错)。每次都用GetRow检查数据是否符合预期，并检查CU消耗是否正确"""
        table_meta = TableMeta('XX', [('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')

        pks = {'PK1': '0' * 20}
        cols = {'C0': 'V' * 2048, 'C1': 'B' * 2048}
        self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)

        #EXIST+COVER
        cu = self.client_test.update_row('XX', Condition("EXPECT_EXIST"), pks, { 'put' : cols })
        self.assert_consumed(cu, CapacityUnit(1, 2))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(2, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        #EXIST+ADD
        cols = {'C2': 'V' * 2048, 'C3': 'B' * 2048}
        cu = self.client_test.update_row('XX', Condition("EXPECT_EXIST"), pks, { 'put' : cols })
        self.assert_consumed(cu, CapacityUnit(1, 2))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(3, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, {'C0': 'V' * 2048, 'C1': 'B' * 2048, 'C2': 'V' * 2048, 'C3': 'B' * 2048})

        #EXIST+DELETE
        cu = self.client_test.update_row('XX', Condition("EXPECT_EXIST"), pks, { 'delete' : ['C0', 'C1'] })
        self.assert_consumed(cu, CapacityUnit(1, 1))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(2, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, {'C2': 'V' * 2048, 'C3': 'B' * 2048})

        #EXIST+MIX
        cols = {'C0': 'V' * 2048, 'C1': None, 'C2': 'V' * 2048}
        cu = self.client_test.update_row('XX', Condition("EXPECT_EXIST"), pks, {'put' : {'C0': 'V' * 2048, 'C2': 'V' * 2048}, 'delete' : ['C1']})
        self.assert_consumed(cu, CapacityUnit(1, 2))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(2, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, {'C0': 'V' * 2048, 'C2': 'V' * 2048, 'C3': 'B' * 2048})

        pks = {'PK1': '0' * 20}
        cols = {'C0': 'V' * 2048, 'C1': 'B' * 2048}
        self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)

        #IGNORE+COVER
        cu = self.client_test.update_row('XX', Condition("IGNORE"), pks, {'put':cols})
        self.assert_consumed(cu, CapacityUnit(0, 2))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(2, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        #IGNORE+ADD
        cols = {'C2': 'V' * 2048, 'C3': 'B' * 2048}
        cu = self.client_test.update_row('XX', Condition("IGNORE"), pks, {'put':cols})
        self.assert_consumed(cu, CapacityUnit(0, 2))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(3, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, {'C0': 'V' * 2048, 'C1': 'B' * 2048, 'C2': 'V' * 2048, 'C3': 'B' * 2048})

        #IGNORE+DELETE
        cu = self.client_test.update_row('XX', Condition("IGNORE"), pks, {'delete':['C0', 'C1']})
        self.assert_consumed(cu, CapacityUnit(0, 1))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(2, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, {'C2': 'V' * 2048, 'C3': 'B' * 2048})

        #IGNORE+MIX
        cu = self.client_test.update_row('XX', Condition("IGNORE"), pks, {'put':{'C0': 'V' * 2048, 'C2': 'V' * 2048}, 'delete':['C1']})
        self.assert_consumed(cu, CapacityUnit(0, 2))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(2, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, {'C0': 'V' * 2048, 'C2': 'V' * 2048, 'C3': 'B' * 2048})

    def test_update_row_when_value_type_changed(self):
        """原有的行包含max个列，值分别为INTEGER, DOUBLE, STRING(8 byte), BOOLEAN, BINARY(8 byte)，测试PutRow包含max个列，值分别为INTEGER, DOUBLE, STRING(8 byte), BOOLEAN, BINARY(8 byte)，每次GetRow检查数据，并验证CU消耗正常。"""
        table_meta = TableMeta('XX', [('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')

        pks = {'PK1': '0'}
        cols = {}
        for i in range(0, restriction.MaxColumnCountForRow):
            cols['C' + str(i)] = 1
        ewcu = self.sum_CU_from_row(pks, cols)
        cu = self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)
        self.assert_consumed(cu, CapacityUnit(0, ewcu))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ewcu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        for k in cols.keys():
            cols[k] = 1.0
        cu = self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)
        self.assert_consumed(cu, CapacityUnit(0, ewcu))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ewcu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        for k in cols.keys():
            cols[k] = 'V' * 8 
        cu = self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)
        self.assert_consumed(cu, CapacityUnit(0, ewcu))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ewcu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        for k in cols.keys():
            cols[k] = True 
        ewcu = self.sum_CU_from_row(pks, cols)
        cu = self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)
        self.assert_consumed(cu, CapacityUnit(0, ewcu))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ewcu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

        for k in cols.keys():
            cols[k] = bytearray('V' * 8) 
        ewcu = self.sum_CU_from_row(pks, cols)
        cu = self.client_test.put_row('XX', Condition("IGNORE"), pks, cols)
        self.assert_consumed(cu, CapacityUnit(0, ewcu))
        cu, rpks, rcols = self.client_test.get_row('XX', pks)
        self.assert_consumed(cu, CapacityUnit(ewcu, 0))
        self.assert_equal(rpks, pks)
        self.assert_equal(rcols, cols)

    def test_update_row_but_expect_row_not_exist(self):
        """UpdateRow的row existence expectation为NOT_EXIST，期望返回OTSParameterInvalid"""
        table_meta = TableMeta('XX', [('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')

        try:
            self.client_test.update_row('XX', Condition("EXPECT_NOT_EXIST"), {'PK1': '0'}, {'put':{'Col0' : 'XXXX'}})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "Invalid condition: EXPECT_NOT_EXIST while updating row.")

    def test_delete_row_but_expect_row_not_exist(self):
        """DeleteRow的row existence expectation为NOT_EXIST，期望返回OTSParameterInvalid"""
        table_meta = TableMeta('XX', [('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))

        try:
            self.client_test.delete_row('XX', Condition("EXPECT_NOT_EXIST"), {'PK1': '0'})
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "Invalid condition: EXPECT_NOT_EXIST while deleting row.")

    def test_get_range_less_than_4K(self):
        """GetRange包含10行，数据不超过4K，或大于4K小于8K，期望消耗CU(1, 0)或者(2, 0)"""
        table_meta = TableMeta('XX', [('PK1', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(
            restriction.MaxReadWriteCapacityUnit,
            restriction.MaxReadWriteCapacityUnit
        ))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('XX')

        rowitems = []
        for i in range(0, 9):
            pk = {'PK1': str(i)}
            col = {'C': 'V' * 400}
            putrow_item = PutRowItem(Condition("IGNORE"), pk, col)
            rowitems.append(putrow_item)
        self.client_test.batch_write_row([{'table_name': 'XX', 'put': rowitems}])

        cu_sum = CapacityUnit(0, 0)
        i = 0
        for row in self.client_test.xget_range('XX', 'FORWARD', {'PK1': INF_MIN}, {'PK1': INF_MAX}, cu_sum):
            epk = {'PK1': str(i)}
            ecol = {'C': 'V' * 400}
            self.assert_equal(row, (epk, ecol))
            i += 1
        self.assert_consumed(cu_sum, CapacityUnit(1, 0))

        rowitems = []
        for i in range(0, 9):
            pk = {'PK1': str(i)}
            col = {'C': 'V' * 800}
            putrow_item = PutRowItem(Condition("IGNORE"), pk, col)
            rowitems.append(putrow_item)
        self.client_test.batch_write_row([{'table_name': 'XX', 'put': rowitems}])

        cu_sum = CapacityUnit(0, 0)
        i = 0
        for row in self.client_test.xget_range('XX', 'FORWARD', {'PK1': INF_MIN}, {'PK1': INF_MAX}, cu_sum):
            epk = {'PK1': str(i)}
            ecol ={ 'C': 'V' * 800}
            self.assert_equal(row, (epk, ecol))
            i += 1
        self.assert_consumed(cu_sum, CapacityUnit(2, 0))

    def _valid_column_name_test(self, table_name, pk, columns):
        #put_row
        consumed = self.client_test.put_row(table_name, Condition("IGNORE"), pk, columns)
        self.assert_consumed(consumed, CapacityUnit(0, self.sum_CU_from_row(pk, columns)))
        #get_row
        consumed, primary_keys, columns_res = self.client_test.get_row(table_name, pk)
        self.assert_consumed(consumed, CapacityUnit(self.sum_CU_from_row(pk, columns), 0))
        self.assert_equal(primary_keys, pk)
        self.assert_equal(columns_res, columns)
        #batch_get_row
        batches = [(table_name, [pk], [])]
        response = self.client_test.batch_get_row(batches)
        expect_row_data_item = RowDataItem(True, "", "", CapacityUnit(self.sum_CU_from_row(pk, columns), 0), pk, columns)
        eresponse = [[expect_row_data_item]]
        self.assert_RowDataItem_equal(response, eresponse)
        #get_range
        inclusive_key = copy.copy(pk)
        exclusive_key = copy.copy(pk)
        for k in pk.keys():
            inclusive_key[k] = INF_MIN
            exclusive_key[k] = INF_MAX
        consumed, next_start_primary_keys, rows = self.client_test.get_range(table_name, 'FORWARD', inclusive_key,exclusive_key, [], None)
        self.assert_consumed(consumed, CapacityUnit(self.sum_CU_from_row(pk, columns), 0))
        self.assert_equal(next_start_primary_keys, None)
        self.assert_equal(rows, [(pk, columns)])

        #update_row
        consumed = self.client_test.update_row(table_name, Condition("IGNORE"), pk, {'put':columns})
        self.assert_consumed(consumed, CapacityUnit(0, self.sum_CU_from_row(pk, columns)))

        #delete_row
        consumed = self.client_test.delete_row(table_name, Condition("IGNORE"), pk)
        self.assert_consumed(consumed, CapacityUnit(0, self.sum_CU_from_row(pk, {})))
        #batch_write_row
        put_row_item = PutRowItem(Condition("IGNORE"), pk, columns)
        update_row_item = UpdateRowItem(Condition("IGNORE"), pk, {'put':columns})
        delete_row_item = DeleteRowItem(Condition("IGNORE"), pk)
        batches_list = [{'put':[put_row_item]}, {'update':[update_row_item]}, {'delete':[delete_row_item]}]
        expect_write_data_item = BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, self.sum_CU_from_row(pk, columns)))
        response_list = [
            {'put':[BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, self.sum_CU_from_row(pk, columns)))]}, 
            {'update':[BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, self.sum_CU_from_row(pk, columns)))]}, 
            {'delete':[BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, self.sum_CU_from_row(pk, {})))]}, 
        ]
        for i in range(len(batches_list)):
            write_batches = batches_list[i]
            write_batches['table_name'] = table_name
            response = self.client_test.batch_write_row([write_batches])
            eresponse = [response_list[i]]
            self.assert_BatchWriteRowResponseItem(response, eresponse)

    def test_max_data_size_row(self):
        """对于每一个行操作API，测试一个数据量最大的行， 主键为max个，每个主键名字长度为max，value为string长度为max，列的数据为max byte，类型分别为INTEGER, BOOLEAN, BINARY, STRING, DOUBLE, 判断消耗的CU符合预期"""
        #这里假设MaxPKColumnNum <=10 ,不然PKname的长度不能符合要求
        pk_schema, pk = self.get_primary_keys(restriction.MaxPKColumnNum, "STRING", pk_name="P" * (restriction.MaxColumnNameLength - 1), pk_value="x" * (restriction.MaxPKStringValueLength))

        table_name = "table_test"
        table_meta = TableMeta(table_name, pk_schema)
        reserved_throughput = ReservedThroughput(CapacityUnit(restriction.MaxReadWriteCapacityUnit, restriction.MaxReadWriteCapacityUnit))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test')
        
        remained_size = restriction.MaxColumnDataSizeForRow - restriction.MaxPKColumnNum * (restriction.MaxColumnNameLength + restriction.MaxPKStringValueLength)
        remained_column_cnt = restriction.MaxColumnCountForRow - restriction.MaxPKColumnNum

        integer_columns = {}
        string_columns = {}
        bool_columns = {}
        binary_columns = {}
        double_columns = {}
        #INTEGER and DOUBLE
        col_key_tmp = "a%0" + str(restriction.MaxColumnNameLength - 1) + "d"
        if (remained_size / (restriction.MaxColumnNameLength + 8)) > remained_column_cnt:
            for i in range(remained_column_cnt):
                integer_columns[col_key_tmp%i] = i
                double_columns[col_key_tmp%i] = i + 0.1
        else:
            for i in range(remained_size / (restriction.MaxColumnNameLength + 8)):
                integer_columns[col_key_tmp%i] = i
                double_columns[col_key_tmp%i] = i + 0.1
        #BOOL
        if (remained_size / (restriction.MaxColumnNameLength + 1)) > remained_column_cnt:
            for i in range(remained_column_cnt):
                bool_columns[col_key_tmp%i] = True
        else:
            for i in range(remained_size / (restriction.MaxColumnNameLength + 1)):
                bool_columns[col_key_tmp%i] = False
        #string
        for i in range(remained_size / (restriction.MaxColumnNameLength + restriction.MaxNonPKStringValueLength)):
            string_columns[col_key_tmp%i] = "X" * restriction.MaxNonPKStringValueLength
        #binary
        for i in range(remained_size / (restriction.MaxColumnNameLength + restriction.MaxBinaryValueLength)):
            binary_columns[col_key_tmp%i] = bytearray(restriction.MaxBinaryValueLength)
        
        for col in [integer_columns, string_columns, bool_columns, binary_columns, double_columns]:
            self._valid_column_name_test(table_name, pk, col)

    def test_batch_get_on_the_same_row(self):
        """创建一个表T，一个行R，数据量为 < 1KB，BatchGetRow([(T, [R, R])])，重复行，期望返回OTSParameterInvalid，再一次BatchGetRow([(T, [R]), (T, [R]])，同名表在不同组，期望返回OTSParameterInvalid"""
        table_name = 'table_test_batch_get_on_the_same_row'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'INTEGER')])
        pk_dict = {'PK0':'a', 'PK1':1}
        column_dict = {'col1': 'M' * 500}
        reserved_throughput = ReservedThroughput(CapacityUnit(10, 10))

        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test_batch_get_on_the_same_row')

        consumed = self.client_test.put_row(table_name, Condition("IGNORE"), pk_dict, column_dict)
        consumed_expect = CapacityUnit(0, 1)
        self.assert_consumed(consumed, consumed_expect)

        try:
            response = self.client_test.batch_get_row([(table_name, [pk_dict, pk_dict], [])])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "The input parameter is invalid.")

        try:
            response = self.client_test.batch_get_row([(table_name, [pk_dict], []), (table_name, [pk_dict], [])])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "Duplicated table name: 'table_test_batch_get_on_the_same_row'.")

    def test_batch_write_on_the_same_row(self):
        """BatchWriteRow，分别为put, delete, update，操作在同一行（在同一个表名下，或者重复的两个表名下），期望返回OTSParameterInvalid"""
        table_name = 'table_test_batch_write_on_the_same_row'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'INTEGER')])
        pk_dict = {'PK0':'a', 'PK1':1}
        pk_dict_2 = {'PK0':'a', 'PK1':2}
        column_dict = {'col1': 'M'}
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))

        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test_batch_write_on_the_same_row')

        put_row_item = PutRowItem(Condition("IGNORE"), pk_dict, {'col1':150})
        update_row_item = UpdateRowItem(Condition("IGNORE"), pk_dict, {'put':{'col1':200}})
        delete_row_item = DeleteRowItem(Condition("IGNORE"), pk_dict)
        test_batch_write_row_list = [{'put':[put_row_item], 'update':[update_row_item]}, {'put':[put_row_item], 'delete':[delete_row_item]}, {'update':[update_row_item], 'delete':[delete_row_item]}]

        for i in range(len(test_batch_write_row_list)):
            try:
                table_item = test_batch_write_row_list[i]
                table_item['table_name'] = table_name
                write_response = self.client_test.batch_write_row([table_item])
                self.assert_false()
            except OTSServiceError as e:
                self.assert_error(e, 400, "OTSParameterInvalid", "The input parameter is invalid.")

        update_row_item_2 = UpdateRowItem(Condition("IGNORE"), pk_dict_2, {'put':{'col1':200}})
        try:
            write_response = self.client_test.batch_write_row([{'table_name':table_name, 'put':[put_row_item]}, {'table_name':table_name, 'update':[update_row_item_2]}])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "Duplicated table name: table_test_batch_write_on_the_same_row.")
        try:
            write_response = self.client_test.batch_write_row([{'table_name':table_name, 'update':[update_row_item]}, {'table_name':table_name, 'update':[update_row_item_2]}])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "Duplicated table name: table_test_batch_write_on_the_same_row.")
        try:
            write_response = self.client_test.batch_write_row([{'table_name':table_name, 'delete':[delete_row_item]}, {'table_name':table_name, 'update':[update_row_item_2]}])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "Duplicated table name: table_test_batch_write_on_the_same_row.")


    def test_no_item_in_batch_ops(self):
        """BatchGetRow和BatchWriteRow没有包含任何行，期望返回OTSParameterInvalid"""
        table_name = 'table_test_no_item_in_batch_ops'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))

        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test_no_item_in_batch_ops')

        try:
            write_response = self.client_test.batch_write_row([{'table_name':table_name}])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "No row specified in table: 'table_test_no_item_in_batch_ops'.")

        try:
            response = self.client_test.batch_get_row([(table_name, [], [])])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "No row specified in table: 'table_test_no_item_in_batch_ops'.")

    def test_no_table_in_batch_ops(self):
        """BatchGetRow和BatchWriteRow没有包含任何行，期望返回OTSParameterInvalid"""
        table_name = 'table_test_no_item_in_batch_ops'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))

        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test_no_item_in_batch_ops')

        try:
            write_response = self.client_test.batch_write_row([])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "No row specified in the request of BatchWriteRow.")

        try:
            response = self.client_test.batch_get_row([])
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, "OTSParameterInvalid", "No row specified in the request of BatchGetRow.")

    def test_get_range_when_direction_is_wrong_for_1_PK(self):
        """一个表有1个PK，测试方向为FORWARD/FORWARD，第一个begin(大于或等于)/(小于或等于)end，PK类型分别为STRING, INTEGER的情况，期望返回OTSParameterInvalid"""
        table_name_string = 'table_test_get_range_when_direction_string'
        table_meta_string = TableMeta(table_name_string, [('PK0', 'STRING')])
        table_name_integer = 'table_test_get_range_when_direction_integer'
        table_meta_integer = TableMeta(table_name_integer, [('PK0', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))

        self.client_test.create_table(table_meta_string, reserved_throughput)
        self.client_test.create_table(table_meta_integer, reserved_throughput)
        self.wait_for_partition_load('table_test_get_range_when_direction_string')
        self.wait_for_partition_load('table_test_get_range_when_direction_integer')

        pk_dict_small = {'PK0':'AAA'}
        pk_dict_big = {'PK0':'DDDD'}
        self._check_get_range_expect_exception(table_name_string, pk_dict_big, pk_dict_small, 'FORWARD', "The input parameter is invalid.")
        self._check_get_range_expect_exception(table_name_string, pk_dict_small, pk_dict_big, 'BACKWARD', "The input parameter is invalid.")

        pk_dict_small = {'PK0':10}
        pk_dict_big = {'PK0':90}
        self._check_get_range_expect_exception(table_name_integer, pk_dict_big, pk_dict_small, 'FORWARD', "The input parameter is invalid.")
        self._check_get_range_expect_exception(table_name_integer, pk_dict_small, pk_dict_big, 'BACKWARD', "The input parameter is invalid.")

    def _check_get_range_expect_exception(self, table_name, inclusive_start_primary_keys, exclusive_end_primary_keys, direction, error_msg, error_code="OTSParameterInvalid", limit=None, columns_to_get=[]):
        try:
            consumed, next_start_primary_keys, rows = self.client_test.get_range(table_name, direction, inclusive_start_primary_keys,
                                                                                 exclusive_end_primary_keys, columns_to_get, limit)
            print inclusive_start_primary_keys
            print exclusive_end_primary_keys
            self.assert_false()
        except OTSServiceError as e:
            self.assert_error(e, 400, error_code, error_msg)

    def test_all_the_ranges_for_2_PK(self):
        """一个表有2个PK, partition key 为 a < b < c < d，可以是STRING, INTEGER(分别测试)，每个partitionkey有2个行ax, ay, bx, by, cx, cy，其中x, y为第二个PK的值，分别为STRING, INTEGER, BOOlEAN。分别测试正反方向(测试反方向时把begine和end互换)的get_range: (a MIN, b MAX)，(b MAX, a MIN)（出错），(b MIN, a MAX)出错，(a MIN, a MAX), (a MAX, a MIN)（出错), (a MAX, b MIN), (a MAX, c MIN), (b x, a x)（出错）, (a x, a y), (a MIN, a y), (a x, a MAX), (a x, c x), (a x, c y), (a y, c x), (a x, a x)（出错），每个成功的操作检查数据返回符合期望，CU消耗符合期望"""
        for first_pk_type in ('STRING', 'INTEGER'):
            for second_pk_type in ('STRING', 'INTEGER'):
                table_name = 'T' + first_pk_type + second_pk_type
                table_meta = TableMeta(table_name, [('PK0', first_pk_type), ('PK1', second_pk_type)])
                self.client_test.create_table(table_meta, ReservedThroughput(CapacityUnit(100, 100)))
                self.wait_for_partition_load(table_name)

                if first_pk_type == 'STRING':
                    a, b, c = 'A', 'B', 'C'
                else:
                    a, b, c = 1, 2, 3

                if second_pk_type == 'STRING':
                    x, y = 'A', 'B'
                else:
                    x, y = 1, 2

                #pk_list = [{'PK0': 'a', 'PK1': 1}, {'PK0': 'a', 'PK1': 2}, {'PK0': 'b', 'PK1': 1}, 
                #           {'PK0': 'b', 'PK1': 2}, {'PK0': 'c', 'PK1': 1}, {'PK0': 'c', 'PK1': 2}]
                #
                row_list = []
                for first_pk_value in [a, b, c]:
                    for second_pk_value in [x, y]:
                        row_list.append(({'PK0' : first_pk_value, 'PK1' : second_pk_value}, {}))
                        self.client_test.put_row(table_name, Condition("IGNORE"), {'PK0' : first_pk_value, 'PK1' : second_pk_value}, {})

                def get_range(begin_pk0, begin_pk1, end_pk0, end_pk1):
                    return [{'PK0' : begin_pk0, 'PK1' : begin_pk1}, {'PK0' : end_pk0, 'PK1' : end_pk1}]

                range_list = [
                    # range                             是否正常 期望rows        
                    (get_range(b, x,       a, x),       False,   row_list[0:1]          ),
                    (get_range(a, x,       a, y),       True,    row_list[0:2]  ),
                    (get_range(a, x,       c, x),       True,    row_list[0:5]),
                    (get_range(a, x,       c, y),       True,    row_list[0:6]),
                    (get_range(a, y,       c, x),       True,    row_list[1:5]),
                    (get_range(a, x,       a, x),       False,   row_list[0:1]          ),
                ]

                range_list_1 = [
                    (get_range(a, INF_MAX, b, INF_MIN), True,    row_list[0:0]          ),
                    (get_range(b, INF_MAX, a, INF_MIN), False,   row_list[0:0]          ),
                    (get_range(b, INF_MAX, a, INF_MAX), False,   row_list[0:0]          ),
                    (get_range(a, INF_MIN, a, INF_MAX), True,    row_list[0:2]),
                    (get_range(a, INF_MAX, a, INF_MIN), False,   row_list[0:0]          ),
                    (get_range(a, INF_MAX, b, INF_MIN), True,    row_list[0:0]          ),
                    (get_range(a, INF_MAX, c, INF_MIN), True,    row_list[2:4]),
                ]

                range_list_2 = [
                    (get_range(a, INF_MIN, a, y),       True,    row_list[0:2]  ),
                    (get_range(a, x,       a, INF_MAX), True,    row_list[0:3]),
                ]
                begin, end = get_range(a, INF_MIN, a, y)
                self._check_xget_range(table_name, begin, end, 'FORWARD', row_list[0:1])
                self._check_xget_range(table_name, end, begin, 'BACKWARD', row_list[0:2])
                begin, end = get_range(a, x, a, INF_MAX)
                self._check_xget_range(table_name, begin, end, 'FORWARD', row_list[0:2])
                self._check_xget_range(table_name, end, begin, 'BACKWARD', row_list[1:2])

                for direction in ['FORWARD', 'BACKWARD']:
                    for range_, is_normal, _expect_rows in range_list:

                        if direction == 'FORWARD':
                            begin, end = range_
                            expect_rows = _expect_rows[0:len(_expect_rows) - 1]
                        else:
                            end, begin = range_
                            expect_rows = _expect_rows[1:(len(_expect_rows))]
                        self.logger.debug("%s  %s  %s  %s" %(begin, end, direction, expect_rows))
                        if is_normal:
                            self._check_xget_range(table_name, begin, end, direction, expect_rows)
                        else:
                            self._check_get_range_expect_exception(table_name, begin, end, direction, "The input parameter is invalid.")

                    for range_, is_normal, _expect_rows in range_list_1:
                        expect_rows = _expect_rows
                        if direction == 'FORWARD':
                            begin, end = range_
                        else:
                            end, begin = range_
                        self.logger.debug("%s  %s  %s  %s" %(begin, end, direction, expect_rows))
                        if is_normal:
                            self._check_xget_range(table_name, begin, end, direction, expect_rows)
                        else:
                            self._check_get_range_expect_exception(table_name, begin, end, direction, "The input parameter is invalid.")

    def test_4_PK_range(self):
        """一个表有4个PK，类型分别是STRING, STRING, INTEGER，测试range: ('A' 'A' 10 False, 'A' 'A' 10 True), ('A' 'A' 10 False, 'A' 'A' 11 True),  ('A' 'A' 10 False, 'A' 'A' 9 True)（出错）, ('A' 'A' 10 MAX, 'A' 'B' 10 MIN), ('A' MIN 10 False, 'B' MAX 2 True)，构造数据让每个区间都有值"""
        table_name = 'table_test_4_PK_range'
        table_meta = TableMeta(table_name, [('PK0', 'STRING'), ('PK1', 'STRING'), ('PK2', 'INTEGER'), ('PK3', 'INTEGER')])
        pk_dict_list = [{'PK0':'A', 'PK1':'A', 'PK2':9, 'PK3':9}, 
                   {'PK0':'A', 'PK1':'A', 'PK2':10, 'PK3':1},
                   {'PK0':'A', 'PK1':'A', 'PK2':10, 'PK3':9},
                   {'PK0':'A', 'PK1':'A', 'PK2':11, 'PK3':9}]
        row_list = []
        putrow_list = []
        for pk_dict in pk_dict_list:
            putrow_list.append(PutRowItem(Condition("IGNORE"), pk_dict, {}))
            row_list.append((pk_dict, {}))
        reserved_throughput = ReservedThroughput(CapacityUnit(100, 100))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test_4_PK_range')

        self.client_test.batch_write_row([{'table_name':table_name, 'put':putrow_list}])

        pk_dict_min = {'PK0':'A', 'PK1':INF_MIN, 'PK2':10, 'PK3':1}
        pk_dict3 = {'PK0':'A', 'PK1':'A', 'PK2':10, 'PK3':INF_MAX}
        pk_dict5 = {'PK0':'A', 'PK1':'B', 'PK2':10, 'PK3':INF_MIN}
        pk_dict_max = {'PK0':'B', 'PK1':INF_MAX, 'PK2':2, 'PK3':9}

        self._check_get_range_expect_exception(table_name, pk_dict_list[1], pk_dict_list[0], 'FORWARD',  "The input parameter is invalid.")

        self._check_xget_range(table_name, pk_dict_list[1], pk_dict_list[2], 'FORWARD', row_list[1:2])
        self._check_xget_range(table_name, pk_dict_list[1], pk_dict_list[3], 'FORWARD', row_list[1:3])
        self._check_xget_range(table_name, pk_dict3, pk_dict5, 'FORWARD', row_list[3:4])
        self._check_xget_range(table_name, pk_dict_min, pk_dict_max, 'FORWARD', row_list[0:4])

    def test_empty_range(self):
        """分别测试PK个数为1，2，3，4的4个表，range中包含的row个数为0的情况，期望返回为空，CU消耗为(1, 0)"""
        table_name_list = ['table_test_empty_range_0', 'table_test_empty_range_1', 'table_test_empty_range_2', 'table_test_empty_range_3']
        pk_schema0, pk_dict0_exclusive = self.get_primary_keys(1, 'STRING', 'PK', INF_MAX)
        pk_schema1, pk_dict1_exclusive = self.get_primary_keys(2, 'STRING', 'PK', INF_MAX)
        pk_schema2, pk_dict2_exclusive = self.get_primary_keys(3, 'STRING', 'PK', INF_MAX)
        pk_schema3, pk_dict3_exclusive = self.get_primary_keys(4, 'STRING', 'PK', INF_MAX)
        pk_schema0, pk_dict0_inclusive = self.get_primary_keys(1, 'STRING', 'PK', INF_MIN)
        pk_schema1, pk_dict1_inclusive = self.get_primary_keys(2, 'STRING', 'PK', INF_MIN)
        pk_schema2, pk_dict2_inclusive = self.get_primary_keys(3, 'STRING', 'PK', INF_MIN)
        pk_schema3, pk_dict3_inclusive = self.get_primary_keys(4, 'STRING', 'PK', INF_MIN)
        pk_schema_list = [pk_schema0, pk_schema1, pk_schema2, pk_schema3]
        pk_dict_inclusive_list = [pk_dict0_inclusive, pk_dict1_inclusive, pk_dict2_inclusive, pk_dict3_inclusive]
        pk_dict_exclusive_list = [pk_dict0_exclusive, pk_dict1_exclusive, pk_dict2_exclusive, pk_dict3_exclusive]
        reserved_throughput = ReservedThroughput(CapacityUnit(10, 9))

        for i in range(len(table_name_list)):
            table_meta = TableMeta(table_name_list[i], pk_schema_list[i])
            self.client_test.create_table(table_meta, reserved_throughput)
            self.wait_for_partition_load(table_name_list[i])
            consumed, next_start_primary_keys, rows = self.client_test.get_range(table_name_list[i], 'FORWARD',  pk_dict_inclusive_list[i], pk_dict_exclusive_list[i])
            self.assert_consumed(consumed, CapacityUnit(1, 0))
            self.assert_equal(next_start_primary_keys, None)
            self.assert_equal(rows, [])

    def test_get_range_limit_invalid(self):
        """测试get_range的limit为0或-1，期望返回错误OTSParameterInvalid"""
        table_name = 'table_test_get_range_limit_invalid'
        table_meta = TableMeta(table_name, [('PK0', 'STRING')])
        reserved_throughput = ReservedThroughput(CapacityUnit(10, 10))

        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load('table_test_get_range_limit_invalid')

        pk_dict_small = {'PK0':'AAA'}
        pk_dict_big = {'PK0':'DDDD'}
        self._check_get_range_expect_exception(table_name, pk_dict_small, pk_dict_big, 'FORWARD', "The limit must be greater than 0.", "OTSParameterInvalid", 0)
        self._check_get_range_expect_exception(table_name, pk_dict_small, pk_dict_big, 'FORWARD', "The limit must be greater than 0.", "OTSParameterInvalid", -1)

    def _check_xget_range(self, table_name, inclusive_start_primary_keys, exclusive_end_primary_keys, direction, expect_rows, limit=10, columns=[]):
        consumed_count =  CapacityUnit(0, 0)
        row_size_sum = 0
        count = 0
        for row in self.client_test.xget_range(table_name, direction, inclusive_start_primary_keys,
                                               exclusive_end_primary_keys, consumed_count, columns, limit):
            if not row in expect_rows:
                self.assert_false()

            row_size_sum = row_size_sum + self.get_row_size(row[0], row[1]) 
            count = count + 1
        self.assert_equal(count, len(expect_rows))
        cu_read = int(math.ceil(row_size_sum * 1.0 / 4096))
        if cu_read == 0:
            cu_read = 1
        consumed_expect = CapacityUnit(cu_read, 0)
        self.assert_consumed(consumed_count, consumed_expect)

    def test_read_empty_row(self):
        """测试对空行的读操作，以及GetRange在行没有对应的列时期望不返回空行"""
        table_name = 'T'
        
        table_meta = TableMeta('T', [('PK0', 'INTEGER')])
        self.client_test.create_table(table_meta, ReservedThroughput(CapacityUnit(20, 20)))
        self.wait_for_partition_load('T')

        consumed = self.client_test.put_row('T', Condition('IGNORE'), {'PK0' : 0}, {})
        self.assert_consumed(consumed, CapacityUnit(0, 1))

        consumed = self.client_test.put_row('T', Condition('IGNORE'), {'PK0' : 1}, {'Col' : 1})
        self.assert_consumed(consumed, CapacityUnit(0, 1))

        consumed, pks, columns = self.client_test.get_row('T', {'PK0' : 0})
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(pks, {'PK0' : 0})
        self.assert_equal(columns, {})
        
        consumed, pks, columns = self.client_test.get_row('T', {'PK0' : 0}, columns_to_get=['Col'])
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(pks, {})
        self.assert_equal(columns, {})

        ret = self.client_test.batch_get_row([('T', [{'PK0' : 0}], None)])
        self.assert_RowDataItem_equal(ret, [[RowDataItem(True, None, None, CapacityUnit(1, 0), {'PK0' : 0}, {})]])
        
        ret = self.client_test.batch_get_row([('T', [{'PK0' : 0}], ['Col'])])
        self.assert_RowDataItem_equal(ret, [[RowDataItem(True, None, None, CapacityUnit(1, 0), {}, {})]])

        consumed, next_pk, row_list = self.client_test.get_range('T', 'FORWARD', {'PK0' : 0}, {'PK0' : 1})
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(next_pk, None)
        self.assert_equal(row_list, [({'PK0' : 0}, {})])

        consumed, next_pk, row_list = self.client_test.get_range('T', 'FORWARD', {'PK0' : 0}, {'PK0' : 1}, columns_to_get=['Col'])
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(next_pk, None)
        self.assert_equal(row_list, [])

        consumed, next_pk, row_list = self.client_test.get_range('T', 'FORWARD', {'PK0' : 0}, {'PK0' : 2})
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(next_pk, None)
        self.assert_equal(row_list, [({'PK0' : 0}, {}), ({'PK0' : 1}, {'Col' : 1})])

        consumed, next_pk, row_list = self.client_test.get_range('T', 'FORWARD', {'PK0' : 0}, {'PK0' : 2}, columns_to_get=['Col'])
        self.assert_consumed(consumed, CapacityUnit(1, 0))
        self.assert_equal(next_pk, None)
        self.assert_equal(row_list, [({}, {'Col' : 1})])

    def test_all_item_in_batch_write_row_failed(self):
        """当batch write row里的所有item全部后端失败时，期望每个item都返回具体的错误，而整个操作返回正常"""
    
        table1 = "table1"
        table_meta1 = TableMeta(table1, [("PK", "INTEGER")])
        reserved_throughput = ReservedThroughput(CapacityUnit(50, 50))
        self.client_test.create_table(table_meta1, reserved_throughput)

        self.wait_for_partition_load('table1')
        put_table1 = PutRowItem(Condition("EXPECT_EXIST"), {"PK": 11}, {"COL": "table1_11"})
        update_table1 = UpdateRowItem(Condition("EXPECT_EXIST"), {"PK": 12}, {"put" : {"COL": "table1_12"}})
        delete_table1 = DeleteRowItem(Condition("EXPECT_EXIST"), {"PK": 13})

        batches = [{
            'table_name':table1, 
            'put' : [put_table1], 
            'update' : [update_table1], 
            'delete' : [delete_table1],
        }]

        response = self.client_test.batch_write_row(batches)
        expect_res = [{
            'put':[BatchWriteRowResponseItem(False, "OTSConditionCheckFail", "Condition check failed.", None)],
            'update':[BatchWriteRowResponseItem(False, "OTSConditionCheckFail", "Condition check failed.", None)],
            'delete':[BatchWriteRowResponseItem(False, "OTSConditionCheckFail", "Condition check failed.", None)],

        }]
        self.assert_BatchWriteRowResponseItem(response, expect_res)

    def test_one_delete_in_update(self):
        """当一行为空时，进行一个UpdateRow，并且仅包含一个Delete操作"""

        table_name = 'T'

        table_meta = TableMeta(table_name, [("PK", "INTEGER")])
        reserved_throughput = ReservedThroughput(CapacityUnit(50, 50))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        cu = self.client_test.update_row(table_name, Condition('IGNORE'), {"PK" : 11}, {"delete" : ["Col0"]})
        self.assert_consumed(cu, CapacityUnit(0, 1))
        
    def test_all_delete_in_update(self):
        """当一行为空时，进行一个UpdateRow，并且包含128个Delete操作"""
        
        table_name = 'T'

        table_meta = TableMeta(table_name, [("PK", "INTEGER")])
        reserved_throughput = ReservedThroughput(CapacityUnit(50, 50))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        columns_to_delete = []
        for i in range(restriction.MaxColumnCountForRow):
            columns_to_delete.append('Col' + str(i))

        cu = self.client_test.update_row(table_name, Condition('IGNORE'), {"PK" : 11}, {"delete" : columns_to_delete})
        pk_size = 10 # 'PK' + int
        col_size = sum([len(col) for col in columns_to_delete])
        expect_write_cu = int(math.ceil((pk_size + col_size) * 1.0 / 4096))
        self.assert_consumed(cu, CapacityUnit(0, expect_write_cu))
    
        cu, pks, cols = self.client_test.get_row(table_name, {"PK" : 11})
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(pks, {})
        self.assert_equal(cols, {})

        
    def test_one_delete_in_update_of_batch_write(self):
        """当一行为空时，进行一个BatchWriteRow的UpdateRow，并且仅包含一个Delete操作"""
        
        table_name = 'T'

        table_meta = TableMeta(table_name, [("PK", "INTEGER")])
        reserved_throughput = ReservedThroughput(CapacityUnit(50, 50))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)

        update_item = UpdateRowItem(Condition("IGNORE"), {"PK": 12}, {"delete" : ["Col0"]})
        batches = [{'table_name' : table_name, 'update' : [update_item]}]
        response = self.client_test.batch_write_row(batches)
        self.assert_BatchWriteRowResponseItem(response, [{'update' : [BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, 1))]}])
    
        cu, pks, cols = self.client_test.get_row(table_name, {"PK" : 12})
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(pks, {})
        self.assert_equal(cols, {})

    
    def test_all_delete_in_update_of_batch_write(self):
        """当一行为空时，进行一个BatchWriteRow的UpdateRow，并且包含128个Delete操作"""
        
        table_name = 'T'

        table_meta = TableMeta(table_name, [("PK", "INTEGER")])
        reserved_throughput = ReservedThroughput(CapacityUnit(50, 50))
        self.client_test.create_table(table_meta, reserved_throughput)
        self.wait_for_partition_load(table_name)
    
        columns_to_delete = []
        for i in range(restriction.MaxColumnCountForRow):
            columns_to_delete.append('Col' + str(i))
    
        update_item = UpdateRowItem(Condition("IGNORE"), {"PK": 12}, {"delete" : columns_to_delete})
        batches = [{'table_name' : table_name, 'update' : [update_item]}]
        response = self.client_test.batch_write_row(batches)
        pk_size = 10 # 'PK' + int
        col_size = sum([len(col) for col in columns_to_delete])
        expect_write_cu = int(math.ceil((pk_size + col_size) * 1.0 / 4096))
        self.assert_BatchWriteRowResponseItem(response, [{'update' : [BatchWriteRowResponseItem(True, "", "", CapacityUnit(0, expect_write_cu))]}])
    
        cu, pks, cols = self.client_test.get_row(table_name, {"PK" : 12})
        self.assert_consumed(cu, CapacityUnit(1, 0))
        self.assert_equal(pks, {})
        self.assert_equal(cols, {})

if __name__ == '__main__':
    unittest.main()
