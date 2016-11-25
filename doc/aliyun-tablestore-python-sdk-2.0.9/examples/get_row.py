# -*- coding: utf8 -*-

from example_config import *
from ots2 import *
import time

table_name = 'GetRowExample'

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
    primary_key = {'gid':1, 'uid':101}
    attribute_columns = {'name':'John', 'mobile':15100000000, 'address':'China', 'age':20}
    condition = Condition('EXPECT_NOT_EXIST') # Expect not exist: put it into table only when this row is not exist.
    consumed = ots_client.put_row(table_name, condition, primary_key, attribute_columns)
    print u'Write succeed, consume %s write cu.' % consumed.write

def get_row(ots_client):
    primary_key = {'gid':1, 'uid':101}
    columns_to_get = ['name', 'address', 'age'] # given a list of columns to get, or empty list if you want to get entire row.
    consumed, primary_key_columns, attribute_columns = ots_client.get_row(table_name, primary_key, columns_to_get)
    print u'Read succeed, consume %s read cu.' % consumed.read

    print u'Value of attribute \'name\': %s' % attribute_columns.get('name')
    print u'Value of attribute \'address\': %s' % attribute_columns.get('address')
    print u'Value of attribute \'age\': %s' % attribute_columns.get('age')

if __name__ == '__main__':
    ots_client = OTSClient(OTS_ENDPOINT, OTS_ID, OTS_SECRET, OTS_INSTANCE)
    create_table(ots_client)

    time.sleep(3) # wait for table ready
    put_row(ots_client)
    get_row(ots_client)
    delete_table(ots_client)

