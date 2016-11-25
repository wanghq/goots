# -*- coding: utf8 -*-

from example_config import *
from ots2 import *
import time

table_name = 'BatchWriteRowExample'

def create_table(ots_client):
    schema_of_primary_key = [('gid', 'INTEGER'), ('uid', 'INTEGER')]
    table_meta = TableMeta(table_name, schema_of_primary_key)
    reserved_throughput = ReservedThroughput(CapacityUnit(0, 0))
    ots_client.create_table(table_meta, reserved_throughput)
    print 'Table has been created.'

def delete_table(ots_client):
    ots_client.delete_table(table_name)
    print 'Table \'%s\' has been deleted.' % table_name

def batch_write_row(ots_client):
    # batch put 10 rows and update 10 rows on exist table, delete 10 rows on a not-exist table.
    put_row_items = []
    for i in range(0, 10):
        primary_key = {'gid':i, 'uid':i+1}
        attribute_columns = {'name':'somebody'+str(i), 'address':'somewhere'+str(i), 'age':i}
        condition = Condition('IGNORE')
        item = PutRowItem(condition, primary_key, attribute_columns)
        put_row_items.append(item)

    update_row_items = []
    for i in range(10, 20):
        primary_key = {'gid':i, 'uid':i+1}
        attribute_columns = {'put': {'name':'somebody'+str(i), 'address':'somewhere'+str(i), 'age':i}}
        condition = Condition('IGNORE')
        item = UpdateRowItem(condition, primary_key, attribute_columns)
        update_row_items.append(item)

    delete_row_items = []
    for i in range(10, 20):
        primary_key = {'gid':i, 'uid':i+1}
        condition = Condition('IGNORE')
        item = DeleteRowItem(condition, primary_key)
        delete_row_items.append(item)

    batch_rows = {'table_name': table_name, 'put': put_row_items, 'update': update_row_items}
    batch_rows_of_another_table = {'table_name': 'notExistTable', 'delete': delete_row_items}
    result = ots_client.batch_write_row([batch_rows, batch_rows_of_another_table])

    print 'check first table\'s put results:'
    for item in result[0]['put']:
        if item.is_ok:
            print 'Put succeed, consume %s write cu.' % item.consumed.write
        else:
            print 'Put failed, error code: %s, error message: %s' % (item.error_code, item.error_message)

    print 'check first table\'s update results:'
    for item in result[0]['update']:
        if item.is_ok:
            print 'Update succeed, consume %s write cu.' % item.consumed.write
        else:
            print 'Update failed, error code: %s, error message: %s' % (item.error_code, item.error_message)

    print 'check second table\'s delete results:'
    for item in result[1]['delete']:
        if item.is_ok:
            print 'Delete succeed, consume %s write cu.' % item.consumed.write
        else:
            print 'Delete failed, error code: %s, error message: %s' % (item.error_code, item.error_message)

if __name__ == '__main__':
    ots_client = OTSClient(OTS_ENDPOINT, OTS_ID, OTS_SECRET, OTS_INSTANCE)
    create_table(ots_client)

    time.sleep(3) # wait for table ready
    batch_write_row(ots_client)
    delete_table(ots_client)

