#!/bin/python
# -*- coding: utf8 -*-

import logging
import unittest

from ots2.client import *
from ots2.metadata import *
from ots2.test.mock_connection import MockConnection

ENDPOINT = 'http://10.97.204.97:8800'
ACCESSID = 'accessid'
ACCESSKEY = 'accesskey'
INSTANCENAME = 'instancename'

class SDKMockTest(unittest.TestCase):

    def setUp(self):
        logger = logging.getLogger('test')
        handler=logging.FileHandler("test.log")
        formatter = logging.Formatter("[%(asctime)s]    [%(process)d]   [%(levelname)s] " \
                    "[%(filename)s:%(lineno)s]   %(message)s")
        handler.setFormatter(formatter)
        logger.addHandler(handler)
        logger.setLevel(logging.DEBUG)

        OTSClient.connection_pool_class = MockConnection
        self.ots_client = OTSClient(ENDPOINT, ACCESSID, ACCESSKEY, INSTANCENAME, logger_name='test')

    def tearDown(self):
        pass

    def test_list_table(self):
        table_list = self.ots_client.list_table();
        self.assertEqual(table_list[0], 'test_table1')
        self.assertEqual(table_list[1], 'test_table2')

    def test_create_table(self):
        table_meta = TableMeta('test_table', [('PK1', 'STRING'), ('PK2', 'INTEGER')])
        reserved_throughput = ReservedThroughput(CapacityUnit(10, 10))
        self.ots_client.create_table(table_meta, reserved_throughput)

        reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
        self.ots_client.create_table(table_meta, reserved_throughput)

    def test_delete_table(self):
        self.ots_client.delete_table('test_table')

    def test_update_table(self):
        reserved_throughput = ReservedThroughput(CapacityUnit(10))
        response = self.ots_client.update_table('test_table', reserved_throughput)
        self.assertEqual(response.reserved_throughput_details.capacity_unit.read, 10);
        self.assertEqual(response.reserved_throughput_details.capacity_unit.write, 10);
        self.assertEqual(response.reserved_throughput_details.last_increase_time, 123456);
        self.assertEqual(response.reserved_throughput_details.last_decrease_time, None);
        self.assertEqual(response.reserved_throughput_details.number_of_decreases_today, 5);

    def test_describe_table(self):
        response = self.ots_client.describe_table('test_table')
        self.assertEqual(response.table_meta.table_name, 'test_table')
        self.assertEqual(response.table_meta.schema_of_primary_key, [('PK1', 'STRING'), ('PK2', 'INTEGER')])
        self.assertEqual(response.reserved_throughput_details.capacity_unit.read, 1000);
        self.assertEqual(response.reserved_throughput_details.capacity_unit.write, 100);
        self.assertEqual(response.reserved_throughput_details.last_increase_time, 123456);
        self.assertEqual(response.reserved_throughput_details.last_decrease_time, 123456);
        self.assertEqual(response.reserved_throughput_details.number_of_decreases_today, 5);

    def test_put_row(self):
        condition = Condition('EXPECT_NOT_EXIST')
        primary_key = {'PK1':'hello', 'PK2':100}
        attribute_columns = {'COL1':'world', 'COL2':1000}
        consumed = self.ots_client.put_row('test_table', condition, primary_key, attribute_columns)
        self.assertEqual(consumed.read, 0)
        self.assertEqual(consumed.write, 10)
    
    def test_get_row(self):
        primary_key = {'PK1':'hello', 'PK2':100}
        columns_to_get = ['COL1', 'COL2']
        consumed, resp_pks, resp_attribute_columns = self.ots_client.get_row('test_table', primary_key, columns_to_get)
        self.assertEqual(consumed.read, 10)
        self.assertEqual(consumed.write, 0)
        self.assertEqual(resp_pks, {'PK1':'Hello', 'PK2':bytearray('World')})
        self.assertEqual(resp_attribute_columns, {'COL1':'test', 'COL2':100})

        consumed, resp_pks, resp_attribute_columns = self.ots_client.get_row('test_table', primary_key)
        self.assertEqual(consumed.read, 10)
        self.assertEqual(consumed.write, 0)
        self.assertEqual(resp_pks, {'PK1':'Hello', 'PK2':bytearray('World')})
        self.assertEqual(resp_attribute_columns, {'COL1':'test', 'COL2':100})

        consumed, resp_pks, resp_attribute_columns = self.ots_client.get_row('test_table', primary_key, [])
        self.assertEqual(consumed.read, 10)
        self.assertEqual(consumed.write, 0)
        self.assertEqual(resp_pks, {'PK1':'Hello', 'PK2':bytearray('World')})
        self.assertEqual(resp_attribute_columns, {'COL1':'test', 'COL2':100})

    def test_binary_type(self):
        bytearr = bytearray.fromhex('2F 77 39 00')
        primary_key = {'PK1':'hello', 'PK2':bytearr}
        columns_to_get = ['COL1', 'COL2']
        consumed, resp_pks, resp_attribute_columns = self.ots_client.get_row('test_table', primary_key, columns_to_get)
        self.assertEqual(consumed.read, 10)
        self.assertEqual(consumed.write, 0)
        self.assertEqual(resp_pks, {'PK1':'Hello', 'PK2':bytearray('World')})
        self.assertEqual(resp_attribute_columns, {'COL1':'test', 'COL2':100})

    def test_update_row(self):
        condition = Condition('IGNORE')
        primary_key = {'PK1':'hello', 'PK2':100}
        update_of_attribute_columns = {
            'put' : {'COL1':'test'},
            'delete' : ['COL2'],
        }
        consumed = self.ots_client.update_row('test_table', condition, primary_key, update_of_attribute_columns)
        self.assertEqual(consumed.read, 0)
        self.assertEqual(consumed.write, 10)

    def test_delete_row(self):
        condition = Condition('EXPECT_EXIST')
        primary_key = {'PK1':'hello', 'PK2':100}
        consumed = self.ots_client.delete_row('test_table', condition, primary_key)
        self.assertEqual(consumed.read, 0)
        self.assertEqual(consumed.write, 1)

    def test_batch_get_row(self):
        primary_key = {'PK1':'hello', 'PK2':100}
        columns_to_get = ['COL1', 'COL2']
        batch_list = [('test_table', [primary_key], columns_to_get)]
        response = self.ots_client.batch_get_row(batch_list)

        resp_row = response[0][0]
        self.assertEqual(resp_row.is_ok, True)
        self.assertEqual(resp_row.error_code, None)
        self.assertEqual(resp_row.error_message, None)
        self.assertEqual(resp_row.consumed.read, 100)
        self.assertEqual(resp_row.primary_key_columns, {'PK1':'Hello', 'PK2':100})
        self.assertEqual(resp_row.attribute_columns, {'COL1':'test', 'COL2':1000})

        resp_row = response[0][1]
        self.assertEqual(resp_row.is_ok, False)
        self.assertEqual(resp_row.error_code, 'ErrorCode')
        self.assertEqual(resp_row.error_message, 'ErrorMessage')
        self.assertEqual(resp_row.consumed, None)
        self.assertEqual(resp_row.primary_key_columns, None)
        self.assertEqual(resp_row.attribute_columns, None)

    def test_batch_write_row(self):
        primary_key = {'PK1':'hello', 'PK2':100}
        attribute_columns = {'COL1':'world', 'COL2':1000}
        condition = Condition('EXPECT_NOT_EXIST')
        put_row_item = PutRowItem(condition, primary_key, attribute_columns)

        condition = Condition('EXPECT_EXIST')
        attribute_columns = {'put': {'COL1':'world', 'COL2':1000}}
        update_row_item1 = UpdateRowItem(condition, primary_key, attribute_columns)
        primary_key = {'PK1':'world', 'PK2':101}
        attribute_columns = {'put': {'COL1':'hello', 'COL2':1001}}
        condition = Condition('IGNORE')
        update_row_item2 = UpdateRowItem(condition, primary_key, attribute_columns)

        condition = Condition('IGNORE')
        delete_row_item = DeleteRowItem(condition, primary_key)
        batch_list = [
            {
                'table_name': 'test_table',
                'put': [put_row_item],
                'update': [update_row_item1, update_row_item2],
                'delete': [delete_row_item]
            }
        ]
        response = self.ots_client.batch_write_row(batch_list)

        self.assertEqual(len(response), 1)
        put_resp_list = response[0]['put']
        self.assertEqual(len(put_resp_list), 1)
        resp_item = put_resp_list[0]
        self.assertEqual(resp_item.is_ok, True) 
        self.assertEqual(resp_item.error_code, None) 
        self.assertEqual(resp_item.error_message, None) 
        self.assertEqual(resp_item.consumed.read, 0) 
        self.assertEqual(resp_item.consumed.write, 10) 

        update_resp_list = response[0]['update']
        self.assertEqual(len(update_resp_list), 2)
        resp_item = update_resp_list[0]
        self.assertEqual(resp_item.is_ok, False) 
        self.assertEqual(resp_item.error_code, 'ErrorCode') 
        self.assertEqual(resp_item.error_message, 'ErrorMessage') 
        self.assertEqual(resp_item.consumed.read, 0) 
        self.assertEqual(resp_item.consumed.write, 100) 
        resp_item = update_resp_list[1]
        self.assertEqual(resp_item.is_ok, True) 
        self.assertEqual(resp_item.error_code, None) 
        self.assertEqual(resp_item.error_message, None) 
        self.assertEqual(resp_item.consumed.read, 0) 
        self.assertEqual(resp_item.consumed.write, 1111) 

        delete_resp_list = response[0]['delete']
        self.assertEqual(len(delete_resp_list), 1)
        resp_item = delete_resp_list[0]
        self.assertEqual(resp_item.is_ok, True) 
        self.assertEqual(resp_item.error_code, None) 
        self.assertEqual(resp_item.error_message, None) 
        self.assertEqual(resp_item.consumed.read, 0) 
        self.assertEqual(resp_item.consumed.write, 1000) 

    def test_get_range(self):
        start_primary_key = {'PK1':'hello', 'PK2':100}
        end_primary_key = {'PK1':INF_MAX, 'PK2':INF_MIN}
        columns_to_get = ['COL1', 'COL2']
        consumed, next_start_pk, row_list = self.ots_client.get_range(
                    'table_name', 'FORWARD', 
                    start_primary_key, end_primary_key, 
                    columns_to_get, limit=100
        )

        self.assertEqual(consumed.read, 100)
        self.assertEqual(consumed.write, 0)
        self.assertEqual(next_start_pk, {'PK1':'NextStart', 'PK2':101})
        (primary_key_columns, attribute_columns) = row_list[0] 
        self.assertEqual(primary_key_columns, {'PK1':'Hello', 'PK2':100})
        self.assertEqual(attribute_columns, {'COL1':'test', 'COL2':1000})

        consumed, next_start_pk, row_list = self.ots_client.get_range(
                    'table_name', 'FORWARD', 
                    start_primary_key, end_primary_key, 
                    None, limit=100
        )

        self.assertEqual(consumed.read, 100)
        self.assertEqual(consumed.write, 0)
        self.assertEqual(next_start_pk, {'PK1':'NextStart', 'PK2':101})
        (primary_key_columns, attribute_columns) = row_list[0] 
        self.assertEqual(primary_key_columns, {'PK1':'Hello', 'PK2':100})
        self.assertEqual(attribute_columns, {'COL1':'test', 'COL2':1000})

        consumed, next_start_pk, row_list = self.ots_client.get_range(
                    'table_name', 'FORWARD', 
                    start_primary_key, end_primary_key, 
                    [], limit=100
        )

        self.assertEqual(consumed.read, 100)
        self.assertEqual(consumed.write, 0)
        self.assertEqual(next_start_pk, {'PK1':'NextStart', 'PK2':101})
        (primary_key_columns, attribute_columns) = row_list[0] 
        self.assertEqual(primary_key_columns, {'PK1':'Hello', 'PK2':100})
        self.assertEqual(attribute_columns, {'COL1':'test', 'COL2':1000})

    def test_xget_range(self):
        consumed_counter = CapacityUnit(0, 0)
        start_primary_key = {'PK1':'hello', 'PK2':100}
        end_primary_key = {'PK1':INF_MAX, 'PK2':INF_MIN}
        columns_to_get = ['COL1', 'COL2']
        range_iterator = self.ots_client.xget_range('table_name', 'FORWARD', 
                    start_primary_key, end_primary_key, consumed_counter,
                    columns_to_get, limit=100
        )

        for (primary_key_columns, attribute_columns) in range_iterator:
            self.assertEqual(primary_key_columns, {'PK1':'Hello', 'PK2':100})
            self.assertEqual(attribute_columns, {'COL1':'test', 'COL2':1000})
            if consumed_counter.read >= 300:
                break
        self.assertEqual(consumed_counter.read, 300)

if __name__ == '__main__':
    unittest.main()
