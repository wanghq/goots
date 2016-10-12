# -*- coding: utf8 -*-

from example_config import *
from ots2 import *
import time

table_name = 'BatchGetRowExample'

def create_table(ots_client):
    schema_of_primary_key = [('gid', 'INTEGER'), ('uid', 'INTEGER')]
    table_meta = TableMeta(table_name, schema_of_primary_key)
    reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
    ots_client.create_table(table_meta, reserved_throughput)
    print 'Table has been created.'

def delete_table(ots_client):
    ots_client.delete_table(table_name)
    print 'Table \'%s\' has been deleted.' % table_name

def put_row(ots_client):
    for i in range(0, 10):
        primary_key = {'gid':i, 'uid':i+1}
        attribute_columns = {'name':'John', 'mobile':i, 'address':'China', 'age':i}
        condition = Condition('EXPECT_NOT_EXIST') # Expect not exist: put it into table only when this row is not exist.
        consumed = ots_client.put_row(table_name, condition, primary_key, attribute_columns)
        print u'Write succeed, consume %s write cu.' % consumed.write

def batch_get_row(ots_client):
    # try get 10 rows from exist table and 10 rows from not-exist table
    columns_to_get = ['git', 'uid', 'name', 'mobile', 'address', 'age']
    rows_to_get = []
    for i in range(0, 10):
        primary_key = {'gid':i, 'uid':i+1}
        rows_to_get.append(primary_key)

    batch_rows = (table_name, rows_to_get, columns_to_get)
    result = ots_client.batch_get_row([batch_rows, ('notExistTable', rows_to_get, [])])

    print 'Check first table\'s result:'
    for item in result[0]:
        if item.is_ok:
            print 'Read succeed, PrimaryKey: %s, Attributes: %s' % (item.primary_key_columns, item.attribute_columns)
        else:
            print 'Read failed, error code: %s, error message: %s' % (item.error_code, item.error_message)

    print 'Check second table\'s result:'
    for item in result[1]:
        if item.is_ok:
            print 'Read succeed, PrimaryKey: %s, Attributes: %s' % (item.primary_key_columns, item.attribute_columns)
        else:
            print 'Read failed, error code: %s, error message: %s' % (item.error_code, item.error_message)

if __name__ == '__main__':
    ots_client = OTSClient(OTS_ENDPOINT, OTS_ID, OTS_SECRET, OTS_INSTANCE)
    create_table(ots_client)

    time.sleep(3) # wait for table ready
    put_row(ots_client)
    batch_get_row(ots_client)
    delete_table(ots_client)

